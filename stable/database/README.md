# NuoDB Database Helm Chart

This chart starts a NuoDB database deployment on a Kubernetes cluster using the Helm package manager. To start a second database under the same NuoDB Admin deployment, start a second database using the same instructions with a new database name. 

## Command

```bash
helm install nuodb/database [--name releaseName] [--set parameter] [--values myvalues.yaml]
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
| `image.pullPolicy` | busybox container pull policy |`Always`|
| `image.pullSecrets` | Specify docker-registry secret names as an array | [] (does not add image pull secrets to deployed pods) |

The `registry` option can be used to connect to private image repositories, such as Artifactory.

The `registry`, `repository`, and `tag` values are combined to form the `image` declaration in the Helm charts.

For example, when using GlusterFS storage class, you would supply the following parameter:

```bash
  ...
  --set buzybox.image.registry=acme-dockerv2-virtual.jfrog.io
  ...
```

For example:

```yaml
busybox:
  image:
    registry: docker.io
    repository: busybox
    tag: latest
    pullPolicy: Always
```

#### nuodb.*

The purpose of this section is to specify the NuoDB image parameters.

The following tables list the configurable parameters for the `nuodb` option:

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `image.registry` | NuoDB container registry | `docker.io` |
| `image.repository` | NuoDB container image name |`nuodb/nuodb-ce`|
| `image.tag` | NuoDB container image tag | `latest` |
| `image.pullPolicy` | NuoDB container pull policy |`Always`|
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
    pullPolicy: Always
```

#### openshift.*

The purpose of this section is to specify parameters specific to OpenShift, e.g. enable features only present in OpenShift.

The following tables list the configurable parameters for the `openshift` option:

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `enabled` | Enable OpenShift features | `false` |

For example, to enable an OpenShift integration, and enable DeploymentConfigs:

