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
| `image.pullSecrets` | Specify docker-registry secret names as an array | `[]` (does not add image pull secrets to deployed pods) |

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
| `image.repository` | NuoDB container image name |`nuodb/nuodb`|
| `image.tag` | NuoDB container image tag | `"6.0"` |
| `image.pullPolicy` | NuoDB container pull policy |`IfNotPresent`|
| `image.pullSecrets` | Specify docker-registry secret names as an array | `[]` (does not add image pull secrets to deployed pods) |
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
    repository: nuodb/nuodb
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
| `tlsKeyStore.passwordKey` | TLS keystore secret password key. One of `tlsKeyStore.password` or `tlsKeyStore.passwordKey` must be used | `password` |
| `tlsClientPEM.secret` | TLS client PEM secret name | `nil` |
| `tlsClientPEM.key` | TLS client PEM secret key | `nil` |
| `tde.secrets` | Transparent Data Encryption secret names used for different databases | `{}` |
| `tde.storagePasswordsDir` | Transparent Data Encryption storage passwords mount path | `/etc/nuodb/tde` |
| `affinityLabels` | List of AP label keys to check for affinity with the engine process, ordered by precedence; the AP with the earliest label key, whose value matches the corresponding engine label, will be chosen to manage the engine process. For example, specifying "zone region" means "find an AP in the same 'zone' and, if there are none, select one in the same 'region'". Supported starting from NuoDB image version 6.0.3 | `"node zone region"` |

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
| `clusterip` | suffix for the clusterIP load-balancer | `clusterip` |
| `balancer` | suffix for the balancer service | `balancer` |
| `nodeport` | suffix for the NodePort service | `nodeport` |

#### database.*

