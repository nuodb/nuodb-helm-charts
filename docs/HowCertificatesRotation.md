# Rotating NuoDB TLS Certificates

## Introduction

When using TLS encryption, it is necessary to rotate key pair certificates before the certificates expire.
To avoid downtime, `NuoDB Admin` processes (APs) with the new certificates should be introduced in a rolling fashion.

This document expands on [Security Model of NuoDB in Kubernetes](./HowtoTLS.md) and explains how to do proper certificate rotation in NuoDB Kubernetes deployments.

> **NOTE**: For information on rotating TLS keys in non-Kubernetes deployments of NuoDB, see [here](https://doc.nuodb.com/nuodb/latest/deployment-models/physical-or-vmware-environments-with-nuodb-admin/domain-operations/enabling-tls-encryption/rotating-key-pair-certificates/). 
This document expands on the product documentation and is specific to this Helm Chart repository.

### Terminology

- `Key` = a combination of a private key with its corresponding X509 certificate chain. 
These are usually saved in a PKCS12 file such as `nuoadmin.p12`.
- `NuoDB Admin` = admin interface for domain and database management.
For information on how to start the NuoDB Admin tier, see [Admin Chart](../stable/admin/README.md).
- `CA` = Certificate Authority

## About Key Rotation

By default in Kubernetes deployments, the same key pair certificate is used by all NuoDB Admin processes (APs), a Certificate Authority (CA) is used to sign them and the CA certificate is trusted by all processes in the domain.
This means that if a CA certificate needs to be renewed both the keystore and the truststore have to be updated for all APs, and by extension, all database processes.

As stated, processes with the new certificates should be introduced in a rolling fashion to avoid downtime, which means that processes with the old key pair certificate can verify processes with the new key pair certificate. Since the new key pair certificates must be verified using a new trusted certificate, existing processes (admin and database) must have the new trusted certificate propagated to their truststores before the new key pair certificates can be introduced.

An overview of how key rotation is performed is as follows:

1. A new shared admin new key pair certificate (and associated client PEM file) is generated.
2. The new certificate is added to the truststore of every AP and every database process.
3. The new certificate is added to the truststore of every SQL client and NuoDB Command client.
4. For each AP, the keystore is replaced. (Make sure that the new key pair certificate is in effect.)
5. For each database process, the process is restarted, so that its certificate chain is based on the new certificate.
6. (Optional) The old certificate is removed from the truststore of every AP and every database process.
7. (Optional) The old certificate is removed from the truststore of every SQL client and NuoDB Command client.

Usually only the server key pair certificate is renewed, which means that the keystore file has to be updated for all APs and database processes.
This reduces the steps needed during key rotation as the new key pair certificate can be verified using the old truststore certificate.
In particular steps 2 and 3 from the overview above are skipped.

### Update Key Pair Certificates

Depending on the used TLS keys management, step 4 in the overview above may be different.
By default NuoDB uses Kubernetes secrets to store TLS keys and expose them to NuoDB processes.

#### Using Kubernetes Secrets

NuoDB APs can reload the keystore without having to be reconfigured and restarted, however, in Kubernetes deployments the TLS keys are mounted in AP containers as a `subPath` volume mount which doesn't allow automatic secret update.
To replace the keystore for all APs, use the following steps:

Create new Kubernetes secrets using the keystore file generated with the renewed server key pair certificates.

```bash
kubectl create secret generic nuodb-keystore-renewed \
  --namespace nuodb \
  --from-file=nuoadmin.p12=/tmp/keys/nuoadmin.p12 \
  --from-literal=password=${PASSWD} -n nuodb

kubectl create secret generic nuodb-client-pem-renewed \
  --namespace nuodb \
  --from-file=nuocmd.pem=/tmp/keys/nuocmd.pem \
  --from-literal=password=${PASSWD} -n nuodb
```

Upgrade the Helm release installed with the [admin](../stable/admin/README.md) chart using the new TLS secrets.
This will perform NuoDB Admin statefulset rolling upgrade so make sure that you are having enough APs to prevent downtime.

```bash
helm upgrade admin nuodb/admin \
  --namespace nuodb \
  --set admin.tlsKeyStore.secret=nuodb-keystore-renewed \
  --set admin.tlsClientPEM.secret=nuodb-client-pem-renewed \
  -f values.yaml
```

Make sure that the new key pair certificates are used by all APs as described in [Verify Domain Certificates](#verify-domain-certificates).

Upgrade the Helm release installed with the [database](../stable/database/README.md) chart using the new TLS secrets.
This will perform database Storage Manager (SM) statefulsets and Transaction Engine (TE) deployment rolling upgrade so make sure that you are having enough APs to prevent downtime.

```bash
helm upgrade database nuodb/database \
  --namespace nuodb \
  --set admin.tlsKeyStore.secret=nuodb-keystore-renewed \
  --set admin.tlsClientPEM.secret=nuodb-client-pem-renewed \
  --set admin.tlsCACert.secret=nuodb-ca-cert-renewed \
  -f values.yaml
```

Make sure that the new key pair certificates are used by all domain processes as described in [Verify Domain Certificates](#verify-domain-certificates).

The client PEM secret and the CA certificate generation can be skipped if they are not rotated.

#### Using HashiCorp Vault

HashiCorp uses sidecar containers to inject keys into another container.
It automatically updates secret values stored in Vault KV store to match the files in the volume mounted inside the container.
This means that it is enough to update the keystore, truststore and client PEM in HashiCorp Vault. 

Ensure that the keystore has the same password (as specified by `admin.tlsKeyStore.password` Helm option) so that the APs can reload the keystore without having to be reconfigured and restarted.

For more information on how to configure TLS with HashiCorp Vault, see [Using HashiCorp Vault for Management of TLS Certificates](./HowToHashiCorpVault.md).

Make sure that the new key pair certificates are used by all APs as described in [Verify Domain Certificates](#verify-domain-certificates).

Restart database processes so that the APs generate new engine certificates signed with the new server certificate chain in the `Intermediate CA` model or the new shared key pair certificates are used in the `Pass-down` model.
Kubernetes 1.15+ supports a rollout restart which can be used to restart database processes in a rolling fashion.

```bash
kubectl rollout restart \
  --namespace nuodb \
  statefulset sm-database-nuodb-cluster0-demo-hotcopy

kubectl rollout restart \
  --namespace nuodb \
  statefulset sm-database-nuodb-cluster0-demo

kubectl rollout restart \
  --namespace nuodb \
  te-database-nuodb-cluster0-demo
```

Alternatively, the database processes can be restarted manually by shutting down the process using `nuocmd shutdown process --start-id ...` command.
Kubernetes will automatically restart the engine container and which will effectively restart the database process.

It is recommended to perform each rollout or database process restart sequentially after verifying the domain process certificates as described in [Verify Domain Certificates](#verify-domain-certificates).

### Verify Domain Certificates

To get the current certificate data, use `nuocmd get certificate-info` command and make sure that the renewed key pair certificates are displayed in the corresponding section.

```bash
kubectl exec -ti admin-nuodb-cluster0-0 \
  --namespace nuodb -- \
  nuocmd --show-json get certificate-info
```

The command displays the key pair or certificate information from the domain in the bellow sections:

- `serverCertificates` - information about certificates used by each AP identified by `serverId`.
- `serverTrusted` - trusted certificate aliases (referenced in `trustedCertificates` section) by each AP identified by `serverId`.
- `processCertificates` - information about certificates used by each database process indentified by `startId`.
- `processTrusted` - trusted certificate aliases (referenced in `trustedCertificates` section) by each database process identified by `startId`.
- `trustedCertificates` - information about trusted certificates in the domain identified by `alias`.

Example output:

```json
{
  /* output omitted */
  "serverCertificates": {
    "admin-nuodb-cluster0-0": {
      "certificatePem": "-----BEGIN CERTIFICATE-----...-----END CERTIFICATE-----",
      "expires": 1591301400000,
      "issuerName": "CN=*.domain",
      "subjectName": "CN=*.domain"
    },
    "admin-nuodb-cluster0-1": {
      "certificatePem": "-----BEGIN CERTIFICATE-----...-----END CERTIFICATE-----",
      "expires": 1591301400000,
      "issuerName": "CN=*.domain",
      "subjectName": "CN=*.domain"
    }
  },
  /* output omitted */
  "processCertificates": {
    "0": {
      "certificatePem": "-----BEGIN CERTIFICATE-----...-----END CERTIFICATE-----",
      "expires": 1591301400000,
      "issuerName": "CN=*.domain",
      "subjectName": "CN=server0"
    },
    "1": {
      "certificatePem": "-----BEGIN CERTIFICATE-----...-----END CERTIFICATE-----",
      "expires": 1591301400000,
      "issuerName": "CN=*.domain",
      "subjectName": "CN=server1"
    },
    "2": {
      "certificatePem": "-----BEGIN CERTIFICATE-----...-----END CERTIFICATE-----",
      "expires": 1591301400000,
      "issuerName": "CN=*.domain",
      "subjectName": "CN=server0"
    },
    "3": {
      "certificatePem": "-----BEGIN CERTIFICATE-----...-----END CERTIFICATE-----",
      "expires": 1591301400000,
      "issuerName": "CN=*.domain",
      "subjectName": "CN=server1"
    }
  },
}
```

## NuoDB TLS Keys Rotation

NuoDB supports various key models for TLS keys used by the `NuoDB Admin` and database processes.

> **NOTE**: For information on how to configure different TLS key models, see [Security Model of NuoDB in Kubernetes](./HowtoTLS.md).

To renew domain keystore and key pair certificates, you can either create new keystore files and certificates on your own or create them using the NuoDB Commands `nuocmd`.
All examples bellow will be using a helper pod started with NuoDB image and `nuocmd` command line tools.

```bash
mkdir /tmp/keys

kubectl run generate-nuodb-certs \
  --image nuodb/nuodb-ce \
  --env="PASSWD=changeMe" \
  --command -- 'tail' '-f' '/dev/null'

kubectl exec -ti generate-nuodb-certs -- \
  mkdir -p /tmp/keys
```

### Rotate Server Certificates

Generate a new server key pair:

```bash
kubectl exec -ti generate-nuodb-certs -- \
  nuocmd create keypair \
    --keystore /tmp/keys/nuoadmin.p12 --store-password "$PASSWD" --ca \
    --dname "CN=*.nuodb.svc.cluster.local, OU=Eng, O=NuoDB, L=Boston, ST=MA, C=US, SERIALNUMBER=67890"
```

In the example above, we specify the `SERIALNUMBER` field so that the new certificate has a different distinguished name from the current certificate.

The new server certificates should be signed by the CA certificate.
If the TLS keys have been generated using the `setup-keys.sh` script, the private/public key pair of the CA are available in the `ca.p12` file.

Generate and sign the new server certificates using the CA certificate chain:

```bash
kubectl cp ca.p12 generate-nuodb-certs:/tmp/keys/

kubectl exec -ti generate-nuodb-certs -- \
  nuocmd sign certificate \
    --keystore /tmp/keys/nuoadmin.p12 --store-password "$PASSWD" \
    --ca-keystore /tmp/keys/ca.p12 --ca-store-password "$PASSWD" \
    --validity 36500 --ca --update
```

The `--ca` argument specified weather the `isCA` extension is enabled in the generated certificate.
This allows the admin to act as intermediate CA and is used in the [Shared Admin Key + Intermediate CA](#shared-admin-key-+-intermediate-ca) model.
You can omit this argument in case the [Shared Admin Key + Pass-down](#share-admin-key-+-pass-down) model is used.

The newly created keystore files can be copied on the client machine which will be used later to update the keystore.

```bash
kubectl cp generate-nuodb-certs:/tmp/keys/. /tmp/keys
```

If the server certificates are signed by a public CA, then follow the provided steps to renew your certificates.

> **NOTE**: For information on how to create the server keystore when using certificates signed by a Public CA, see [here](https://doc.nuodb.com/nuodb/latest/deployment-models/physical-or-vmware-environments-with-nuodb-admin/domain-operations/enabling-tls-encryption/using-certificates-signed-by-a-public-certificate-authority/).

Update domain keystore with the renewed certificates as described in section [Update Key Pair Certificates](#update-key-pair-certificates).

### Rotate NuoDB Commands Certificate

The NuoDB Commands client key can be renewed together with the NuoDB Admin server key pair certificates or separately as needed.

Generate new client key pair:

```bash
kubectl exec -ti generate-nuodb-certs -- \
  nuocmd create keypair \
    --keystore /tmp/keys/nuocmd.p12 --store-password "$PASSWD" \
    --dname "CN=nuocmd.nuodb.com, OU=Eng, O=NuoDB, L=Boston, ST=MA, C=US, SERIALNUMBER=67890" --validity 36500

kubectl exec -ti generate-nuodb-certs -- \
  nuocmd show certificate \
    --keystore /tmp/keys/nuocmd.p12 --store-password "$PASSWD" > nuocmd.pem

kubectl exec -ti generate-nuodb-certs -- \
  nuocmd show certificate \
    --keystore /tmp/keys/nuocmd.p12 --store-password "$PASSWD" --cert-only > nuocmd.cert
```

The newly created keystore files can be copied on the client machine and then on some AP pod.

```bash
kubectl cp generate-nuodb-certs:/tmp/keys/. /tmp/keys

kubectl cp /tmp/keys/nuocmd.cert nuodb/admin-nuodb-cluster0-0:/tmp/nuocmd.cert
```

Add the new certificate to the truststore of all admin and database processes:

```bash
kubectl exec -ti admin-nuodb-cluster0-0 \
  nuocmd add trusted-certificate \
    --alias nuocmd_prime --cert /tmp/nuocmd.cert --timeout 10
```

This adds the new client certificate under the alias `nuocmd_prime`.
It is not possible to replace the client certificate because all processes in the domain are currently using it.
This single command invocation causes the trusted certificate to be propagated to all APs and all database processes.
The `--timeout` argument specifies the amount of time to wait for the new certificate to be propagated to the truststore of every process in the domain.

Update the client key as described in section [Update Key Pair Certificates](#update-key-pair-certificates).

### Rotate CA Certificate

Generate new CA certificate:

```bash
kubectl exec -ti generate-nuodb-certs -- \
  nuocmd create keypair \
    --keystore /tmp/keys.ca.p12 --store-password "$PASSWD" \
    --dname "CN=ca.nuodb.com, OU=Eng, O=NuoDB, L=Boston, ST=MA, C=US, SERIALNUMBER=67890" \
    --validity 36500 --ca

kubectl exec -ti generate-nuodb-certs -- \
  nuocmd show certificate \
    --keystore /tmp/keys.ca.p12 --store-password "$PASSWD" --cert-only > ca.cert
```

Generate a new server key pair:

```bash
kubectl exec -ti generate-nuodb-certs -- \
  nuocmd create keypair \
    --keystore /tmp/keys/nuoadmin.p12 --store-password "$PASSWD" --ca \
    --dname "CN=*.nuodb.svc.cluster.local, OU=Eng, O=NuoDB, L=Boston, ST=MA, C=US, SERIALNUMBER=67890"
```

Generate and sign the new server certificates using the new CA certificate chain:

```bash
kubectl exec -ti generate-nuodb-certs -- \
  nuocmd sign certificate \
    --keystore /tmp/keys/nuoadmin.p12 --store-password "$PASSWD" \
    --ca-keystore /tmp/keys/ca.p12 --ca-store-password "$PASSWD" \
    --validity 36500 --update
```

The newly created keystore files can be copied on the client machine and then on some AP pod.

```bash
kubectl cp generate-nuodb-certs:/tmp/keys/. /tmp/keys

kubectl cp /tmp/keys/ca.cert nuodb/admin-nuodb-cluster0-0:/tmp/ca.cert
```

If the server certificates are signed by a public CA, then follow the provided steps to renew your certificates.

> **NOTE**: For information on how to create the server keystore when using certificates signed by a Public CA, see [here](https://doc.nuodb.com/nuodb/latest/deployment-models/physical-or-vmware-environments-with-nuodb-admin/domain-operations/enabling-tls-encryption/using-certificates-signed-by-a-public-certificate-authority/).

Add the new CA certificate to the truststore of all admin and database processes:

```bash
kubectl exec -ti admin-nuodb-cluster0-0 \
  nuocmd add trusted-certificate \
    --alias ca_prime --cert /tmp/ca.cert --timeout 10
```

This adds the new client certificate under the alias `ca_prime`.
It is not possible to replace the CA certificate because all processes in the domain and SQL clients are currently using it.
This single command invocation causes the trusted certificate to be propagated to all APs and all database processes.
The `--timeout` argument specifies the amount of time to wait for the new certificate to be propagated to the truststore of every process in the domain.

Update domain keystore with the renewed certificates as described in [Update Key Pair Certificates](#update-key-pair-certificates).

### Cleanup

Backup the `ca.p12` keystore file as it will be needed during the next keys rotation.

Remove the helper pod used to generate certificates and the local copy of the keystore files:

```bash
kubectl delete pod generate-nuodb-certs

rm -rf /tmp/keys
```

Make sure that no AP, domain process or SQL client uses the old key pair certificates as described in [Verify Domain Certificates](#verify-domain-certificates).
It's recommended to remove the old certificates from the domain for security reasons which involves:

- removing the old Kubernetes secrets
- removing already rotated old CA or client certificates from truststore
- removing already rotated old CA or client certificates from SQL clients and NuoDB Commands clients

To remove the old certificate from the truststore of every AP and every database process use the following command:

```bash
kubectl exec -ti admin-nuodb-cluster0-0 \
  nuocmd remove trusted-certificate --alias ca
```