```yaml
openshift:
  enabled: true
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

#### backup.*

The following tables list the configurable parameters of the `backup` portion of the `database` chart and their default values.

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `persistence.accessModes` | Volume access modes enabled (must match capabilities of the storage class) | `ReadWriteOnce` |
| `persistence.size` | Amount of disk space allocated for database backup storage | `20Gi` |
| `persistence.storageClass` | Storage class for volume backing database backup storage | `-` |

#### database.*

The following tables list the configurable parameters of the `database` chart and their default values.

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `nameOverride` | Provide a name in place of `database-daemonset` |`""`|
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
| `persistence.accessModes` | Volume access modes enabled (must match capabilities of the storage class) | `ReadWriteOnce` |
| `persistence.size` | Amount of disk space allocated for database archive storage | `20Gi` |
| `persistence.storageClass` | Storage class for volume backing database archive storage | `-` |
| `configFilesPath` | Directory path where `configFiles.*` are found | `/etc/nuodb/` |
| `configFiles.*` | See below. | `{}` |
| `sm.logPersistence.enabled` | Whether to enable persistent storage for logs | `false` |
| `sm.logPersistence.overwriteBackoff.copies` | How many copies of the crash directory to keep within windowMinutes | `3` |
| `sm.logPersistence.overwriteBackoff.windowMinutes` | The window within which to keep the number of crash copies | `120` |
| `sm.logPersistence.accessModes` | Volume access modes enabled (must match capabilities of the storage class) | `ReadWriteOnce` |
| `sm.logPersistence.size` | Amount of disk space allocated for log storage | `60Gi` |
| `sm.logPersistence.storageClass` | Storage class for volume backing log storage | `-` |
| `sm.hotCopy.replicas` | SM replicas with hot-copy enabled | `1` |
| `sm.hotCopy.enablePod` | Create DS/SS for hot-copy enabled SMs | `true` |
| `sm.hotcopy.deadline` | Deadline for a hotcopy job to start (seconds) | `1800` |
| `sm.hotcopy.timeout` | Timeout for a started hotcopy job to complete (seconds) | `1800` |
| `sm.hotcopy.successHistory` | Number of successful Jobs to keep | `5` |
| `sm.hotcopy.failureHostory` | Number of failed jobs to keep | `5` |
| `sm.hotcopy.backupDir` | Directory path where backiupsets will be stored | `/var/opt/nuodb/backup` |
| `sm.hotcopy.backupGroup` | Name of the backup group | `{{ .Values.cloud.cluster.name }}` |
| `sm.hotcopy.fullSchedule` | cron schedule for FULL hotcopy jobs | `35 22 * * 6` |
| `sm.hotcopy.incrementalSchedule` | cron schedule for INCREMENTAL hotcopy jobs | `35 22 * * 0-5` |
| `sm.hotcopy.persistence.size` | size of the hotcopy storage PV | `20Gi` |
| `sm.hotcopy.persistence.accessModes` | access modes for the hotcopy storage PV | `[ ReadWriteOnce ]` |
| `sm.hotcopy.persistence.size` | size of the hotcopy storage PV | `20Gi` |
| `sm.hotcopy.journalBackup.enabled` | Is `journal hotcopy` enabled - true/false | `false` |
| `sm.hotcopy.journalBackup.intervalMinutes` | Frequency of running `journal hotcopy` (minutes) | `15` |
| `sm.hotcopy.journalBackup.deadline` | Deadline for a `journal hotcopy` job to start (seconds) | `90` |
| `sm.hotcopy.journalBackup.timeout` | Timeout for a started `journal hotcopy` to complete (seconds) | `950` |
| `sm.hotcopy.coldStorage.credentials` | Credentials for accessing backup cold storage (user:password) | `""` |
| `sm.noHotCopy.replicas` | SM replicas with hot-copy disabled | `0` |
| `sm.noHotCopy.enablePod` | Create DS/SS for non-hot-copy SMs | `true` |
| `sm.memoryOption` | SM engine memory (*future deprecation*) | `"8g"` |
| `sm.labels` | Labels given to the SMs started | `{}` |
| `sm.engineOptions` | Additional NuoDB engine options | `{}` |
| `sm.resources` | Labels to apply to all resources | `{}` |
| `sm.affinity` | Affinity rules for NuoDB SM | `{}` |
| `sm.nodeSelector` | Node selector rules for NuoDB SM | `{}` |
| `sm.tolerations` | Tolerations for NuoDB SM | `[]` |
| `te.externalAccess.enabled` | Whether to deploy a Layer 4 cloud load balancer service for the admin layer | `false` |
| `te.externalAccess.internalIP` | Whether to use an internal (to the cloud) or external (public) IP address for the load balancer | `nil` |
| `te.dbServices.enabled` | Whether to deploy clusterip and headless services for direct TE connections (defaults true) | `nil` |
| `te.replicas` | TE replicas | `1` |
| `te.memoryOption` | TE engine memory (*future deprecation*) | `"8g"` |
| `te.labels` | Labels given to the TEs started | `""` |
| `te.engineOptions` | Additional NuoDB engine options | `""` |
| `te.resources` | Labels to apply to all resources | `{}` |
| `te.affinity` | Affinity rules for NuoDB TE | `{}` |
| `te.nodeSelector` | Node selector rules for NuoDB TE | `{}` |
| `te.tolerations` | Tolerations for NuoDB TE | `[]` |
| `te.otherOptions` | Additional key/value Docker options | `{}` |
| `sm.affinityNoHotCopyDS` | Affinity rules for non-hot-copy SMs (DaemonSet) | `{}` |
| `sm.affinityHotCopyDS` | Affinity rules for hot-copy enabled SMs (DaemonSet) | `{}` |
| `sm.nodeSelectorHotCopyDS` | Node selector rules for hot-copy enabled SMs (DaemonSet) | `{}` |
| `sm.nodeSelectorNoHotCopyDS` | Node selector rules for non-hot-copy SMs (DaemonSet) | `{}` |
| `sm.tolerationsDS` | Tolerations for SMs (DaemonSet) | `[]` |
| `sm.otherOptions` | Additional key/value Docker options | `{}` |

#### database.configFiles.*

The purpose of this section is to detail how to provide alternate configuration files for NuoDB. NuoDB has several configuration files that may be modified to suit.

There are two sets of configuration files documented:

- [Admin Configuration for a Particular Host][1]
- [Database Configuration for a Particular Host][2]

Any file located in `database.configFilesPath` can be replaced; the YAML key corresponds to the file name being created or replaced.

| Key | Description | Default |
| ----- | ----------- | ------ |
| `nuodb.config` | [NuoDB database options][6] | `nil` |

#### database.serviceSuffix

The purpose of this section is to allow customisation of the names of the clusterIP and balancer database services (load-balancers).

| Key | Description | Default |
| ----- | ----------- | ------ |
| `clusterip` | suffix for the clusterIP load-balancer | .Values.admin.serviceSuffix.clusterip |
| `balancer` | suffix for the balancer service | .Values.admin.serviceSuffix.balancer |

### Running

Verify the Helm chart:

```bash
helm install nuodb/database --name database \
    --debug --dry-run
