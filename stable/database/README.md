# NuoDB Database Helm Chart

This chart starts a NuoDB database deployment on a Kubernetes cluster using the Helm package manager. To start a second database under the same NuoDB Admin deployment, start a second database using the same instructions with a new database name. 

## Command

```bash
helm install [name] nuodb/database [--generate-name] [--set parameter] [--values myvalues.yaml]
```

## Software Version Prerequisites

Please visit the **[NuoDB Helm Chart main page](https://github.com/nuodb/nuodb-helm-charts/#software-release-requirements)** for software version prerequisites.

## Installing the Chart

All configurable parameters for each top-level scope are detailed below, organized by scope.

#### global.*

The purpose of this section is to specify global settings.

The following tables list the configurable parameters for the `global` option:

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `global.imageRegistry` | Global Docker image registry | `nil` |
| `global.imagePullSecrets` | Global Docker registry secret names as an array | `[]` (does not add image pull secrets to deployed pods) |

#### cloud.*

The purpose of this section is to specify the cloud provider, and specify the availability zones where a solution is deployed.

The following tables list the configurable parameters for the `cloud` option:

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `provider` | Cloud provider; permissible values include: `azure`, `amazon`, or `google` |`nil`|
| `zones` | List of availability zones to deploy to |`[]`|
| `cluster.name` | logical name of the cluster. Useful in multi-cluster configs | `cluster0` |

For example, for the Google Cloud:

```yaml
cloud:
  provider: google
  zones:
    - us-central1-a
    - us-central1-b
    - us-central1-c
  cluster:
    name: cluster0
```

#### busybox.*

The purpose of this section is to specify the Busybox image parameters.

The following tables list the configurable parameters for the `busybox` option:

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `image.registry` | busybox container registry | `docker.io` |
| `image.repository` | busybox container image name |`busybox`|
| `image.tag` | busybox container image tag | `latest` |
| `image.pullPolicy` | busybox container pull policy |`IfNotPresent`|
| `image.pullSecrets` | Specify docker-registry secret names as an array | [] (does not add image pull secrets to deployed pods) |

The `registry` option can be used to connect to private image repositories, such as Artifactory.

The `registry`, `repository`, and `tag` values are combined to form the `image` declaration in the Helm charts.

For example, when using GlusterFS storage class, you would supply the following parameter:

```bash
  ...
  --set busybox.image.registry=acme-dockerv2-virtual.jfrog.io
  ...
```

For example:

```yaml
busybox:
  image:
    registry: docker.io
    repository: busybox
    tag: latest
    pullPolicy: IfNotPresent
```

#### nuodb.*

The purpose of this section is to specify the NuoDB image parameters.

The following tables list the configurable parameters for the `nuodb` option:

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `image.registry` | NuoDB container registry | `docker.io` |
| `image.repository` | NuoDB container image name |`nuodb/nuodb-ce`|
| `image.tag` | NuoDB container image tag | `4.0.8` |
| `image.pullPolicy` | NuoDB container pull policy |`IfNotPresent`|
| `image.pullSecrets` | Specify docker-registry secret names as an array | [] (does not add image pull secrets to deployed pods) |
| `serviceAccount` | The name of the service account used by NuoDB Pods | `nuodb` |
| `addRoleBinding` | Whether to add role and role-binding giving `serviceAccount` access to Kubernetes APIs (Pods, PersistentVolumes, PersistentVolumeClaims, StatefulSets) | `true` |

The `registry` option can be used to connect to private image repositories, such as Artifactory.

The `registry`, `repository`, and `tag` values are combined to form the `image` declaration in the Helm charts.

For example, when using GlusterFS storage class, you would supply the following parameter:

```bash
  ...
  --set nuodb.image.registry=acme-dockerv2-virtual.jfrog.io
  ...
```

For example:

```yaml
nuodb:
  image:
    registry: docker.io
    repository: nuodb/nuodb-ce
    tag: latest
    pullPolicy: IfNotPresent
```

#### admin.*

The purpose of this section is to specify the NuoDB Admin parameters.

The following tables list the configurable admin parameters for a database and their default values.

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `domain` | NuoDB admin cluster name | `nuodb` |
| `namespace` | Namespace where admin is deployed; when peering to an existing admin cluster provide its project name | `nuodb` |
| `configFilesPath` | Directory path where `configFiles.*` are found | `/etc/nuodb/` |
| `configFiles.*` | See below. | `{}` |
| `tlsCACert.secret` | TLS CA certificate secret name | `nil` |
| `tlsCACert.key` | TLS CA certificate secret key | `nil` |
| `tlsKeyStore.secret` | TLS keystore secret name | `nil` |
| `tlsKeyStore.key` | TLS keystore secret key | `nil` |
| `tlsKeyStore.password` | TLS keystore secret password | `nil` |
| `tlsClientPEM.secret` | TLS client PEM secret name | `nil` |
| `tlsClientPEM.key` | TLS client PEM secret key | `nil` |
| `tde.secrets` | Transparent Data Encryption secret names used for different databases | `{}` |
| `tde.storagePasswordsDir` | Transparent Data Encryption storage passwords mount path | `/etc/nuodb/tde` |

#### admin.configFiles.*

The purpose of this section is to detail how to provide alternative configuration files for NuoDB. 

There are two sets of configuration files documented:

- [Admin Configuration for a Particular Deployment][1]
- [Database Configuration for a Particular Deployment][2]

Any file located in `admin.configFilesPath` can be replaced; the YAML key corresponds to the file name being created or replaced.

| Key | Description | Default |
| ----- | ----------- | ------ |
| `nuodb.lic` | [NuoDB license file content; defaults to NuoDB CE Edition][3] | `nil` |
| `nuoadmin.conf` | [NuoDB Admin host properties][4] | `nil` |
| `nuodb-types.config`| [Type mappings for the NuoDB Migrator tool][5] | `nil` |
| `nuoadmin.logback.xml` | Logging configuration. NuoDB recommends using the default settings. | `nil` |

#### admin.serviceSuffix

The purpose of this section is to allow customisation of the names of the clusterIP and balancer admin services (load-balancers).

| Key | Description | Default |
| ----- | ----------- | ------ |
| `clusterip` | suffix for the clusterIP load-balancer | "clusterip" |
| `balancer` | suffix for the balancer service | "balancer" |

#### database.*

The following tables list the configurable parameters of the `database` chart and their default values.

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `nameOverride` | Provide a name in place of `database-dbname` |`""`|
| `fullNameOverride` | Provide a name to substitute for the full names of resources |`""`|
| `name` | Database name | `demo` |
| `rootUser` | Database root user | `dba` |
| `rootPassword` | Database root password | `secret` |
| `securityContext.enabled` | Enable security context | `false` |
| `securityContext.runAsUser` | User ID for the container | `1000` |
| `securityContext.fsGroup` | Group ID for the container | `1000` |
| `securityContext.capabilities` | Enable capabilities for the container - disregards `securityContext.enabled` | `[]` |
| `env` | Import ENV vars inside containers | `[]` |
| `envFrom` | Import ENV vars from configMap(s) | `[]` |
| `lbConfig.prefilter` | Database load balancer prefilter expression | `nil` |
| `lbConfig.default` | Database load balancer default query | `nil` |
| `persistence.accessModes` | Volume access modes enabled (must match capabilities of the storage class) | `ReadWriteOnce` |
| `persistence.size` | Amount of disk space allocated for database archive storage | `20Gi` |
| `persistence.storageClass` | Storage class for volume backing database archive storage | `-` |
| `configFilesPath` | Directory path where `configFiles.*` are found | `/etc/nuodb/` |
| `configFiles.*` | See below. | `{}` |
| `podAnnotations` | Annotations to pass through to the SM an TE pods | `nil` |
| `primaryRelease` | One primary and multiple secondary database Helm releases for the same database name may be deployed into the same Kubernetes namespace. Set to `false` when deploying secondary database Helm releases. | `true` |
| `autoImport.*` | Enable and configure the automatic initializing of the initial database state | `disabled` |
| `autoImport.source` | The source - typically a URL - of the database copy to import | `""` |
| `autoImport.credentials` | Authentication for the download of `source` in the form `user`:`password` | '""'|
| `autoImport.stripLevels` | The number of levels to strip off pathnames when unpacking a TAR file of an archive | `1` |
| `autoImport.type` | Type of content in `source`. One of `stream` -> exact copy of an archive; or `backupset` -> a NuoDB hotcopy backupset | 'backupset' |
| `autoRestore.*` | Enable and configure the automatic re-initialization of a single archive in a running database - see the options in `autoImport` | `disabled` |
| `sm.logPersistence.enabled` | Whether to enable persistent storage for logs | `false` |
| `sm.logPersistence.overwriteBackoff.copies` | How many copies of the crash directory to keep within windowMinutes | `3` |
| `sm.logPersistence.overwriteBackoff.windowMinutes` | The window within which to keep the number of crash copies | `120` |
| `sm.logPersistence.accessModes` | Volume access modes enabled (must match capabilities of the storage class) | `ReadWriteOnce` |
| `sm.logPersistence.size` | Amount of disk space allocated for log storage | `60Gi` |
| `sm.logPersistence.storageClass` | Storage class for volume backing log storage.  This storage class must be pre-configured in the cluster | `-` |
| `sm.hotCopy.replicas` | SM replicas with hot-copy enabled | `1` |
| `sm.hotCopy.enablePod` | Create DS/SS for hot-copy enabled SMs | `true` |
| `sm.hotCopy.enableBackups` | Enable full/incremental/journal backups | `true` |
| `sm.hotCopy.deadline` | Deadline for a hotcopy job to start (seconds) | `1800` |
| `sm.hotCopy.timeout` | Timeout for a started `full` or `incremental` hotcopy job to complete (seconds). The default timeout of "0" will force the backup jobs to wait forever for the requested hotcopy operation to complete | `0` |
| `sm.hotCopy.successHistory` | Number of successful Jobs to keep | `5` |
| `sm.hotCopy.failureHistory` | Number of failed jobs to keep | `5` |
| `sm.hotCopy.backupDir` | Directory path where backupsets will be stored | `/var/opt/nuodb/backup` |
| `sm.hotCopy.backupGroup` | Name of the backup group | `{{ .Values.cloud.cluster.name }}` |
| `sm.hotCopy.fullSchedule` | cron schedule for FULL hotcopy jobs | `35 22 * * 6` |
| `sm.hotCopy.incrementalSchedule` | cron schedule for INCREMENTAL hotcopy jobs | `35 22 * * 0-5` |
| `sm.hotCopy.restartPolicy` | Restart policy for backup related JOBs and CRON JOBs | `OnFailure` |
| `sm.hotCopy.persistence.size` | size of the hotcopy storage PV | `20Gi` |
| `sm.hotCopy.persistence.accessModes` | access modes for the hotcopy storage PV | `[ ReadWriteOnce ]` |
| `sm.hotCopy.persistence.size` | size of the hotcopy storage PV | `20Gi` |
| `sm.hotCopy.journalBackup.enabled` | Is `journal hotcopy` enabled - true/false | `false` |
| `sm.hotCopy.journalBackup.intervalMinutes` | Frequency of running `journal hotcopy` (minutes) | `15` |
| `sm.hotCopy.journalBackup.deadline` | Deadline for a `journal hotcopy` job to start (seconds) | `90` |
| `sm.hotCopy.journalBackup.timeout` | Timeout for a started `journal hotcopy` to complete (seconds). The default timeout of "0" will force the backup jobs to wait forever for the requested hotcopy operation to complete | `0` |
| `sm.hotCopy.coldStorage.credentials` | Credentials for accessing backup cold storage (user:password) | `""` |
| `sm.hotCopy.journalPath.enabled` | Whether to enable separate SM journal directory. For more info, read the [Journal HowTo](../../docs/HowToArchiveJournal.md) | `false` |
| `sm.hotCopy.journalPath.accessModes` | Volume access modes enabled (must match capabilities of the storage class) | `ReadWriteOnce` |
| `sm.hotCopy.journalPath.size` | Amount of disk space allocated for SM journal | `20Gi` |
| `sm.hotCopy.journalPath.storageClass` | Storage class for SM journal.  This storage class must be pre-configured in the cluster | `-` |
| `sm.noHotCopy.replicas` | SM replicas with hot-copy disabled | `0` |
| `sm.noHotCopy.enablePod` | Create DS/SS for non-hot-copy SMs | `true` |
| `sm.noHotCopy.journalPath.enabled` | Whether to enable separate SM journal directory. For more info, read the [Journal HowTo](../../docs/HowToArchiveJournal.md) | `false` |
| `sm.noHotCopy.journalPath.accessModes` | Volume access modes enabled (must match capabilities of the storage class) | `ReadWriteOnce` |
| `sm.noHotCopy.journalPath.size` | Amount of disk space allocated for SM journal | `20Gi` |
| `sm.noHotCopy.journalPath.storageClass` | Storage class for SM journal.  This storage class must be pre-configured in the cluster | `-` |
| `sm.labels` | Labels given to the SMs started | `{}` |
| `sm.engineOptions` | Additional NuoDB engine options | `{}` |
| `sm.resources` | Labels to apply to all resources | `{}` |
| `sm.affinity` | Affinity rules for NuoDB SM | `{}` |
| `sm.nodeSelector` | Node selector rules for NuoDB SM | `{}` |
| `sm.tolerations` | Tolerations for NuoDB SM | `[]` |
| `sm.readinessTimeoutSeconds` | SM readiness probe timeout, sometimes needs adjusting depending on environment and pod resources | `5` |
| `te.externalAccess.enabled` | Whether to deploy a Layer 4 service for the database | `false` |
| `te.externalAccess.internalIP` | Whether to use an internal (to the cloud) or external (public) IP address for the load balancer. Only applies to external access of type `LoadBalancer` | `nil` |
| `te.externalAccess.type` | The service type used to enable external database access. The supported types are `NodePort` and `LoadBalancer` (defaults to `LoadBalancer`) | `nil` |
| `te.dbServices.enabled` | Whether to deploy clusterip and headless services for direct TE connections (defaults true) | `nil` |
| `te.logPersistence.enabled` | Whether to enable persistent storage for logs | `false` |
| `te.logPersistence.overwriteBackoff.copies` | How many copies of the crash directory to keep within windowMinutes | `3` |
| `te.logPersistence.overwriteBackoff.windowMinutes` | The window within which to keep the number of crash copies | `120` |
| `te.logPersistence.accessModes` | Volume access modes enabled (must match capabilities of the storage class).  This is expected to be ReadWriteMany.  Not all storage providers support this mode. | `ReadWriteMany` |
| `te.logPersistence.size` | Amount of disk space allocated for log storage | `60Gi` |
| `te.logPersistence.storageClass` | Storage class for volume backing log storage.  This storage class must be pre-configured in the cluster | `-` |
| `te.replicas` | TE replicas | `1` |
| `te.labels` | Labels given to the TEs started | `""` |
| `te.engineOptions` | Additional NuoDB engine options | `""` |
| `te.resources` | Labels to apply to all resources | `{}` |
| `te.affinity` | Affinity rules for NuoDB TE | `{}` |
| `te.nodeSelector` | Node selector rules for NuoDB TE | `{}` |
| `te.tolerations` | Tolerations for NuoDB TE | `[]` |
| `te.otherOptions` | Additional key/value Docker options | `{}` |
| `sm.otherOptions` | Additional key/value Docker options | `{}` |
| `te.readinessTimeoutSeconds` | TE readiness probe timeout, sometimes needs adjusting depending on environment and pod resources | `5` |
| `automaticProtocolUpgrade.enabled` | Enable automatic database protocol upgrade and a Transaction Engine (TE) restart as an upgrade finalization step done by Kubernetes Aware Admin (KAA). Applicable for NuoDB major versions upgrade only. Requires NuoDB 4.2.3+ | `false` |
| `automaticProtocolUpgrade.tePreferenceQuery` | LBQuery expression to select the TE that will be restarted after a successful database protocol upgrade. Defaults to random Transaction Engine (TE) in MONITORED state | `""` |

#### database.configFiles.*

The purpose of this section is to detail how to provide alternate configuration files for NuoDB. NuoDB has several configuration files that may be modified to suit.

There are two sets of configuration files documented:

- [Admin Configuration for a Particular Host][1]
- [Database Configuration for a Particular Host][2]

Any file located in `database.configFilesPath` can be replaced; the YAML key corresponds to the file name being created or replaced.

| Key | Description | Default |
| ----- | ----------- | ------ |
| `nuodb.config` | [NuoDB database options][6] | `nil` |


### Running

The purpose of this section is to allow customisation of the names of the clusterIP and balancer database services (load-balancers).

| Key | Description | Default |
| ----- | ----------- | ------ |
| `clusterip` | suffix for the clusterIP load-balancer | .Values.admin.serviceSuffix.clusterip |
| `balancer` | suffix for the balancer service | .Values.admin.serviceSuffix.balancer |

#### nuocollector.*

The purpose of this section is to specify the NuoDB monitoring parameters.

The following tables list the configurable parameters for the `nuocollector` option of the admin chart and their default values.

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `enabled` | Whether to enable NuoDB monitoring using sidecar containers |`false`|
| `image.registry` | NuoDB Collector container registry | `docker.io` |
| `image.repository` | NuoDB Collector container image name |`nuodb/nuodb-collector`|
| `image.tag` | NuoDB Collector container image tag | `1.1.0` |
| `image.pullPolicy` | NuoDB Collector container pull policy |`IfNotPresent`|
| `watcher.registry` | ConfigMap watcher container registry | `docker.io` |
| `watcher.repository` | ConfigMap watcher container image name |`kiwigrid/k8s-sidecar`|
| `watcher.tag` | ConfigMap watcher container image tag | `0.1.259` |
| `watcher.pullPolicy` | ConfigMap watcher container pull policy |`IfNotPresent`|
| `plugins.database` | NuoDB Collector additional plugins for database services |`{}`|

### Running

Verify the Helm chart:

```bash
helm install database nuodb/database --debug --dry-run
```

Deploy a database without backups:

**Tip**: If you plan to deploy NuoDB Insights visual monitoring, add the `--set nuocollector.enabled=true` switch as below.

```bash
helm install database nuodb/database \
    --set database.sm.hotCopy.replicas=0 --set database.sm.noHotCopy.replicas=1 --set nuocollector.enabled=true
```

Wait until the deployment completes:

```bash
helm status database
```

Verify the pods are running:

```bash
$ kubectl get pods
NAME                                              READY   STATUS    RESTARTS   AGE
admin-nuodb-cluster0-0                            1/1     Running   0          62s
disable-thp-transparent-hugepage-qc7mr            1/1     Running   0          71s
sm-database-nuodb-cluster0-demo-0                 3/3     Running   0          54s
te-database-nuodb-cluster0-demo-69d46b46c-2mznz   3/3     Running   0          54s
```


The command displays the NuoDB Pods running on the Kubernetes cluster. When completed, both the TE and the storage containers should show a **STATUS** of `Running`, and with 0 **RESTARTS**.

Verify the connected states of the database domain:

```bash
$ kubectl exec -it admin-nuodb-cluster0-0 -- nuocmd show domain

server version: 4.0.8-2-881d0e5d44, server license: Community
server time: 2021-02-04T19:42:45.729, client token: d7029e4f34b18b0e3e30444267d4f3ae89b60a3e
Servers:
  [admin-nuodb-cluster0-0] admin-nuodb-cluster0-0.nuodb.nuodb.svc.cluster.local:48005 [last_ack = 0.58] [member = ADDED] [raft_state = ACTIVE] (LEADER, Leader=admin-nuodb-cluster0-0, log=0/20/20) Connected *
Databases:
  demo [state = RUNNING]
    [SM] sm-database-nuodb-cluster0-demo-0/10.1.0.130:48006 [start_id = 0] [server_id = admin-nuodb-cluster0-0] [pid = 168] [node_id = 1] [last_ack =  0.40] MONITORED:RUNNING
    [TE] te-database-nuodb-cluster0-demo-69d46b46c-2mznz/10.1.0.129:48006 [start_id = 1] [server_id = admin-nuodb-cluster0-0] [pid = 98] [node_id = 2] [last_ack =  4.38] MONITORED:RUNNING
```

The command displays the status of NuoDB processes. The Servers section lists admin processes; they should all be **Connected**, one will be the **LEADER** and other designated as a **FOLLOWER**.

**Tip**: Wait until all processes are be in a **RUNNING** state.

to scale-out the TEs, run:

```bash
$ kubectl scale deployment te-database-nuodb-cluster0-demo --replicas=2
deployment.extensions/te-database-nuodb-cluster0-demo scaled
```

## Cleaning Up Archive References

This will clear the archive references and metadata from the admin layer if the default demo database was recreated

```bash
kubectl exec -it admin-nuodb-cluster0-0  -- /bin/bash

$ nuocmd get archives --db-name demo
$ nuocmd delete database --db-name demo
$ nuocmd delete archive --archive-id 0 --purge
$ nuocmd show domain
```

Then you must also clear the PVCs:

```bash
kubectl delete pvc archive-volume-sm-database-nuodb-cluster0-demo-0
kubectl delete pvc archive-volume-sm-database-nuodb-cluster0-demo-hotcopy-0
kubectl delete pvc backup-volume-sm-database-nuodb-cluster0-demo-hotcopy-0
```

## Uninstalling the Chart

To uninstall/delete the deployment:

```bash
helm delete database
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

[1]: #adminconfigfiles
[2]: #databaseconfigfiles
[3]: http://doc.nuodb.com/Latest/Content/Nuoadmin-Obtaining-and-Installing-NuoDB-Licenses.htm
[4]: http://doc.nuodb.com/Latest/Content/Nuoadmin-Host-Properties.htm
[5]: http://doc.nuodb.com/Latest/Content/Data-Type-Mappings.htm
[6]: http://doc.nuodb.com/Latest/Default.htm#Database-Options.htm
