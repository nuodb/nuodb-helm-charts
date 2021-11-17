# NuoDB Restore Helm Chart

This chart starts a Job to restore a NuoDB database from an existing backup, in a Kubernetes cluster using the Helm package manager.

## Command

```bash
helm install nuodb/restore [--generate-name | --name releaseName] [--set parameter] [--values myvalues.yaml]
```

## Software Version Prerequisites

Please visit the **[NuoDB Helm Chart main page](https://github.com/nuodb/nuodb-helm-charts/#software-release-requirements)** for software version prerequisites.

### Installing the Chart

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

For example, for the Google Cloud:

```yaml
cloud:
  provider: google
  zones:
    - us-central1-a
    - us-central1-b
    - us-central1-c
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

#### admin.*

The following tables list the configurable parameters of the backup chart and their default values.

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `admin.domain` | NuoDB admin cluster name | `domain` |
| `admin.namespace` | Namespace where admin is deployed; when peering to an existing admin cluster provide its project name | `nuodb` |
| `tlsCACert.secret` | TLS CA certificate secret name | `nil` |
| `tlsCACert.key` | TLS CA certificate secret key | `nil` |
| `tlsClientPEM.secret` | TLS client PEM secret name | `nil` |
| `tlsClientPEM.key` | TLS client PEM secret key | `nil` |

#### database.*

The following tables list the configurable parameters of the database chart and their default values.

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `name` | Database name | `demo` |

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

#### restore.*

The following tables list the configurable parameters of the restore chart and their default values.

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `restore.type` | What type of restore to perform: [ "database" \| "archive" ]. A "database" restore restarts the entire database at a previous state. An "archive" restore restores/repairs a SINGLE archive in a RUNNING database. | `"database"` |
| `restore.target` | The database name to request restore operations for | `"demo"` |
| `restore.source` | The source which will be restored from. Supported values are [ _backupset_ \| _URL_ \| `:latest` \| `<backup group>:latest` ]. URL restore source is expected to be in form of `protocol://authority/path`. Otherwise the value is interpreted as backupset name. The URL can point to a downloadable _tag.gz_ file containing a hotcopy backupset or an exact copy of a NuoDB archive (called _stream_). | `:latest` |
| `restore.credentials` | Credentials to use for a URL source (user:password) | `""` |
| `restore.stripLevels` | The number of levels to strip off path names when unpacking a TAR file of an archive or backup set | `"1"` |
| `restore.autoRestart` | Whether to automatically restart the database and trigger the restore (true/false). Only valid for a "database" restore | `true` |
| `restore.labels` | Process labels used to filter the complete set of archiveIds that should be restored which then defines the new state of the database upon restore. If multiple labels are defined, _any_ of them have to match so that the archive ID is selected for restore. By default all SM processes in the database are selected. The setting works with NuoDB 4.2+. | `""` |
| `restore.archiveIds` | Complete set of archiveIds that should be restored which then defines the new state of the database upon restore. Either `restore.labels` or `restore.archiveIds` should be specified. The setting works with NuoDB 4.2+ | `[]` |
| `restore.manual` | If set to "true", archives restore should be done manually once the SM pod is started. The engine startup will block waiting for the user to complete the archives restore. The setting works with NuoDB 4.2+ | `false` |

## Detailed Steps

- First install Helm charts for Admin and Database.

### Identify the Backupset

While installing the restore chart, we may wish to find the backupset to be used. This can be done before the existing database is shut down, or afterwards using a simple Pod. The following instructions assume the former approach.

#### While Database is Running

To get the recent backupset we need to exec into the `sm-database-cashews-demo-backup-0` pod:

```bash
$ kubectl exec -it sm-database-cashews-demo-backup-0 -- ls -al var/opt/nuodb/backup

total 32
drwxrwx--- 5 root  root  4096 Sep  4 16:03 .
drwxrwx--- 6 root  root  4096 Sep  4 15:57 ..
drwxr-xr-x 4 nuodb root  4096 Sep  4 16:02 20190904T160241
drwxr-xr-x 7 nuodb root  4096 Sep  4 16:05 20190904T160352
drwx------ 2 root  root 16384 Sep  4 15:57 lost+found
```

Then pick the desired recent backupset and copy its name into the source value in the file of restore; as the backupsets naturally sort, the last one is the latest. In the above example it would be: `20190710T190517`

In your values.yaml file, the place where you'd put the backupset name is as follows:

```bash
source: 20190904T160352
```

Or if you are using command line parameters, the setting would be:

```bash
... --set restore.source=20190904T160352 ...
```

In multi-cluster deployments, different backup cron jobs execute with different schedules.
Specific backupsets are produced for database deployment in each cluster.
To perform a database restore in multi-cluster deployment, select the hotcopy SMs in one of the clusters and provide a backupset available in that cluster or `:group-latest` as `restore.source`.

### Install Restore chart

If we have not set the values in `restore/values.yaml` then we can override it while installing the restore chart:

```bash
helm install nuodb/restore --name restore \
  ${values_option} \
  --set admin.domain=${DOMAIN_NAME} \
  --set restore.target=demoDb0 \
  --set restore.source=:latest \
  --set cloud.zones={us-central1-a}
```

The job should finish with a `Completed` status.

### Optionally manually restart the database

If `restore.autoRestart` was set to `true`, then the `restore` chart will restart the database, and the restore will proceed automatically.

However, if `restore.autoRestart` is set to `false`, then you retain control to manually stop and restart the pods you wish.

* You could shut down all processes in any order using `nuocmd shutdown database`. This will cause k8s to automatically restart all TE and SM pods in any order
* Alternatively, you could scale-down TE and SM pods, and then scale up the SM pods in the order of your choosing; and then scale-up the TE pods - again in the order of your choosing.

### Validate the restore

Verify the restore completed successfully; view the log output from the restarted SM pods, it should contain something similar to the following:

```bash
Finished restoring /var/opt/nuodb/backup/20190619T101450 to /var/opt/nuodb/archive/nuodb/demo. Created archive with archive ID 8
```

### Manual archive restore

For complex requirements where automatic archive restore is not sufficient, `restore.manual="true"` can be set during restore chart installation.
The NuoDB engines startup will block waiting for the user to perform the archives restore manually.
All archive ids requested for restore can be seen using:

```bash
nuodocker get restore-requests --db-name demo
```

Once the archive restore is done, it should be marked as completed using the following command:

```bash
nuodocker complete restore --db-name demo --archive-ids <archive_id>
```

The SM startup will unblock and the restored archives will be used to define the new database state.

## Uninstalling the Chart

To uninstall/delete the deployment:

```bash
helm del --purge restore
```