```

Deploy a database without backups:

```bash
helm install nuodb/database --name database \
    --set database.sm.hotcopy.replicas=0 --set database.sm.nohotcopy.replicas=1
```

Wait until the deployment completes:

```bash
helm status database
```

Verify the pods are running:

```bash
$ kubectl get pods
NAME                                               READY   STATUS      RESTARTS   AGE
admin-nuodb-0                                      1/1     Running     0          18h
disable-thp-transparent-hugepage-59f7q             1/1     Running     0          18h
sm-database-cashews-demo-0                         1/1     Running     0          18h
sm-database-cashews-demo-hotcopy-0                 1/1     Running     0          18h
te-database-cashews-demo-599ff97797-dtqkk          1/1     Running     0          18h
tiller-deploy-88ff958dd-pgsjn                      1/1     Running     0          23h
```

The command displays the NuoDB Pods running on the Kubernetes cluster. When completed, both the TE and the storage containers should show a **STATUS** of `Running`, and with 0 **RESTARTS**.

Verify the connected states of the database domain:

```bash
$ kubectl exec -it admin-nuodb-0 -- nuocmd show domain

server version: 4.0-2-ef765f7906, server license: Enterprise
server time: 2019-08-29T13:31:10.325, client token: b2c99602e831c0ad61e9becd518e4d5b323d6b3f
Servers:
  [admin-cashews-0] admin-cashews-0.cashews.nuodb.svc:48005 [last_ack = 1.81] [member = ADDED] [raft_state = ACTIVE] (LEADER, Leader=admin-cashews-0, log=0/6535/6535) Connected *
Databases:
  demo [state = RUNNING]
    [SM] sm-database-cashews-demo-0/10.28.7.84:48006 [start_id = 0] [server_id = admin-cashews-2] [pid = 87] [node_id = 2] [last_ack = 10.17] MONITORED:RUNNING
    [SM] sm-database-cashews-demo-hotcopy-0/10.28.2.172:48006 [start_id = 1] [server_id = admin-cashews-0] [pid = 87] [node_id = 1] [last_ack =  4.49] MONITORED:RUNNING
    [TE] te-database-cashews-demo-599ff97797-dtqkk/10.28.3.166:48006 [start_id = 2] [server_id = admin-cashews-0] [pid = 86] [node_id = 3] [last_ack =  3.68] MONITORED:RUNNING
```

The command displays the status of NuoDB processes. The Servers section lists admin processes; they should all be **Connected**, one will be the **LEADER** and other designated as a **FOLLOWER**.

  **Tip**: Wait until all processes are be in a **RUNNING** state.

Now to scale the TEs is simple enough:

```bash
$ kubectl scale deployment te-database-nuodb-demo --replicas=2
deployment.extensions/te-database-nuodb-demo scaled
```

## Cleaning Up Archive References

This will clear the archive references and metadata from the admin layer if the default demo database was recreated

```bash
kubectl exec -it admin-nuodb-0  -- /bin/bash

$ nuocmd get archives --db-name demo
$ nuocmd delete archive --archive-id 0 --purge
$ nuocmd delete database --db-name demo
$ nuocmd show domain
```

Then you must also clear the PVCs:

```bash
kubectl delete pvc archive-volume-sm-database-nuodb-demo-0
kubectl delete pvc archive-volume-sm-database-nuodb-demo-hotcopy-0
kubectl delete pvc backup-volume-sm-database-nuodb-demo-hotcopy-0
```

## Uninstalling the Chart

To uninstall/delete the deployment:

```bash
helm del --purge database
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

[1]: #adminconfigfiles
[2]: #databaseconfigfiles
[3]: http://doc.nuodb.com/Latest/Content/Nuoadmin-Obtaining-and-Installing-NuoDB-Licenses.htm
[4]: http://doc.nuodb.com/Latest/Content/Nuoadmin-Host-Properties.htm
[5]: http://doc.nuodb.com/Latest/Content/Data-Type-Mappings.htm
[6]: http://doc.nuodb.com/Latest/Default.htm#Database-Options.htm
