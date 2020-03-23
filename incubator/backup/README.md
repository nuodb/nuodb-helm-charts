# NuoDB Backup Helm Chart

This chart starts a NuoDB database backup on existing NuoDB storage managers in a Kubernetes cluster using the Helm package manager.

## Command

```bash
helm install nuodb/backup [--name releaseName] [--set parameter] [--values myvalues.yaml]
```

## Software Version Prerequisites

Please visit the **[NuoDB Helm Chart main page](https://github.com/nuodb/nuodb-helm-charts/#software-release-requirements)** for software version prerequisites.

## Installing the Chart

All values.yaml configurable parameters for each top-level scope are detailed below, organized by scope.

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

#### admin.*

The purpose of this section is to specify the NuoDB Admin parameters.

The following tables list the configurable parameters for the `admin` option of the admin chart and their default values.

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `domain` | NuoDB admin cluster name | `nuodb` |
| `namespace` | Namespace where admin is deployed; when peering to an existing admin cluster provide its project name | `nuodb` |
| `tlsCACert.secret` | TLS CA certificate secret name | `nil` |
| `tlsCACert.key` | TLS CA certificate secret key | `nil` |
| `tlsClientPEM.secret` | TLS client PEM secret name | `nil` |
| `tlsClientPEM.key` | TLS client PEM secret key | `nil` |

For example:

```yaml
admin:
  domain: nuodb
```

#### database.*

The purpose of this section is to specify the NuoDB Admin parameters.

The following tables list the configurable parameters for the `admin` option of the admin chart and their default values.

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `name` | Database name |`demo`|

For example, to enable an OpenShift integration, and enable routes:

```yaml
database:
  name: demo
```

#### backup.*

The purpose of this section is to specify the NuoDB backup parameters.

The following tables list the configurable parameters for the `backup` option:

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `fullSchedule` | Cron schedule for full backups | `35 22 */1 * 6` |
| `incrementalSchedule` | Cron schedule for incremental backups |`35 22 */1 * 0-5`|
| `deadline` | Deadline for a job in seconds | `1800` |
| `successHistory` | how many completed jobs should be kept |`5`|
| `failureHistory` | how many failed jobs should be kept |`3`|
| `timeout` | Nuo Docker timeout (should match deadline) |`1800`|
| `backupDir` | Directory location for backups |`/var/opt/nuodb/backup`|

See the following online resource for an interactive tool to create cron schedules:

  [Crontab Guru Tool](https://crontab.guru/)

### Running

Verify the Helm chart:

```bash
helm install nuodb/backup --name backup \
    --debug --dry-run
```

Run a full backup:

```bash
job_name=full-job-$(cat /dev/urandom | tr -dc 'a-z0-9' | fold -w 8 | head -n 1)
kubectl create job $job_name \
  --from=cronjob/full-backup-demo-cronjob
```

Run a incremental backup:

```bash
job_name=incr-job-$(cat /dev/urandom | tr -dc 'a-z0-9' | fold -w 8 | head -n 1)
kubectl create job $job_name \
  --from=cronjob/incremental-backup-demo-cronjob
```

The command deploys NuoDB on the Kubernetes cluster in the default configuration. The configuration section lists the parameters that can be configured during installation.

  **Tip**: List all releases using `helm list`

Wait until the deployment completes:

```bash
helm status backup
```

Verify the pods are running:

```bash
$ kubectl get pods
NAME                           READY     STATUS    RESTARTS   AGE
demo-storage-wo-backup-q7n6n   1/1       Running   0          9m
east-0                         1/1       Running   0          29m
te-east-1-l62jv                1/1       Running   0          9m
tiller-86c4495fcc-lczdp        1/1       Running   0          5h
```

The command displays the NuoDB Pods running on the Kubernetes cluster. When completed, both the TE and the storage containers should show a **STATUS** of `Running`, and with 0 **RESTARTS**.

## Validate the backup

Check for `.inc` and `full` which means backup is successful and also for completed pods in the console:

```bash
kubectl exec -it  sm-demo-backup-0 -- /bin/bash
[root@sm-demo-backup-0 /]# cd /var/opt/nuodb/backup
[root@sm-demo-backup-0 backup]# ls
20190605T200147  lost+found
[root@sm-demo-backup-0 backup]# ls 20190605T200147
1.inc  2.inc  3.inc  4.inc  full  state.xml  tmp
```

The number of `.inc` files depends on how many times you trigger the incremental backup from the console.

## Uninstalling the Chart

To uninstall/delete the deployment:

```bash
helm del --purge backup
kubectl delete jobs --all
```

The command removes all the Kubernetes components associated with the chart and deletes the release.
