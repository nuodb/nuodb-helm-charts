# NuoDB Transparent Data Encryption in Kubernetes

## Introduction

NuoDB offers strong protection for both _data at rest_ and _data in transit_ allowing enterprises to implement and enforce needed company security policies. This helps you to protect your valuable company data. Transport Layer Security (TLS) secures data in transit within NuoDB domain by encrypting all network communication and the data which is in motion. Transparent Data Encryption (TDE) allows data at rest to be encrypted before being written to disk.

Before TDE encryption can be enabled on a database, `NuoDB Admin` must be configured with a _storage password_ for each database.
`NuoDB Admin` is responsible for propagating storage passwords within the admin layer and to database engines.
All SMs and TEs in the database will use this storage password as a key to secure the encrypted data.
Starting from NuoDB v4.1.2 and NuoDB Helm Charts v3.1.0 Transparent Data Encryption is supported in NuoDB Kubernetes deployments.

Enabling TDE consists of two steps
1. [Supply storage passwords](#supplying-storage-passwords)
2. [Configure database layer](#enable-transparent-data-encryption)

> **NOTE**: For information on enabling TDE in non-Kubernetes deployments of NuoDB, see [here](https://doc.nuodb.com/nuodb/latest/database-administration/transparent-data-encryption/configuring-transparent-data-encryption/). 
This document expands on the product documentation and is specific to this Helm Chart repository.

### Terminology

- `NuoDB Admin` = admin interface for domain and database management. Started by the [Admin Chart](../stable/admin/README.md).
- `Storage password` = key used by database engines to secure the encrypted data.
- `Target password` = Password that needs to be set as _current_ database storage password.
- `Historical password`= Password used to encrypt _any_ archive in the database or a backupset. 
There could be multiple historical passwords.
- `TDE monitor service` = NuoDB Admin auxiliary service which can read storage passwords configuration from files on disk.


## Supplying storage passwords 

The domain administrator is responsible for configuring TDE storage passwords. They can be supplied as Kubernetes secret or injected inside the pods using 3rd party software like HashiCorp Vault.
If using third-party party software, the storage passwords need to be injected and stored as files inside AP and SM pods. NuoDB supports file per password approach explained in [using Kubernetes secrets](#using-kubernetes-secrets) and single configuration file per database approach explained in [HashiCorp Vault integration](#using-hashicorp-vault).
Choosing the right approach depends on the capabilities and flexibility that any other 3rd party tool securing the storage passwords have.

> **IMPORTANT**: NuoDB, Inc. cannot recover Transparent Data Encryption storage passwords, and the database cannot be loaded if its storage password is lost. It is the user’s responsibility to keep track of storage passwords, including historical passwords for encrypted hot copy backups.

### Using Kubernetes secrets

Storage passwords can be supplied using Kubernetes secret.
One secret per database should be created in the Kubernetes namespace where NuoDB is deployed.
It must hold the `target` and optional `historical` storage passwords for that database.
The configured secret name is mounted inside all AP, hotcopy SM and non-hotcopy SM pods. The mount path is `/etc/nuodb/tde/<db-name>` by default and the static portion of the path can be changed using `admin.tde.storagePasswordsDir` variable.
NuoDB Admin monitors all directories and files under `/etc/nuodb/tde` path and will configure storage passwords available in the mounted files. The following files can be used:

- `/etc/nuodb/tde/<db-name>/target` = Holds the database target storage password.
- `/etc/nuodb/tde/<db-name>/historical*` _(optional)_ = Holds a database historical storage password. Multiple files can be used to supply multiple historical passwords for each database. Historical passwords are discussed in detailin the  [Change storage passwords](#change-storage-passwords) section.

Create a Kubernetes secret holding TDE passwords:

```bash
kubectl create secret generic demo-tde-secret -n nuodb \
  --from-literal target='topSecret'
```

The secret name need to be set in the [Admin chart](../stable/admin/README.md) and [Database chart](../stable/database/README.md) during installation:

```bash
helm install -n nuodb admin stable/admin \
  --set admin.tde.secrets.demo=demo-tde-secret

helm install -n nuodb database stable/database \
  --set admin.tde.secrets.demo=demo-tde-secret
```

### Using HashiCorp Vault

[Vault Agent Sidecar Injector](https://www.vaultproject.io/docs/platform/k8s/injector) can be used to inject Vault secrets as files inside Kubernetes pods. 

The scope of this document doesn't include steps on how to install and setup HashiCorp Vault. 

> **NOTE**: For a quick tutorial on how to deploy Vault in Kubernetes using agent injector, see [here](https://www.hashicorp.com/blog/injecting-vault-secrets-into-kubernetes-pods-via-a-sidecar).

The secret injection is controlled via Kubernetes annotations which can be specified for [Admin chart](../stable/admin/README.md) or [Database chart](../stable/database/README.md) using respectively `admin.podAnnotations` and `database.podAnnotations` variables.

In addition to single file per password explained in [Using Kubernetes secrets](#using-kubernetes-secrets) section, NuoDB Admin supports a single configuration file per database `/etc/nuodb/tde/<db-name>/tde.json` which holds all storage passwords for a single NuoDB database.

The file is encoded in JSON format and takes precedence over the file per password approach. 
Example file content:

```json
{
  "historical-20201203":"secret",
  "historical-20201204":"superSecret",
  "target":"topSecret"
}
```

Use a single secret per database stored in Vault KV store.

```bash
vault kv put v2.nuodb.com/TDE/demo \
  target=topSecret \
  historical-20201204=superSecret \
  historical-20201203=secret

```
The secret name can be different from the example above.
In this case _v2_ illustrates that Vault KV store version 2 is used and _v2.nuodb.com/TDE_ path can group several Vault secrets for several NuoDB databases.

Vault policy should be created and attached to Vault Kubernetes role to allow `read` capability for the secret created above by NuoDB service account. In the examples further in this section, `nuodb` Vault Kubernetes role is used.

The `tde.json` configuration file can be injected using the following Helm values file snippet:

```yaml
...

admin:
  podAnnotations:
    vault.hashicorp.com/agent-inject: "true"
    vault.hashicorp.com/role: nuodb
    vault.hashicorp.com/agent-inject-file-demo: tde.json
    vault.hashicorp.com/secret-volume-path-demo: /etc/nuodb/tde/demo
    vault.hashicorp.com/agent-inject-secret-demo: v2.nuodb.com/TDE/demo
    vault.hashicorp.com/agent-inject-template-demo: |
      {{- with secret "v2.nuodb.com/TDE/demo" -}}
      {{- .Data.data | toJSON }}
      {{- end -}}
 
database:
  podAnnotations:
    vault.hashicorp.com/agent-inject: "true"
    vault.hashicorp.com/role: nuodb
    vault.hashicorp.com/agent-inject-file-demo: tde.json
    vault.hashicorp.com/secret-volume-path-demo: /etc/nuodb/tde/demo
    vault.hashicorp.com/agent-inject-secret-demo: v2.nuodb.com/TDE/demo
    vault.hashicorp.com/agent-inject-template-demo: |
      {{- with secret "v2.nuodb.com/TDE/demo" -}}
      {{- .Data.data | toJSON }}
      {{- end -}}

...
```

Templates for Vault KV store version 1 and version 2 are different..
If Vault KV store version 1 is used, `.Data.data` references in inject templates should be replaced with `.Data`.

> **NOTE**: For information on how different annotations can be used, see [here](https://www.vaultproject.io/docs/platform/k8s/injector/annotations).

If `tde.conf` needs to be created using several secrets or special keys from the secret needs to be filtered out, [scratch](https://github.com/hashicorp/consul-template#scratch) Consul template API functions can be used. 
For example this template will filter out any other keys in a Vault secret:

```yaml
    vault.hashicorp.com/agent-inject-template-demo: |
      {{- with secret "v2.nuodb.com/TDE/demo" -}}
      {{- range $k, $v := .Data.data -}}
      {{- if or (eq $k "target") ($k | regexMatch "historical.*") }}
      {{- scratch.MapSet "passwords" $k $v }}
      {{- end }}
      {{- end }}
      {{- scratch.Get "passwords" | toJSON }}
      {{- end -}}
```


## Enable Transparent Data Encryption

NuoDB Admin will configure supplied storage passwords and will propagate them to engines.
To verify that the target password is configured correctly, the following command should succeed returning no errors.

```bash
kubectl exec admin-nuodb-cluster0-0 -- nuocmd check data-encryption \
  --db-name demo \
  --password 'topSecret'
```

Now that the database has a storage password, encryption can be enabled using `ALTER DATABASE CHANGE ENCRYPTION TYPE`. Database administrator privileges are required in order to use this command.
The SQL statement should be executed manually by obtaining SQL connection towards the target database.

```bash
kubectl exec -n nuodb admin-nuodb-cluster0-0 -- bash -c \
  'echo "alter database change encryption type AES128;" | /opt/nuodb/bin/nuosql demo --user dba --password secret'
```

SMs will start data encryption in the background.

> **NOTE**: For information on how to confirm if TDE is enabled, see [here](https://doc.nuodb.com/nuodb/latest/database-administration/transparent-data-encryption/configuring-transparent-data-encryption/#_confirming_that_tde_is_enabled).


## Change storage passwords

Database storage passwords can be changed at any time by changing the secret supplying them. The password rotation will be performed automatically in the background.

The current storage password must be specified as historical so that it can be used if:

- some of the SMs are down during password rotation
- database restore needs to be performed from a backupset encrypted with the previous password

If Kubernetes secret is used to supply the storage passwords, update it to match the desired state.

```bash
kubectl create secret generic demo-tde-secret -n nuodb \
  --from-literal target='superSecret' \
  --from-literal historical-20201215='topSecret' \
  --dry-run=client -o yaml | kubectl apply -f -
```

Kubernetes will update automatically mounted secrets inside the containers. The delay depends on _kubelet sync period_ + _kubelet cache propagation delay_.

> **NOTE**: For information on Kubernetes secrets automatic update, see [here](https://kubernetes.io/docs/concepts/configuration/secret/#mounted-secrets-are-updated-automatically).

If any other 3rd party software is used to inject the storage passwords configuration, update the passwords and ensure that the configuration is updated inside the NuoDB pods.

> **NOTE**: For information on HashiCorp Vault secrets renewals, see [here](https://www.vaultproject.io/docs/agent/template#renewals-and-updating-secrets).

To verify that password rotation is done successfully, the following command should succeed returning no errors.

```bash
kubectl exec admin-nuodb-cluster0-0 -- nuocmd check data-encryption \
  --db-name demo \
  --password 'superSecret' \
  --timeout 90
```

Database archives do not need to be re-encrypted when the storage password is changed.

## Working with encrypted backup

Before importing or restoring from an encrypted backupset or cold archive, the storage password used to encrypt it must be present in the NuoDB Admin either as target or historical database storage password.
All configured database storage passwords will be provided to `nuoarchive restore` when restoring from a backupset.
NuoDB Helm Charts 3.1.0+ handles encrypted backup restore automatically.
