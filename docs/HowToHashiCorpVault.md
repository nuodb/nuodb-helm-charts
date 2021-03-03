# Using HashiCorp Vault for Management of TLS Certificates

## Introduction

NuoDB supports TLS encryption for all processes in the domain.
`NuoDB Admin` is responsible for propagating certificates to database processes, so to enable TLS encryption for all processes, it is necessary to configure NuoDB Admin with a set of certificates, and also configure NuoDB Command (`nuocmd`) clients to be able to communicate with `NuoDB Admin`.

This document expands on [Security Model of NuoDB in Kubernetes](./HowtoTLS.md) and explains the usage of HashiCorp Vault for the management of TLS keys.

HashiCorp uses sidecar containers to inject keys into another container.
For more info, explore the blog post [Injecting Vault Secrets](https://www.hashicorp.com/blog/injecting-vault-secrets-into-kubernetes-pods-via-a-sidecar)


> **NOTE**: For information on enabling  TLS encryption in non-Kubernetes deployments of NuoDB, see [here](http://doc.nuodb.com/Latest/Content/Nuoadmin-Configuring-TLS-Security.htm). 
This document expands on the product documentation and is specific to this Helm Chart repository.

### Terminology

- `Key` = a combination of a private key with its corresponding X509 certificate chain. 
These are usually saved in a PKCS12 file such as `nuoadmin.p12`.
- `NuoDB Admin` = admin interface for domain and database management. Started by the [Admin Chart](../stable/admin/README.md).
- `CA` = Certificate Authority

## Generating NuoDB Keys
You can either create your own TLS keys or create them using the convenience script provided with the docker image:

```
docker run --rm -d --name create-tls-keys nuodb/nuodb-ce:latest -- tail -f /dev/null 

docker exec -it create-tls-keys bash -c "mkdir /tmp/keys && \ 
cd /tmp/keys && DEFAULT_PASSWORD=changeIt setup-keys.sh"

docker cp create-tls-keys:/tmp/keys /tmp/keys
docker stop create-tls-keys
```

The convenience script will generate 4 files:
- `nuodb-keystore.p12` containing the X509 that identifies the admin;
- `nuodb-truststore.p12` usually containing the root CA and the primordial admin user;
- `ca.cert` containing the public certificate of the Certificate Authority;
- `nuocmd.pem` which is the private/public keypair for the primordial admin user.

## Installation of HashiCorp Vault

This document will use the `dev` mode of HashiCorp vault for simplicity.
This mode is not suited for production and `seals` should be used instead.
For more info on `HashiCorp tokens` please refer to the [Documentation](https://www.vaultproject.io/docs/concepts/seal)

### Install Vault using Helm

First, install the HashiCorp Vault official helm chart.

```
kubectl create namespace vault
helm repo add hashicorp https://helm.releases.hashicorp.com
helm install vault hashicorp/vault -n vault --set server.dev.enabled=true
```

### Create Vault Policy
Next, connect to Vault and configure a policy named “nuodb-policy”.
This is a very non-restrictive policy, and in a production setting, you would typically want to lock this down more, but it serves as an example while you play around with this feature.

```
kubectl exec -it -n vault vault-0 -- sh
$ cat <<EOF > /home/vault/nuodb-policy.hcl
path "nuodb.com*" {
  capabilities = ["read"]
}
EOF

$ vault policy write nuodb-policy /home/vault/nuodb-policy.hcl
```

### Enable Kubernetes Integration

The next step is to enable the [Vault Kubernetes Auth](https://www.vaultproject.io/docs/auth/kubernetes) method.
The `nuodb-policy` created above in step [Create Vault Policy](###create-vault-policy) will get attached to the `nuodb` namespace.

```
kubectl exec -it -n vault vault-0 -- sh
$ vault auth enable kubernetes

$ vault write auth/kubernetes/config \
   token_reviewer_jwt="$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)" \
   kubernetes_host=https://${KUBERNETES_PORT_443_TCP_ADDR}:443 \
   kubernetes_ca_cert=@/var/run/secrets/kubernetes.io/serviceaccount/ca.crt

$ vault write auth/kubernetes/role/nuodb \
   bound_service_account_names=nuodb \
   bound_service_account_namespaces=nuodb \
   policies=nuodb-policy \
   ttl=1h
```

### Create a NuoDB Secrets Engine

The final step is to create a [Secrets KV Engine](https://www.vaultproject.io/docs/secrets/kv/index.html) to store the NuoDB keys.
All keys will be base64 encoded.

```
kubectl exec -it -n vault vault-0 -- sh
$ vault secrets enable -version=2 -path=nuodb.com kv
```

## Adding NuoDB Keys to Vault

In the previous step [Generating NuoDB Keys](##generating-nuodb-keys) we generated a set of keys required to run NuoDB.
We saved those keys in `/tmp/keys`.
In this step, we will save these keys in HC Vault.

```
kubectl exec -it -n vault vault-0 -- \
    vault kv put nuodb.com/TLS \
    tlsClientPEM=`cat /tmp/keys/nuocmd.pem | base64` \
    tlsCACert=`cat /tmp/keys/ca.cert | base64` \
    tlsKeyStorePassword=changeIt \
    tlsTrustStorePassword=changeIt \
    tlsKeyStore=`cat /tmp/keys/nuoadmin.p12 | base64` \
    tlsTrustStore=`cat /tmp/keys/nuoadmin-truststore.p12 | base64`
```

## Configuration of NuoDB

HashiCorp Vault uses annotations to identify pods that require a Vault Agent.
For more info please consult the [Template Documentation](https://www.vaultproject.io/docs/platform/k8s/injector/annotations#vault-hashicorp-com-agent-inject-template)

### Annotations for the NuoDB Admin Tier

The NuoDB admin tier requires all 4 key files. For convenience we list them here again:
- `nuodb-keystore.p12` containing the X509 that identifies the admin;
- `nuodb-truststore.p12` usually containing the root CA and the primordial admin user;
- `ca.cert` containing the public certificate of the Certificate Authority;
- `nuocmd.pem` which is the private/public keypair for the primordial admin user.

Both the keystore and the truststore PKCS12 files are password protected.
The passwords need to be exported as environmental variables.

```
$ cat vault-annotations-admin.yaml
admin:
  podAnnotations:
    vault.hashicorp.com/agent-inject: "true"
    vault.hashicorp.com/agent-inject-secret-ca.cert: nuodb.com/TLS
    vault.hashicorp.com/agent-inject-secret-nuoadmin-truststore.p12: nuodb.com/TLS
    vault.hashicorp.com/agent-inject-secret-nuoadmin-truststore.password: nuodb.com/TLS
    vault.hashicorp.com/agent-inject-secret-nuoadmin.p12: nuodb.com/TLS
    vault.hashicorp.com/agent-inject-secret-nuoadmin.password: nuodb.com/TLS
    vault.hashicorp.com/agent-inject-secret-nuocmd.pem: nuodb.com/TLS
    vault.hashicorp.com/agent-inject-template-ca.cert: |
      {{- with secret "nuodb.com/TLS" -}}
        {{ .Data.data.tlsCACert | base64Decode }}
      {{- end }}
    vault.hashicorp.com/agent-inject-template-nuoadmin-truststore.p12: |
      {{- with secret "nuodb.com/TLS" -}}
        {{ .Data.data.tlsTrustStore | base64Decode }}
      {{- end }}
    vault.hashicorp.com/agent-inject-template-nuoadmin-truststore.password: |
      {{- with secret "nuodb.com/TLS" -}}
        export NUODB_TRUSTSTORE_PASSWORD=”{{ .Data.data.tlsTrustStorePassword }}”
      {{- end }}
    vault.hashicorp.com/agent-inject-template-nuoadmin.p12: |
      {{- with secret "nuodb.com/TLS" -}}
        {{ .Data.data.tlsKeyStore | base64Decode }}
      {{- end }}
    vault.hashicorp.com/agent-inject-template-nuoadmin.password: |
      {{- with secret "nuodb.com/TLS" -}}
        export NUODB_KEYSTORE_PASSWORD=”{{ .Data.data.tlsKeyStorePassword }}”
      {{- end }}
    vault.hashicorp.com/agent-inject-template-nuocmd.pem: |
      {{- with secret "nuodb.com/TLS" -}}
        {{ .Data.data.tlsClientPEM | base64Decode }}
      {{- end }}
    vault.hashicorp.com/role: nuodb
    vault.hashicorp.com/secret-volume-path: /etc/nuodb/keys
```

Start the NuoDB admin tier with 3 processes and the Vault annotations:
```
$ helm install -n nuodb --set admin.replicas=3 -f vault-annotations-admin.yaml admin nuodb/admin
```

The NuoDB admin pods should now contain two init containers and two containers.
Validate that the pods are ready:
```
$ kubectl get pods -n nuodb
NAME                     READY   STATUS    RESTARTS   AGE
admin-nuodb-cluster0-0   2/2     Running   0          73s
admin-nuodb-cluster0-1   2/2     Running   0          73s
admin-nuodb-cluster0-2   2/2     Running   0          73s
``` 

Validate that the domain is healthy using nuocmd.
```
$ kubectl exec -it -n nuodb admin-nuodb-cluster0-0 -c admin -- nuocmd show domain
server version: 4.0.8-2-881d0e5d44, server license: Community
server time: 2021-03-03T20:30:53.909, client token: 23ce1d3ac8ce652a6cb6aa3f7df1918538326c4e
Servers:
  [admin-nuodb-cluster0-0] admin-nuodb-cluster0-0.nuodb.nuodb.svc.cluster.local:48005 [last_ack = 1.73] [member = ADDED] [raft_state = ACTIVE] (LEADER, Leader=admin-nuodb-cluster0-0, log=0/13/13) Connected *
  [admin-nuodb-cluster0-1] admin-nuodb-cluster0-1.nuodb.nuodb.svc.cluster.local:48005 [last_ack = 1.73] [member = ADDED] [raft_state = ACTIVE] (FOLLOWER, Leader=admin-nuodb-cluster0-0, log=0/13/13) Connected
  [admin-nuodb-cluster0-2] admin-nuodb-cluster0-2.nuodb.nuodb.svc.cluster.local:48005 [last_ack = 1.73] [member = ADDED] [raft_state = ACTIVE] (FOLLOWER, Leader=admin-nuodb-cluster0-0, log=0/13/13) Connected
Databases:
```

Validate that there are no errors in the Vault container.
```
$ kubectl logs -n nuodb  admin-nuodb-cluster0-0 -c vault-agent
```

### Annotations for the NuoDB Engine Tier

With the admin running, we can now start the NuoDB engines (Storage Managers and Transaction Engines).
To do so, we will use the Database Helm chart with the following Vault annotations:
```
$ cat vault-annotations-database.yaml
database:
  podAnnotations:
    vault.hashicorp.com/agent-inject: "true"
    vault.hashicorp.com/agent-inject-secret-ca.cert: nuodb.com/TLS
    vault.hashicorp.com/agent-inject-secret-nuocmd.pem: nuodb.com/TLS
    vault.hashicorp.com/agent-inject-template-ca.cert: |
      {{- with secret "nuodb.com/TLS" -}}
        {{ .Data.data.tlsCACert | base64Decode }}
      {{- end }}
    vault.hashicorp.com/agent-inject-template-nuocmd.pem: |
      {{- with secret "nuodb.com/TLS" -}}
        {{ .Data.data.tlsClientPEM | base64Decode }}
      {{- end }}
    vault.hashicorp.com/role: nuodb
    vault.hashicorp.com/secret-volume-path: /etc/nuodb/keys
```

The engine pod only requires the client credentials (`nuocmd.pem`) and the public CA certificate (`ca.cert`).
These files are not password protected.

```
helm install -n nuodb -f vault-annotations-database.yaml database nuodb/database
```

The NuoDB domain should now consist of 3 admin tier pods, 1 TE and 1 SM.
All pods should contain 2 init containers and 2 containers.
Validate that the pods are ready:
```
$ kubectl get pods -n nuodb
NAME                                               READY   STATUS    RESTARTS   AGE
admin-nuodb-cluster0-0                             2/2     Running   0          10m
admin-nuodb-cluster0-1                             2/2     Running   0          10m
admin-nuodb-cluster0-2                             2/2     Running   0          10m
sm-database-nuodb-cluster0-demo-hotcopy-0          2/2     Running   0          99s
te-database-nuodb-cluster0-demo-556697c994-z6m5d   2/2     Running   0          99s
``` 

Validate that the domain is healthy using nuocmd.
```
$ kubectl exec -it -n nuodb admin-nuodb-cluster0-0 -c admin -- nuocmd show domain
server version: 4.0.8-2-881d0e5d44, server license: Community
server time: 2021-03-03T20:39:31.002, client token: 20f397da114c6e55567ea5d2f53660941f308bea
Servers:
  [admin-nuodb-cluster0-0] admin-nuodb-cluster0-0.nuodb.nuodb.svc.cluster.local:48005 [last_ack = 1.27] [member = ADDED] [raft_state = ACTIVE] (LEADER, Leader=admin-nuodb-cluster0-0, log=0/29/29) Connected *
  [admin-nuodb-cluster0-1] admin-nuodb-cluster0-1.nuodb.nuodb.svc.cluster.local:48005 [last_ack = 1.27] [member = ADDED] [raft_state = ACTIVE] (FOLLOWER, Leader=admin-nuodb-cluster0-0, log=0/29/29) Connected
  [admin-nuodb-cluster0-2] admin-nuodb-cluster0-2.nuodb.nuodb.svc.cluster.local:48005 [last_ack = 1.27] [member = ADDED] [raft_state = ACTIVE] (FOLLOWER, Leader=admin-nuodb-cluster0-0, log=0/29/29) Connected
Databases:
  demo [state = RUNNING]
    [SM] sm-database-nuodb-cluster0-demo-hotcopy-0/10.1.2.17:48006 [start_id = 0] [server_id = admin-nuodb-cluster0-1] [pid = 118] [node_id = 1] [last_ack =  9.22] MONITORED:RUNNING
    [TE] te-database-nuodb-cluster0-demo-556697c994-z6m5d/10.1.2.16:48006 [start_id = 1] [server_id = admin-nuodb-cluster0-2] [pid = 44] [node_id = 2] [last_ack =  3.20] MONITORED:RUNNING
```