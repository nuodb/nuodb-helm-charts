# NuoDB Restore Helm Chart

This chart starts a Job to restore a NuoDB database from an existing backup, in a Kubernetes cluster using the Helm package manager.

## TL;DR;

```bash
helm install nuodb/restore
```

## Prerequisites

- Kubernetes 1.9+
- An existing NuoDB Admin cluster has been provisioned

## Installing the Chart

### Configuration

The configuration is structured where configuration values are implemented following a single-definition rule, that is, values are structured and scoped, and shared across charts; e.g. for admin, its parameters are specified once in a single values file which is used for all the charts, and the database chart can use admin values for configuring connectivity of engines to a specific admin process. The same goes for other values **shared** amongst Helm charts. A few key points here:

- values files have structure, values are scoped
- different values files for different deployments
- values files follow the single definition rule (no repeats)
- global configuration exists under its own scoped section
- each chart has its own scoped section named after it
- cloud information is used to drive availability zones (particularly)

All configurable parameters for each top-level scope is detailed below, organized by scope.

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

#### admin.*

The following tables list the configurable parameters of the backup chart and their default values.

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `admin.domain` | NuoDB admin cluster name | `domain` |
| `admin.namespace` | OpenShift project where admin is deployed; when peering to an existing admin cluster provide its project name | `nuodb` |
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

#### restore.*

The following tables list the configurable parameters of the backup chart and their default values.

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `timeout` | Job deadline (timeout) | 1800 |
| `backupPvc` | PVC name for the backup volume | `nil` |
| `archivePvc` | PVC name for the archive volume | `nil` |
| `backupSet` | The backup set name to restore from | `nil` |

## Detailed Steps

- First install Helm charts for Admin, Database, YCSB, Influx and Backup.

### Save Existing PV/PVC Files

Save the `archive-volume-sm-database-domain-demo-backup-0` and `backup-volume-sm-database-nuodb-demo-backup-0` PVC and their attached PV of SM for further use.

First the archive PV/PVC files:

```bash
kubectl get pv $(kubectl get pv | grep archive-volume-sm-database-cashews-demo-backup-0 | awk '{ print $1 }') -o yaml > archive-pv.yaml
sed -i '' 's/name: .*$/name: archive-pvc/g' archive-pv.yaml
sed -i '' 's/volumeName: .*$/volumeName: archive-pv/g' archive-pv.yaml

kubectl get pvc archive-volume-sm-database-cashews-demo-backup-0 -o yaml > archive-pvc.yaml
sed -i '' 's/name: .*$/name: archive-pvc/g' archive-pvc.yaml
sed -i '' 's/volumeName: .*$/volumeName: archive-pv/g' archive-pvc.yaml
volumeName: pvc-affdb4ee-cf2c-11e9-be96-42010a800132
```

Then the backup PV/PVC files:

```bash
kubectl get pv $(kubectl get pv | grep backup-volume-sm-database-cashews-demo-backup-0 | awk '{ print $1 }') -o yaml > backup-pv.yaml
sed

kubectl get pvc backup-volume-sm-database-cashews-demo-backup-0 -o yaml > backup-pvc.yaml
sed
```

### Identify the Backupset

While installing the restore chart, we will need the backupset of most recent backup taken. This can be done before the existing database is shut down, or afterwards using a simple Pod. The following instructions assume the former approach.

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

Then pick the most recent backupset and copy it in the values file of restore; as the backupsets naturally sort, the last one is the latest. In the above example it would be: `20190710T190517`

In your values.yaml file, the place where you'd put the backupset name is as follows:

```bash
backupSet: 20190904T160352
```

Or if you are using command line parameters, the setting would be:

```bash
... --set restore.backupSet=20190904T160352 ...
```

### Set the Availability Zone in Values

Restore MUST be done in the same region as the backup job as the volumes cannot move between zones.

Therefore we need to have an affinity rule set up to select one node (`kubernetes.io/hostname`) in the zone, or specify the zone (`failure-domain.beta.kubernetes.io/zone`). But it must be constrained minimally by the zone.

