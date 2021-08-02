# Rotating NuoDB TLS Certificates

<details>
<summary>Table of Contents</summary>
<!-- TOC -->

- [Rotating NuoDB TLS Certificates](#rotating-nuodb-tls-certificates)
    - [Introduction](#introduction)
        - [Terminology](#terminology)
    - [About Key Rotation](#about-key-rotation)
        - [Update Key Pair Certificates](#update-key-pair-certificates)
            - [Using Kubernetes Secrets](#using-kubernetes-secrets)
            - [Using HashiCorp Vault](#using-hashicorp-vault)
        - [Verify Domain Certificates](#verify-domain-certificates)
    - [NuoDB TLS Key Rotation](#nuodb-tls-key-rotation)
        - [Rotate CA Certificate](#rotate-ca-certificate)
        - [Rotate Server Certificate](#rotate-server-certificate)
        - [Rotate Client Certificate](#rotate-client-certificate)
        - [Cleanup](#cleanup)

<!-- /TOC -->
</details>

## Introduction

When using TLS encryption, it is necessary to rotate key pair certificates before the certificates expire.
To avoid downtime, NuoDB Admin Processes (APs) with the new certificates should be introduced in a rolling fashion.

This document expands on [Security Model of NuoDB in Kubernetes](./HowtoTLS.md) and explains how to configure proper certificate rotation in NuoDB Kubernetes deployments.

> **NOTE**: For non-Kubernetes deployments, see  [Enabling TLS encryption, Rotating TLS Key Pair Certificates](https://doc.nuodb.com/nuodb/latest/deployment-models/physical-or-vmware-environments-with-nuodb-admin/domain-operations/enabling-tls-encryption/rotating-key-pair-certificates/).
This page expands on the product documentation and is specific to this Helm Chart repository.

### Terminology

- `Key` = a combination of a private key with its corresponding X509 certificate chain. 
These keys are usually saved in a PKCS12 file such as `nuoadmin.p12`.
- `NuoDB Admin` = The administrative interface for domain and database management.
For information on how to start the NuoDB Admin tier, see [Admin Chart](../stable/admin/README.md).
- `CA` = Certificate Authority

## About Key Rotation

NuoDB supports various key models for TLS keys used by the `NuoDB Admin` and database processes.

> **NOTE**: For information on how to configure different TLS key models, see [Security Model of NuoDB in Kubernetes](./HowtoTLS.md).

The [Shared Admin Key + Intermediate CA](./docs/HowtoTLS.md#intermediate-ca) model is used by default in NuoDB Kubernetes deployments.
This means that the same key pair certificate is used by all APs, a Certificate Authority (CA) is used to sign them and the CA certificate is trusted by all processes in the domain.

If a CA certificate needs to be renewed, then both the keystore and the truststore must be updated for all APs, and by extension, all database processes.
As stated, processes with the new certificates should be introduced in a rolling fashion to avoid downtime, which means that processes with the old key pair certificate can verify processes with the new key pair certificate.
Since the new key pair certificates must be verified using a new trusted certificate, existing processes (admin and database) must have the new trusted certificate propagated to their truststores before the new key pair certificates can be introduced.

To simplify updating and propagating new trusted certificates, run `nuocmd add trusted-certificate`.
This command adds a trusted certificate to the NuoDB Admin server and causes the trusted certificate to be added and propagated to all APs and all database processes.

An overview of how CA certificate rotation is performed is as follows:

1. A new CA certificate is generated.
2. A new shared admin new key pair certificate is generated.
3. _(Optional)_ A new client key pair certificate and PEM file are generated.
4. _(Only if step 3 is performed)_ The new client certificate is added to the truststore of every AP.
5. The new CA certificate is added to the truststore of every AP and every database process.
6. The new CA certificate is added to the truststore of every SQL client.
7. The new CA certificate is added to the truststore of every NuoDB Command client.
8. For each AP, the keystore is replaced. Make sure that the new key pair certificate is in effect.
9. For each database process, the process is restarted, so that its certificate chain is based on the new certificate.
10. _(Only if step 3 is performed)_ For each NuoDB Command client, the client PEM file is replaced. Make sure that the new key pair certificate is in effect.
11. _(Optional)_ The old CA certificate is removed from the truststore of every AP and every database process.
12. _(Optional)_ The old CA certificate is removed from the truststore of every SQL client and NuoDB Command client.
13. _(Optional)_ The old client certificate is removed from the truststore of every AP.

Usually, only the server key pair certificate is renewed, which means that only the keystore file has to be updated for all APs and database processes.
This reduces the steps needed during key rotation as the new key pair certificate can be verified using the old truststore certificate.

An overview of how server key pair certificate rotation is performed is as follows:

1. A new shared admin new key pair certificate is generated.
2. _(Optional)_ A new client key pair certificate and PEM file are generated.
3. For each AP, the keystore is replaced. Make sure that the new key pair certificate is in effect.
4. For each database process, the process is restarted, so that its certificate chain is based on the new certificate.
5. _(Only if step 3 is performed)_ For each NuoDB Command client, the client PEM file is replaced. Make sure that the new key pair certificate is in effect.
6. _(Only if step 3 is performed)_ The old client certificate is removed from the truststore of every AP.

Complete step by step examples can be found in [NuoDB TLS Key Rotation](#nuodb-tls-key-rotation) section.

### Update Key Pair Certificates

Depending on the used TLS keys management solution, updating the key pair certificates may be different.
The sections bellow describe in detail how this is performed in Kubernetes deployments installed with NuoDB Helm Charts.

#### Using Kubernetes Secrets

By default, NuoDB uses Kubernetes secrets to store TLS keys and expose them to NuoDB processes.
The TLS keys are mounted in AP containers as a `subPath` volume mount which doesn't allow automatic secret updates.
This means that the secrets should be updated using the Kubernetes controller rolling upgrade strategy triggered by the `helm upgrade` command.

To replace the keystore for all APs, use the steps below.

Create new Kubernetes secrets using the keystore file generated with the renewed server key pair certificates.

```bash
kubectl create secret generic nuodb-keystore-renewed \
  --namespace nuodb \
  --from-file=nuoadmin.p12=/tmp/keys/nuoadmin.p12 \
  --from-literal=password=${PASSWD}

kubectl create secret generic nuodb-ca-cert-renewed \
  --namespace nuodb \
  --from-file=ca.cert=/tmp/keys/ca.cert

kubectl create secret generic nuodb-client-pem-renewed \
  --namespace nuodb \
  --from-file=nuocmd.pem=/tmp/keys/nuocmd.pem
```

> **NOTE**: Skip the generation of the secrets for which the key pair certificates haven't been rotated.

Upgrade the Helm release installed with the [admin](../stable/admin) chart using the new TLS secrets.
This will perform NuoDB Admin statefulset rolling upgrade so make sure that you are having enough APs to prevent downtime.

```bash
helm upgrade admin nuodb/admin \
  --namespace nuodb \
  --set admin.tlsKeyStore.secret=nuodb-keystore-renewed \
  --set admin.tlsClientPEM.secret=nuodb-client-pem-renewed \
  --set admin.tlsCACert.secret=nuodb-ca-cert-renewed \
  -f values.yaml
```

Wait for the AP statefulset rollout to finish and ensure that the new key pair certificates are used by all APs as described in [Verify Domain Certificates](#verify-domain-certificates).

Upgrade the Helm release installed with the [database](../stable/database) chart using the new TLS secrets.
This task will perform a rolling upgrade on all database Storage Manager (SM) statefulsets and Transaction Engine (TE) deployments, therefore ensure that enough additional database processes are running to prevent downtime.

```bash
helm upgrade database nuodb/database \
  --namespace nuodb \
  --set admin.tlsKeyStore.secret=nuodb-keystore-renewed \
  --set admin.tlsClientPEM.secret=nuodb-client-pem-renewed \
  --set admin.tlsCACert.secret=nuodb-ca-cert-renewed \
  -f values.yaml
```

> **NOTE**: Adjust the values depending on the key pair certificates that have been rotated.

Wait for all database processes to restart, report `Ready` and ensure that the new key pair certificates are used by all domain processes as described in [Verify Domain Certificates](#verify-domain-certificates).

#### Using HashiCorp Vault

HashiCorp uses sidecar containers to inject keys into another container.
It automatically updates secret values stored in the Vault KV store to match the files in the volume mounted inside the container.
NuoDB APs will reload the keystore without the need for reconfiguration and restart.
Therefore, it is enough to update the keystore, client PEM and CA certificate in HashiCorp Vault.

Ensure that the keystore has the same password (as specified by `admin.tlsKeyStore.password` Helm option) so that the APs can reload the keystore without having to be reconfigured and restarted.

For more information on how to configure TLS with HashiCorp Vault, see [Using HashiCorp Vault for Management of TLS Certificates](./HowToHashiCorpVault.md).

Make sure that the new key pair certificates are used by all APs as described in [Verify Domain Certificates](#verify-domain-certificates).

Restart database processes so that the APs generate new engine certificates signed with the new server certificate chain in the [Intermediate CA](./HowtoTLS.md#intermediate-ca) model or the new shared key pair certificates are used in the [Pass-down](./HowtoTLS.md#pass-down) model.
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
  deployment te-database-nuodb-cluster0-demo
```

Alternatively, the database processes can be restarted manually by shutting down each process using `nuocmd shutdown process --start-id ...` command.
Kubernetes will automatically restart the engine container which will effectively restart the database process.

It is recommended to perform each rollout or database process restart sequentially after verifying the domain process certificates as described in [Verify Domain Certificates](#verify-domain-certificates).

### Verify Domain Certificates

To get the current certificate data in the NuoDB domain, use `nuocmd get certificate-info` command and make sure that the renewed key pair certificates are displayed in the corresponding sections.

```bash
kubectl exec -ti admin-nuodb-cluster0-0 \
  --namespace nuodb -- \
  nuocmd --show-json get certificate-info
```

The command displays the certificates information in the below sections:

- `processCertificates` - information about certificates used by each database process identified by `startId`.
- `processTrusted` - trusted certificate aliases (referenced in `trustedCertificates` section) for each database process identified by `startId`.
- `serverCertificates` - information about certificates used by each AP identified by `serverId`.
- `serverTrusted` - trusted certificate aliases (referenced in `trustedCertificates` section) for each AP identified by `serverId`.
- `trustedCertificates` - information about trusted certificates in the domain identified by `alias`.

Example output:

```json
{
  "processCertificates": {
    "11": {
      "caPathLength": -1,
      "certificatePem": "-----BEGIN CERTIFICATE-----...-----END CERTIFICATE-----",
      "expires": 1659087997000,
      "expiresTimestamp": "2022-07-29T09:46:37.000+0000",
      "issuerName": "CN=*.nuodb.svc.cluster.local, OU=Eng, O=NuoDB, L=Boston, ST=MA, C=US, SERIALNUMBER=67890",
      "subjectName": "CN=172.17.0.8"
    },
    "12": {
      "caPathLength": -1,
      "certificatePem": "-----BEGIN CERTIFICATE-----...-----END CERTIFICATE-----",
      "expires": 1659087998000,
      "expiresTimestamp": "2022-07-29T09:46:38.000+0000",
      "issuerName": "CN=*.nuodb.svc.cluster.local, OU=Eng, O=NuoDB, L=Boston, ST=MA, C=US, SERIALNUMBER=67890",
      "subjectName": "CN=172.17.0.9"
    },
    "9": {
      "caPathLength": -1,
      "certificatePem": "-----BEGIN CERTIFICATE-----...-----END CERTIFICATE-----",
      "expires": 1659087974000,
      "expiresTimestamp": "2022-07-29T09:46:14.000+0000",
      "issuerName": "CN=*.nuodb.svc.cluster.local, OU=Eng, O=NuoDB, L=Boston, ST=MA, C=US, SERIALNUMBER=67890",
      "subjectName": "CN=172.17.0.5"
    }
  },
  "processTrusted": {
    "11": [
      "ca_prime",
      "nuocmd_prime"
    ],
    "12": [
      "ca_prime",
      "nuocmd_prime"
    ],
    "9": [
      "ca_prime",
      "nuocmd_prime"
    ]
  },
  "serverCertificates": {
    "admin-nuodb-cluster0-0": {
      "caPathLength": 2147483647,
      "certificatePem": "-----BEGIN CERTIFICATE-----...-----END CERTIFICATE-----",
      "expires": 4781150383000,
      "expiresTimestamp": "2121-07-05T09:19:43.000+0000",
      "issuerName": "CN=ca.nuodb.com, OU=Eng, O=NuoDB, L=Boston, ST=MA, C=US, SERIALNUMBER=67890",
      "subjectName": "CN=*.nuodb.svc.cluster.local, OU=Eng, O=NuoDB, L=Boston, ST=MA, C=US, SERIALNUMBER=67890"
    },
    "admin-nuodb-cluster0-1": {
      "caPathLength": 2147483647,
      "certificatePem": "-----BEGIN CERTIFICATE-----...-----END CERTIFICATE-----",
      "expires": 4781150383000,
      "expiresTimestamp": "2121-07-05T09:19:43.000+0000",
      "issuerName": "CN=ca.nuodb.com, OU=Eng, O=NuoDB, L=Boston, ST=MA, C=US, SERIALNUMBER=67890",
      "subjectName": "CN=*.nuodb.svc.cluster.local, OU=Eng, O=NuoDB, L=Boston, ST=MA, C=US, SERIALNUMBER=67890"
    }
  },
  "serverTrusted": {
    "admin-nuodb-cluster0-0": [
      "ca_prime",
      "nuocmd_prime"
    ],
    "admin-nuodb-cluster0-1": [
      "ca_prime",
      "nuocmd_prime"
    ]
  },
  "trustedCertificates": {
    "ca_prime": {
      "caPathLength": 2147483647,
      "certificatePem": "-----BEGIN CERTIFICATE-----...-----END CERTIFICATE-----",
      "expires": 4781150383000,
      "expiresTimestamp": "2121-07-05T09:19:43.000+0000",
      "issuerName": "CN=ca.nuodb.com, OU=Eng, O=NuoDB, L=Boston, ST=MA, C=US, SERIALNUMBER=67890",
      "subjectName": "CN=ca.nuodb.com, OU=Eng, O=NuoDB, L=Boston, ST=MA, C=US, SERIALNUMBER=67890"
    },
    "nuocmd_prime": {
      "caPathLength": -1,
      "certificatePem": "-----BEGIN CERTIFICATE-----...-----END CERTIFICATE-----",
      "expires": 4781151591000,
      "expiresTimestamp": "2121-07-05T09:39:51.000+0000",
      "issuerName": "CN=nuocmd.nuodb.com, OU=Eng, O=NuoDB, L=Boston, ST=MA, C=US, SERIALNUMBER=67890",
      "subjectName": "CN=nuocmd.nuodb.com, OU=Eng, O=NuoDB, L=Boston, ST=MA, C=US, SERIALNUMBER=67890"
    }
  }
}
```

## NuoDB TLS Key Rotation

To renew the domain key pair certificates and the server keystore, you can either create new keystore files and certificates on your own or create them using NuoDB commands.
All examples below will be using a helper pod started with NuoDB image and `nuocmd` command line tools.

```bash
PASSWD=changeMe
mkdir /tmp/keys

kubectl run generate-nuodb-certs \
  --image nuodb/nuodb-ce \
  --env="PASSWD=changeMe" \
  --command -- 'tail' '-f' '/dev/null'

kubectl exec -ti generate-nuodb-certs -- \
  mkdir -p /tmp/keys
```

### Rotate CA Certificate

Generate new CA certificate:

```bash
kubectl exec -ti generate-nuodb-certs -- \
  nuocmd create keypair \
    --keystore /tmp/keys/ca.p12 --store-password "$PASSWD" \
    --dname "CN=ca.nuodb.com, OU=Eng, O=NuoDB, L=Boston, ST=MA, C=US, SERIALNUMBER=67890" \
    --validity 36500 --ca

kubectl exec -ti generate-nuodb-certs -- bash -c \
  'nuocmd show certificate \
    --keystore /tmp/keys/ca.p12 --store-password "$PASSWD" \
    --cert-only > /tmp/keys/ca.cert'
```

If the CA certificate is owned by a public CA, then follow the official steps to obtain the new certificate.

The newly created keystore files can be copied on the client machine and then on some AP pod.

```bash
kubectl cp generate-nuodb-certs:/tmp/keys/. /tmp/keys

kubectl cp /tmp/keys/ca.cert nuodb/admin-nuodb-cluster0-0:/tmp/ca.cert
```

Add the new CA certificate to the truststore of all admin and database processes:

```bash
kubectl exec -ti admin-nuodb-cluster0-0 -- \
  nuocmd add trusted-certificate \
    --alias ca_prime --cert /tmp/ca.cert --timeout 30
```

The `--timeout` argument specifies the amount of time to wait for the new certificate to be propagated to the truststore of every process in the domain.

Ensure that the new CA certificate with alias `ca_prime` is trusted by all APs and database processes as described in [Verify Domain Certificates](#verify-domain-certificates).

Add the new CA certificate to the truststore of every SQL client.

Generate and sign the new server certificates using the new CA certificate chain as described in [Rotate Server Certificate](#rotate-server-certificate).

### Rotate Server Certificate

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
In case the CA certificate has been renewed right now, use the new CA keystore generated using the [Rotate CA Certificate](#rotate-ca-certificate) steps.

Generate and sign the new server certificates using the CA certificate chain:

```bash
kubectl cp ca.p12 generate-nuodb-certs:/tmp/keys/

kubectl exec -ti generate-nuodb-certs -- \
  nuocmd sign certificate \
    --keystore /tmp/keys/nuoadmin.p12 --store-password "$PASSWD" \
    --ca-keystore /tmp/keys/ca.p12 --ca-store-password "$PASSWD" \
    --validity 36500 --ca --update
```

The `--ca` argument specifies whether the `isCA` extension is enabled in the generated certificate.
This allows the admin to act as an intermediate CA and is used in the [Intermediate CA](./HowtoTLS.md#intermediate-ca) model.
You can omit this argument in case the [Pass-down](./HowtoTLS.md#pass-down) model is used.

The newly created keystore files can be copied on the client machine which will be used later to update the keystore secret.

```bash
kubectl cp generate-nuodb-certs:/tmp/keys/. /tmp/keys
```

If the server certificates are signed by a public CA, then follow the official steps to renew your certificates.

> **NOTE**: For information on how to create the server keystore when using certificates signed by a Public CA, see [here](https://doc.nuodb.com/nuodb/latest/deployment-models/physical-or-vmware-environments-with-nuodb-admin/domain-operations/enabling-tls-encryption/using-certificates-signed-by-a-public-certificate-authority/).

Update domain keystore with the renewed certificates as described in section [Update Key Pair Certificates](#update-key-pair-certificates).

### Rotate Client Certificate

The NuoDB Command tool client key can be renewed together with the NuoDB Admin server key pair certificates or separately as needed.

Generate new client key pair:

```bash
kubectl exec -ti generate-nuodb-certs -- \
  nuocmd create keypair \
    --keystore /tmp/keys/nuocmd.p12 --store-password "$PASSWD" \
    --dname "CN=nuocmd.nuodb.com, OU=Eng, O=NuoDB, L=Boston, ST=MA, C=US, SERIALNUMBER=67890" \
    --validity 36500

kubectl exec -ti generate-nuodb-certs -- bash -c \
  'nuocmd show certificate \
    --keystore /tmp/keys/nuocmd.p12 --store-password "$PASSWD" > /tmp/keys/nuocmd.pem'

kubectl exec -ti generate-nuodb-certs -- bash -c \
  'nuocmd show certificate \
    --keystore /tmp/keys/nuocmd.p12 --store-password "$PASSWD" \
    --cert-only > /tmp/keys/nuocmd.cert'
```

The newly created keystore files can be copied on the client machine and then to an AP pod.

```bash
kubectl cp generate-nuodb-certs:/tmp/keys/. /tmp/keys

kubectl cp /tmp/keys/nuocmd.cert nuodb/admin-nuodb-cluster0-0:/tmp/nuocmd.cert
```

Add the new certificate to the truststore of all admin and database processes:

```bash
kubectl exec -ti admin-nuodb-cluster0-0 -- \
  nuocmd add trusted-certificate \
    --alias nuocmd_prime --cert /tmp/nuocmd.cert --timeout 30
```

The `--timeout` argument specifies the amount of time to wait for the new certificate to be propagated to the truststore of every process in the domain.

Esure that the new client certificate with alias `nuocmd_prime` is trusted by all APs and database processes as described in [Verify Domain Certificates](#verify-domain-certificates).

Update the client key as described in section [Update Key Pair Certificates](#update-key-pair-certificates).

### Cleanup

Backup the `ca.p12` keystore file as it will be needed during the next keys rotation.

Remove the helper pod used to generate certificates and the local copy of the keystore files:

```bash
kubectl delete pod generate-nuodb-certs

unset PASSWD

rm -rf /tmp/keys
```

Ensure that AP, domain process, and SQL client do not use the old key pair certificates as described in [Verify Domain Certificates](#verify-domain-certificates).
It's recommended to remove the old certificates from the domain for security reasons which involves:

- removing the old Kubernetes secrets
- removing already rotated old CA or client certificates from truststore
- removing already rotated old CA or client certificates from SQL clients and NuoDB Commands clients

To remove the old certificate from the truststore of every AP and every database process use the following command:

```bash
kubectl exec -ti admin-nuodb-cluster0-0 -- \
  nuocmd remove trusted-certificate --alias <alias>
```