The following tables list the configurable parameters of the `database` chart and their default values.

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `nameOverride` | Provide a name in place of `database-dbname` |`""`|
| `fullNameOverride` | Provide a name to substitute for the full names of resources |`""`|
| `name` | Database name | `demo` |
| `rootUser` | Database root user | `dba` |
| `rootPassword` | Database root password | `secret` |
| `archiveType` | Archive Type (either `""` or `"lsa"`) | `""` |
| `priorityClasses.sm` | The `priorityClassName` to specify SM pods | `""` |
| `priorityClasses.te` | The `priorityClassName` to specify TE pods | `""` |
| `initContainers.runInitDisk` | Whether to run the `init-disk` init container to set volume permissions | `true` |
| `initContainers.runInitDiskAsRoot` | Whether to run the `init-disk` init container as root | `true` |
| `initContainers.resources` | Kubernetes resource requests and limits set on the database init containers | `{}` |
| `securityContext.fsGroupOnly` | Creates a security context for Pods containing only the `securityContext.fsGroup` value | `false` |
| `securityContext.runAsNonRootGroup` | Creates a security context for Pods containing a non-root user and group (1000:1000) along with the `securityContext.fsGroup` value | `false` |
| `securityContext.enabled` | Creates a security context for Pods containing the `securityContext.runAsUser` and `securityContext.fsGroup` values | `false` |
| `securityContext.runAsUser` | The user ID for the Pod security context created if `securityContext.enabled` is `true`. | `1000` |
| `securityContext.fsGroup` | The `fsGroup` for the Pod security context created if any of `securityContext.fsGroupOnly`, `securityContext.runAsNonRootGroup`, or `securityContext.enabled` are `true`. | `1000` |
| `securityContext.enabledOnContainer` | Whether to create SecurityContext for containers | `false` |
| `securityContext.capabilities` | Capabilities for to engine container security context | `{ add: [], drop: [] }` |
| `securityContext.privileged` | Run the NuoDB database containers in privileged mode. Processes in privileged containers are essentially equivalent to root on the host | `false` |
| `securityContext.allowPrivilegeEscalation` | Whether a process can gain more privileges than its parent process. This boolean directly controls if the `no_new_privs` flag will be set on the container process | `false` |
| `securityContext.readOnlyRootFilesystem` | Whether to mount the root filesystem as read-only | `false` |
| `env` | Import ENV vars inside containers | `[]` |
| `envFrom` | Import ENV vars from configMap(s) | `[]` |
| `lbConfig.prefilter` | Database load balancer prefilter expression | `nil` |
| `lbConfig.default` | Database load balancer default query | `nil` |
| `persistence.accessModes` | Volume access modes enabled (must match capabilities of the storage class) | `ReadWriteOnce` |
| `persistence.size` | Amount of disk space allocated for database archive storage | `20Gi` |
| `persistence.storageClass` | Storage class for volume backing database archive storage | `-` |
| `persistence.archiveDataSource.*` | The data source to use to initialize the archive volume. This takes precedence over the archive snapshot resolved from `database.snapshotRestore`. | `disabled` |
| `persistence.journalDataSource.*` | The data source to use to initialize the journal volume. This takes precedence over the journal snapshot resolved from `database.snapshotRestore`. | `disabled` |
| `persistence.validateDataSources` | Whether to validate that data sources resolved from `database.snapshotRestore` or specified explicitly using `persistence.archiveDataSource` and `persistence.journalDataSource` actually exist. This is useful because data source references that do not exist are silently ignored by Kubernetes. | `true` |
| `persistence.preprovisionVolumes` | Whether to explicitly provision PVCs for all SMs specified by `database.sm.noHotCopy.replicas`. If data sources are configured, then they will appear in the preprovisioned PVCs of the non-hotcopy SMs but not in the `volumeClaimTemplates` section, so that the existence of the volume snapshots is not required in order to scale up the non-hotcopy SM StatefulSet. | `false` |
| `configFilesPath` | Directory path where `configFiles.*` are found | `/etc/nuodb/` |
| `configFiles.*` | See below. | `{}` |
| `podAnnotations` | Annotations to pass through to the SM an TE pods | `nil` |
| `primaryRelease` | One primary and multiple secondary database Helm releases for the same database name may be deployed into the same Kubernetes namespace. Set to `false` when deploying secondary database Helm releases. | `true` |
| `autoImport.*` | Enable and configure the automatic initializing of the initial database state | `disabled` |
| `autoImport.source` | The source - typically a URL - of the database copy to import | `""` |
| `autoImport.credentials` | Authentication for the download of `source` in the form `user`:`password` | `""` |
| `autoImport.stripLevels` | The number of levels to strip off pathnames when unpacking a TAR file of an archive | `1` |
| `autoImport.type` | Type of content in `source`. One of `stream` -> exact copy of an archive; or `backupset` -> a NuoDB hotcopy backupset | `backupset` |
| `autoRestore.*` | Enable and configure the automatic re-initialization of a single archive in a running database - see the options in `autoImport` | `disabled` |
| `backupHooks.enabled` | Whether to enable the backup hooks sidecar for non-hotcopy SMs | `false` |
| `backupHooks.resources` | Kubernetes resource requests and limits set on the backup hook sidecar container | `{}` |
| `backupHooks.freezeMode` | The freeze mode to be used when executing backup hooks. Supported modes are `hotsnap`, `fsfreeze` and `suspend`. Defaults to `hostsnap` if empty | `""` |
| `backupHooks.timeout` | Timeout in seconds after which the archive will be automatically unfrozen | `30` |
| `backupHooks.customHandlers` | Custom handlers to register on HTTP server in sidecar | `[]` |
| `backupHooks.customHandlers[*].method` | The HTTP request method to match on | |
| `backupHooks.customHandlers[*].path` | The HTTP request path to match on, which may contain path parameters in the form `{param_name}` | |
| `backupHooks.customHandlers[*].script` | The script to invoke when handling the matched request, which may reference path parameters, query parameters, or the request payload (as `$payload`). If the same variable name appears as a query and path parameter, or a path parameter appears named `$payload`, the path parameter takes precedence. | |
| `backupHooks.customHandlers[*].statusMappings` | Mapping of script exit codes to HTTP status codes | |
| `snapshotRestore.backupId` | The backup ID being restored, which is set to enable restore from data sources | `""` |
| `snapshotRestore.snapshotNameTemplate` | The template used to resolve the names of snapshots to use as data sources for the archive and journal PVCs. The template can reference `backupId` and `volumeType`, which is one of `archive`, `journal`. | `{{.backupId}}-{{.volumeType}}` |
| `ephemeralVolume.enabled` | Whether to create a generic ephemeral volume rather than emptyDir for any storage that does not outlive the pod | `false` |
| `ephemeralVolume.size` | The size of the generic ephemeral volume to create | `1Gi` |
| `ephemeralVolume.sizeToMemory` | Whether to size the generic ephemeral volume based on the `resources.limits.memory` setting of the process so that at least one core file is retained for the lifetime of the pod | `false` |
| `ephemeralVolume.storageClass` | The storage class to use for the generic ephemeral volume | `nil` |
| `sm.logPersistence.enabled` | Whether to enable persistent storage for logs | `false` |
| `sm.logPersistence.overwriteBackoff.copies` | How many copies of the crash directory to keep within windowMinutes | `3` |
| `sm.logPersistence.overwriteBackoff.windowMinutes` | The window within which to keep the number of crash copies | `120` |
| `sm.logPersistence.accessModes` | Volume access modes enabled (must match capabilities of the storage class) | `ReadWriteOnce` |
| `sm.logPersistence.size` | Amount of disk space allocated for log storage | `60Gi` |
| `sm.logPersistence.storageClass` | Storage class for volume backing log storage.  This storage class must be pre-configured in the cluster | `-` |
| `sm.hotCopy.replicas` | SM replicas with hot-copy enabled | `1` |
| `sm.hotCopy.enablePod` | Create StatefulSet for hot-copy enabled SMs | `true` |
| `sm.hotCopy.enableBackups` | Enable full/incremental/journal backups | `true` |
| `sm.hotCopy.deadline` | Deadline for a hotcopy job to start (seconds) | `1800` |
| `sm.hotCopy.timeout` | Timeout for a started `full` or `incremental` hotcopy job to complete (seconds). The default timeout of "0" will force the backup jobs to wait forever for the requested hotcopy operation to complete | `0` |
| `sm.hotCopy.successHistory` | Number of successful Jobs to keep | `5` |
| `sm.hotCopy.failureHistory` | Number of failed jobs to keep | `5` |
| `sm.hotCopy.backupDir` | Directory path where backupsets will be stored | `/var/opt/nuodb/backup` |
| `sm.hotCopy.jobResources` | Kubernetes resource requests and limits set on the backup jobs | `{}` |
| `sm.hotCopy.backupGroupPrefix` | Prefix for the automatically generated backup groups | `{{ .Values.cloud.cluster.name }}` |
| `sm.hotCopy.backupGroups` | Backup groups configuration. By default a backup group per HCSM is created automatically | `{}` |
| `sm.hotCopy.backupGroups.<name>.labels` | Space separated process labels used to select the Storage Managers which are part of this backup group. _Any_ label key and value should match for the SM to be selected | `nil` |
| `sm.hotCopy.backupGroups.<name>.processFilter` | [LBQuery expression](https://doc.nuodb.com/nuodb/latest/client-development/load-balancer-policies/#lbquery-expression-syntax) used to select the Storage Managers which are part of this backup group. | `nil` |
| `sm.hotCopy.backupGroups.<name>.fullSchedule` | cron schedule for _FULL_ hot copy performed by this backup group. If not defined the `sm.hotCopy.fullSchedule` setting will be used | `nil` |
| `sm.hotCopy.backupGroups.<name>.incrementalSchedule` | cron schedule for _INCREMENTAL_ hot copy performed by this backup group. If not defined the `sm.hotCopy.incrementalSchedule` setting will be used | `nil` |
| `sm.hotCopy.backupGroups.<name>.journalSchedule` | cron schedule for _JOURNAL_ hot copy performed by this backup group. If not defined the `sm.hotCopy.journalBackup.journalSchedule` setting will be used | `nil` |
| `sm.hotCopy.fullSchedule` | cron schedule for _FULL_ hotcopy jobs | `35 22 * * 6` |
| `sm.hotCopy.incrementalSchedule` | cron schedule for _INCREMENTAL_ hotcopy jobs | `35 22 * * 0-5` |
| `sm.hotCopy.restartPolicy` | Restart policy for backup related JOBs and CRON JOBs | `OnFailure` |
| `sm.hotCopy.persistence.size` | size of the hotcopy storage PV | `20Gi` |
| `sm.hotCopy.persistence.accessModes` | access modes for the hotcopy storage PV | `[ ReadWriteOnce ]` |
| `sm.hotCopy.persistence.size` | size of the hotcopy storage PV | `20Gi` |
| `sm.hotCopy.journalBackup.enabled` | Is `journal hotcopy` enabled - true/false | `false` |
| `sm.hotCopy.journalBackup.journalSchedule` | cron schedule for _JOURNAL_ hotcopy jobs. When journal backup is enabled, an SM will retain each journal file on disk until it is journal hot copied into a backup set. This means that journal hot copy must be executed periodically to prevent SMs from running out of disk space on the journal volume | `?/15 * * * *` |
| `sm.hotCopy.journalBackup.deadline` | Deadline for a `journal hotcopy` job to start (seconds) | `90` |
| `sm.hotCopy.journalBackup.timeout` | Timeout for a started `journal hotcopy` to complete (seconds). The default timeout of "0" will force the backup jobs to wait forever for the requested hotcopy operation to complete | `0` |
| `sm.hotCopy.coldStorage.credentials` | Credentials for accessing backup cold storage (user:password) | `""` |
| `sm.hotCopy.journalPath.enabled` | Whether to enable separate SM journal directory. For more info, read the [Journal HowTo](../../docs/HowToArchiveJournal.md) | `false` |
| `sm.hotCopy.journalPath.accessModes` | Volume access modes enabled (must match capabilities of the storage class) | `ReadWriteOnce` |
| `sm.hotCopy.journalPath.size` | Amount of disk space allocated for SM journal | `20Gi` |
| `sm.hotCopy.journalPath.storageClass` | Storage class for SM journal.  This storage class must be pre-configured in the cluster | `-` |
| `sm.noHotCopy.replicas` | SM replicas with hot-copy disabled | `0` |
| `sm.noHotCopy.enablePod` | Create StatefulSet for non-hot-copy SMs | `true` |
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
| `sm.topologySpreadConstraints` | Topology spread constraints for NuoDB SM | `[]` |
| `sm.otherOptions` | Additional key/value Docker options | `{}` |
| `sm.readinessProbe.initialDelaySeconds` | The initial delay in seconds for the readiness probe. | `5` |
| `sm.readinessProbe.periodSeconds` | The period in seconds for the readiness probe. | `5` |
| `sm.readinessProbe.failureThreshold` | The number of times that the readiness probe must fail before the container is marked as unready. | `3` |
| `sm.readinessProbe.successThreshold` | The number of times that the readiness probe must success before the container is marked as ready. | `1` |
| `sm.readinessProbe.timeoutSeconds` | The timeout in seconds for an invocation of the readiness probe. | `5` |
| `sm.storageGroup.enabled` | Enable Table Partitions and Storage Groups (TPSG) for all SMs in this database Helm release | `false` |
| `sm.storageGroup.name` | The name of the storage group. Only alphanumeric and underscore (`_`) characters are allowed. By default the Helm release name is used | `{{ .Release.Name }}` |
| `te.enablePod` | Create deployment for TEs. By default, the TE Deployment is disabled if TP/SG is enabled and this is a "secondary" release. | `nil` |
| `te.externalAccess.enabled` | Whether to deploy a Layer 4 service for the database | `false` |
| `te.externalAccess.internalIP` | Whether to use an internal (to the cloud) or external (public) IP address for the load balancer. Only applies to external access of type `LoadBalancer` | `nil` |
| `te.externalAccess.type` | The service type used to enable external database access. The supported types are `NodePort` and `LoadBalancer` (defaults to `LoadBalancer`) | `nil` |
| `te.externalAccess.annotations` | Annotations to pass through to the Service of type `LoadBalancer` | `{}` |
| `te.ingress.enabled` | Whether to deploy an Ingress resources for the NuoDB Database. Supported starting with Kubernetes v1.19.0 and NuoDB image version 4.2.3 | `false` |
| `te.ingress.hostname` | The fully qualified domain name (FQDN) of the network host used by SQL clients to reach this database | `""` |
| `te.ingress.className` | The associated IngressClass name defines which Ingress controller will implement the resource | `""` |
| `te.ingress.annotations` | Custom annotations that are set on the Ingress resource. SSL passthrough feature should be configured so that SQL clients TLS connections are send directly to the Transaction Engines | `{ ingress.kubernetes.io/ssl-passthrough: "true" }` |
| `te.dbServices.enabled` | Whether to deploy clusterip and headless services for direct TE connections (defaults true) | `nil` |
| `te.logPersistence.enabled` | Whether to enable persistent storage for logs | `false` |
| `te.logPersistence.overwriteBackoff.copies` | How many copies of the crash directory to keep within windowMinutes | `3` |
| `te.logPersistence.overwriteBackoff.windowMinutes` | The window within which to keep the number of crash copies | `120` |
| `te.logPersistence.accessModes` | Volume access modes enabled (must match capabilities of the storage class).  This is expected to be ReadWriteMany.  Not all storage providers support this mode. | `ReadWriteMany` |
| `te.logPersistence.size` | Amount of disk space allocated for log storage | `60Gi` |
| `te.logPersistence.storageClass` | Storage class for volume backing log storage.  This storage class must be pre-configured in the cluster | `-` |
| `te.replicas` | Number of Transaction Engine (TE) replicas. A non-zero value is discarded if TE autoscaling is enabled. | `1` |
| `te.labels` | Labels given to the TEs started | `""` |
| `te.engineOptions` | Additional NuoDB engine options | `""` |
| `te.resources` | Labels to apply to all resources | `{}` |
| `te.affinity` | Affinity rules for NuoDB TE | `{}` |
| `te.nodeSelector` | Node selector rules for NuoDB TE | `{}` |
| `te.tolerations` | Tolerations for NuoDB TE | `[]` |
| `te.topologySpreadConstraints` | Topology spread constraints for NuoDB TE | `[]` |
| `te.otherOptions` | Additional key/value Docker options | `{}` |
| `te.readinessProbe.initialDelaySeconds` | The initial delay in seconds for the readiness probe. | `5` |
| `te.readinessProbe.periodSeconds` | The period in seconds for the readiness probe. | `5` |
| `te.readinessProbe.failureThreshold` | The number of times that the readiness probe must fail before the container is marked as unready. | `3` |
| `te.readinessProbe.successThreshold` | The number of times that the readiness probe must success before the container is marked as ready. | `1` |
| `te.readinessProbe.timeoutSeconds` | The timeout in seconds for an invocation of the readiness probe. | `5` |
| `te.autoscaling.minReplicas` | The lower limit for the number of TE replicas to which the autoscaler can scale down. | `1` |
| `te.autoscaling.maxReplicas` | The upper limit for the number of TE replicas to which the autoscaler can scale up. It cannot be less than the minReplicas. | `3` |
| `te.autoscaling.hpa.enabled` | Whether to enable auto-scaling for TE deployment by using HPA resource. | `false` |
| `te.autoscaling.hpa.targetCpuUtilization` | The target average CPU utilization value across all TE pods, represented as a percentage. | `80` |
| `te.autoscaling.hpa.behavior` | Configures the scaling behavior of the target in both Up and Down directions. | `...` |
| `te.autoscaling.hpa.annotations` | Custom annotations set on the HPA resource | `{}` |
| `te.autoscaling.hpa.behavior.scaleUp.stabilizationWindowSeconds` | The number of seconds for which past recommendations should be considered while scaling up. | `300` |
| `te.autoscaling.keda.enabled` | Whether to enable auto-scaling for TE deployment by using KEDA ScaledObject resource. | `false` |
| `te.autoscaling.keda.pollingInterval` | The interval in seconds to check each trigger on. | `30` |
| `te.autoscaling.keda.cooldownPeriod` | The period in seconds to wait after the last trigger reported active before scaling the resource back to 0. | `300` |
| `te.autoscaling.keda.fallback` | The number of replicas to fall back to if a scaler is in an error state. See https://keda.sh/docs/latest/reference/scaledobject-spec/ | `{}` |
| `te.autoscaling.keda.triggers` | List of triggers to activate scaling of the target resource. See https://keda.sh/docs/latest/scalers/ | `[]` |
| `te.autoscaling.keda.annotations` | Custom annotations set on the ScaledObject resource | `{}` |
| `automaticProtocolUpgrade.enabled` | Enable automatic database protocol upgrade and a Transaction Engine (TE) restart as an upgrade finalization step done by Kubernetes Aware Admin (KAA). Applicable for NuoDB major versions upgrade only. Requires NuoDB 4.2.3+ | `false` |
| `automaticProtocolUpgrade.tePreferenceQuery` | LBQuery expression to select the TE that will be restarted after a successful database protocol upgrade. Defaults to random Transaction Engine (TE) in MONITORED state | `""` |
| `resourceLabels` | Custom labels attached to the Kubernetes resources installed by this Helm chart. The labels are immutable and can't be changed with Helm upgrade | `{}` |

#### database.configFiles.*

The purpose of this section is to detail how to provide alternate configuration files for NuoDB. NuoDB has several configuration files that may be modified to suit.

There are two sets of configuration files documented:

- [Admin Configuration for a Particular Host][1]
- [Database Configuration for a Particular Host][2]

Any file located in `database.configFilesPath` can be replaced; the YAML key corresponds to the file name being created or replaced.

| Key | Description | Default |
| ----- | ----------- | ------ |
| `nuodb.config` | [NuoDB database options][6] | `nil` |

#### database.serviceSuffix.*

The purpose of this section is to allow customisation of the names of the clusterIP and balancer database services (load-balancers).

| Key | Description | Default |
| ----- | ----------- | ------ |
| `clusterip` | suffix for the clusterIP load-balancer | `{{ .Values.admin.serviceSuffix.clusterip }}` |
| `balancer` | suffix for the balancer service | `{{ .Values.admin.serviceSuffix.balancer }}` |
| `nodeport` | suffix for the nodePort service | `{{ .Values.admin.serviceSuffix.nodeport }}` |

#### `database.snapshotRestore.*`, `database.persistence.archiveDataSource.*`, and `database.persistence.journalDataSource.*`

The `database.snapshotRestore` section can be specified to restore a database from volume snapshots. `database.snapshotRestore.snapshotNameTemplate` can be used to resolve the names of the archive and journal volume snapshots within the same namespace.

To restore a database from volume snapshots within a different namespace or using PVCs as data sources rather than volume snapshots, the `database.persistence.archiveDataSource` and `database.persistence.journalDataSource` sections can be used to explicitly specify the data sources. When specifying a PVC as a data source, the `apiGroup` should be set to empty.

| Key | Description | Default |
| ----- | ----------- | ------ |
| `name` | Backup or volume name. The entire dataSource will be omitted if this value is empty | `nil` |
| `namespace` | Namespace containing the source | `nil` |
| `kind` | Data source kind | `VolumeSnapshot` |
| `apiGroup` | APIGroup is the group for the resource. If APIGroup is not specified, the specified Kind must be in the core API group. | `snapshot.storage.k8s.io` |

#### database.podMonitor.*

The purpose of this section is to allow metrics from database pods to be scraped by Prometheus server.

| Key | Description | Default |
| ----- | ----------- | ------ |
| `enabled` | Whether to enable PodMonitor resource for database pods. | `false` |
| `labels` | Labels to assign to the PodMonitor resource. Prometheus must be configured with label selector that matches the defined labels on the resource. | `{}` |
| `jobLabel` | The label to use to retrieve the job name from. | `app` |
| `podTargetLabels` | The labels which are transferred from the associated Kubernetes Pod object onto the ingested metrics. | `[]` |
| `portName` | The Pod port name which exposes the endpoint. | `http-metrics` |
| `interval` | Interval at which Prometheus scrapes the metrics from the database pods. | `10s` |
| `path` | HTTP path from which to scrape for metrics. | `/metrics` |
| `interval` | Interval at which Prometheus scrapes the metrics from the database pods. | `10s` |
| `scrapeTimeout` | Timeout after which Prometheus considers the scrape to be failed. If empty, Prometheus uses the global scrape timeout unless it is less than the target’s scrape interval value in which the latter is used. | `""` |
| `scheme` | HTTP scheme to use for scraping. | `http` |
| `tlsConfig` | TLS configuration to use when scraping the target. | `{}` |
| `relabelings` | The relabeling rules to apply to the samples before ingestion. | `[]` |
| `metricRelabelings` | The relabeling rules to apply to the samples before ingestion. | `[]` |
| `basicAuth` | Configures the Basic Authentication credentials to use when scraping. | `{}` |

#### database.legacy

Features in this section have been deprecated but not yet removed.

| Key | Description | Default |
| ----- | ----------- | ------ |
| `headlessService.enabled` | Create a headless service for this database. Use the TE group `ClusterIP` service instead. | `false` |
| `directService.enabled` | Create a service for direct connections to the database. Use the TE group `ClusterIP` service instead. | `false` |

#### nuocollector.*

The purpose of this section is to specify the NuoDB monitoring parameters.

The following tables list the configurable parameters for the `nuocollector` option of the admin chart and their default values.

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `enabled` | Whether to enable NuoDB monitoring using sidecar containers |`false`|
| `image.registry` | NuoDB Collector container registry | `docker.io` |
| `image.repository` | NuoDB Collector container image name |`nuodb/nuodb-collector`|
| `image.tag` | NuoDB Collector container image tag | `2.0.0` |
| `image.pullPolicy` | NuoDB Collector container pull policy |`IfNotPresent`|
| `watcher.registry` | ConfigMap watcher container registry | `docker.io` |
| `watcher.repository` | ConfigMap watcher container image name |`kiwigrid/k8s-sidecar`|
| `watcher.tag` | ConfigMap watcher container image tag | `0.1.259` |
| `watcher.pullPolicy` | ConfigMap watcher container pull policy |`IfNotPresent`|
| `plugins.database` | NuoDB Collector additional plugins for database services |`{}`|
| `resources` | Kubernetes resource requests and limits used for the nuocollector sidecar |`{}`|
| `ports` | Ports to expose on the nuocollector container |`[]`|
| `env` | Import environment variables inside nuocollector container |`[]`|

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

```console
$ kubectl get pods
NAME                                              READY   STATUS    RESTARTS   AGE
admin-nuodb-cluster0-0                            1/1     Running   0          62s
disable-thp-transparent-hugepage-qc7mr            1/1     Running   0          71s
sm-database-nuodb-cluster0-demo-0                 3/3     Running   0          54s
te-database-nuodb-cluster0-demo-69d46b46c-2mznz   3/3     Running   0          54s
```


The command displays the NuoDB Pods running on the Kubernetes cluster. When completed, both the TE and the storage containers should show a **STATUS** of `Running`, and with 0 **RESTARTS**.

Verify the connected states of the database domain:

```console
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

```console
$ kubectl scale deployment te-database-nuodb-cluster0-demo --replicas=2
deployment.extensions/te-database-nuodb-cluster0-demo scaled
```

## Cleaning Up Archive References

This will clear the archive references and metadata from the admin layer if the default demo database was recreated

```
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
[3]: https://doc.nuodb.com/nuodb/latest/deployment-models/physical-or-vmware-environments-with-nuodb-admin/installing-nuodb/obtaining-and-installing-an-enterprise-edition-license
[4]: https://doc.nuodb.com/nuodb/latest/reference-information/configuration-files/host-properties-nuoadmin.conf
[5]: https://doc.nuodb.com/nuodb/latest/reference-information/configuration-files/data-type-mappings-nuodb-types.config
[6]: https://doc.nuodb.com/nuodb/latest/reference-information/database-options