You can only specify one zone in the zone-list in the `restore/values.yaml` file; this zone MUST match the above.

Set the zone in the values file of restore. This zone would be the zone mentioned in the PV attached to `archive-volume-sm-database-cashews-demo-backup-0` PVC; for example:

From the PV we find the zone:

```bash
$ cat backup-pv.yaml | grep failure-domain.beta.kubernetes.io/zone:
    failure-domain.beta.kubernetes.io/zone: us-central1-a

Then in the `values.yaml` file we would set the zone to:

```yaml
cloud:
  zones:
    - us-central1-a
```

Or if you are using command line parameters:

```bash
... --set cloud.zones={us-central1-a} ...
```

### Set the Backup and Archive PVC in Values

The backup and archive PVC need to be specified for the restore job.

Set these in the values file as follows (noting differences in the domain name as necessary):

```yaml
restore:

  backupPvc: backup-volume-sm-database-cashews-demo-backup-0

  archivePvc: archive-volume-sm-database-cashews-demo-backup-0
```

Or, if using command line parameters:

```bash
... --set restore.backupPvc=backup-volume-sm-database-cashews-demo-backup-0 --set restore.archivePvc=archive-volume-sm-database-cashews-demo-backup-0 ...
```

### Delete the Database and Related Workloads

```bash
helm del --purge demo-ycsb backup database
```

### Delete the Archives

To delete the archives at the admin layer, we need to exec into the admin pod and run the following commands:

```bash
kubectl exec -it admin-cashews-0 -- /bin/bash
$ nuocmd show archives
$ nuocmd delete database --db-name demo
```

Make sure the archives are clear if not we can delete the archives:

```bash
$ nuocmd delete archive --archive-id 0 --purge
$ nuocmd show archives
```

### Install Restore chart

If we have not set the values in `restore/values.yaml` then we can override it while installing the restore chart:

```bash
helm install nuodb/restore -n restore \
  ${values_option} \
  --set admin.domain=${DOMAIN_NAME} \
  --set restore.backupPvc=backup-volume-sm-database-cashews-demo-backup-0 \
  --set restore.archivePvc=archive-volume-sm-database-cashews-demo-backup-0 \
  --set cloud.zones={us-central1-a}
```

The job should finish with a `Completed` status.

### Validate the restore

Verify the restore completed successfully; view the log output from the restore Job pod, it should read similar to the following:

```bash
Finished restoring /var/opt/nuodb/backup/20190619T101450 to /var/opt/nuodb/archive/nuodb/demo. Created archive with archive ID 8
```

Verify the archive entry has been created:

```bash
$ kubectl exec -it admin-cashews-0 -- nuocmd show archives
[2] <NO VALUE> : /var/opt/nuodb/archive/cashews/demo @ demo [journal_path = ] [snapshot_archive_path = ] NOT_RUNNING
```

### Delete the Helm Restore Deployment and Validate Admin TOMBSTONE

First drop the restore deployment:

```bash
helm del --purge restore
```

Then validate the admin lists the database as being in a TOMBSTONE state:

```bash
kubectl exec -it admin-cashews-0 -- nuocmd show domain
```

### Update the PV/PVC Objects

Delete the archive-pv, backup-pv, archive-pvc and backup-pvc:

```bash
kubectl delete pvc archive-volume-sm-database-cashews-demo-0 backup-volume-sm-database-cashews-demo-backup-0
kubectl delete pv pvc-b0007157-cf2c-11e9-be96-42010a800132 pvc-b00fbb79-cf2c-11e9-be96-42010a800132
```

Add a new storage class for restored database to use:

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: nuodb-archive-from-snapshot
provisioner: kubernetes.io/azure-disk
parameters:
  storageaccounttype: Standard_LRS
  kind: Managed
  cachingmode: ReadWrite
reclaimPolicy: Retain
volumeBindingMode: Immediate
```

It is important that `volumeBindingMode: Immediate` is specified.

Update the archive-pv.yaml which was created in previous step and add label `database: demo` to it. Also change the storage class name as `nuodb-archive-from-snapshot` in the archive-pv.yaml.

Below is the example for `archive-pv.yaml`:

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  labels:
    failure-domain.beta.kubernetes.io/region: eastus
    failure-domain.beta.kubernetes.io/zone: '1'
    database: demo
  name: pvc-5c3532ef-a335-11e9-97e4-52126adb550c
spec:
  accessModes:
  - ReadWriteOnce
  azureDisk:
    cachingMode: ReadWrite
    diskName: kubernetes-dynamic-pvc-5c3532ef-a335-11e9-97e4-52126adb550c
    diskURI: /subscriptions/876c0bdf-be67-4156-877c-c4362fef9e92/resourceGroups/MC_helm-test_helm-test_eastus/providers/Microsoft.Compute/disks/kubernetes-dynamic-pvc-5c3532ef-a335-11e9-97e4-52126adb550c
    fsType: ""
    kind: Managed
    readOnly: false
  capacity:
    storage: 20Gi
  persistentVolumeReclaimPolicy: Retain
  storageClassName: nuodb-archive-from-snapshot
  ```

Be mindful of the zone the PV is located; use the zones setting to place the new database in the zone where the PV is located:

```yaml
labels:
     failure-domain.beta.kubernetes.io/region: eastus
     failure-domain.beta.kubernetes.io/zone: '1
```

Apply the `archive-pv.yaml`. Note: Only PV needs to be applied and not the pvc as it will be created by the new database.

```bash
kubectl apply -f archive-pv.yaml
```

### Start a new Database

Start a restored database:

```bash
helm install nuodb/database -n restored-database \
  ${values_option} \
  --set admin.domain=${DOMAIN_NAME} \
  --set database.persistence.storageClass=nuodb-archive-from-snapshot \
  --set cloud.zones={us-central1-a} \
  --set database.isManualVolumeProvisioning=true \
  --set backup.persistence.enabled=false
```

Validate the PV if it got bound to the PVC:

```bash
kubectl describe pv pvc-f91ce1b9-a255-11e9-90fe-42010a800179
```

Check the domain and the archives:

```bash
kubectl exec -it admin-cashews-0 -- nuocmd show domain
kubectl exec -it admin-cashews-0 -- nuocmd show archives
```

**Tip**: List all releases using `helm list`

Wait until the deployment completes:

## Uninstalling the Chart

To uninstall/delete the deployment:

```bash
helm del --purge restore
kubectl delete jobs --all
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## References

1. Kubernetes Concepts: [Persistent Volumes][0]
2. Kubernetes How To: [How to Use Preexisting Persistent Disks][1]

[0]: https://cloud.google.com/kubernetes-engine/docs/concepts/persistent-volumes
[1]: https://cloud.google.com/kubernetes-engine/docs/how-to/persistent-volumes/preexisting-pd
[2]: images/restore-workflow.png
[3]: https://github.com/kubernetes/examples/blob/master/staging/volumes/azure_disk/README.md
[4]: https://github.com/kubernetes/examples/blob/master/staging/volumes/azure_disk/azure.yaml
[5]: https://kubernetes.io/docs/concepts/storage/storage-classes/#aws-ebs
[6]: https://kubernetes.io/docs/concepts/storage/storage-classes/#azure-disk
[7]: https://kubernetes.io/docs/concepts/storage/storage-classes/#gce-pd
[8]: https://kubernetes.io/docs/concepts/storage/volumes/#awselasticblockstore
[9]: https://kubernetes.io/docs/concepts/storage/volumes/#azuredisk
[10]: https://kubernetes.io/docs/concepts/storage/volumes/#gcepersistentdisk
[11]: https://github.com/MicrosoftDocs/azure-docs/blob/master/articles/aks/azure-disks-dynamic-pv.md

[internal-0]: http://confluence.internal.nuodb.com/display/EN2/Restore+in+Openshift
[internal-1]: http://confluence.internal.nuodb.com/pages/viewpage.action?pageId=34079913
