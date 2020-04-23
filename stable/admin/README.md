# NuoDB Admin Helm Chart

This chart starts a NuoDB Admin deployment on a Kubernetes cluster using the Helm package manager and must be running before attempting to start a NuoDB database.

## Command

```bash
helm install nuodb/admin [--name releaseName] [--set parameter] [--values myvalues.yaml]
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
| `cluster.domain` | Kubernetes domain as specified during kubernetes set up. Useful in multi-cluster configs | `cluster.local` |
| `cluster.entrypointDomain` | Kubernetes domain for the NuoDB Entrypoint Admin Process. Useful in multi-cluster configs | `cluster.local` |


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
    pullPolicy: IfNotPresent
```

#### nuodb.*

The purpose of this section is to specify the NuoDB image parameters.

The following tables list the configurable parameters for the `nuodb` option:

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `image.registry` | NuoDB container registry | `docker.io` |
| `image.repository` | NuoDB container image name |`nuodb/nuodb-ce`|
| `image.tag` | NuoDB container image tag | `latest` |
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

The following tables list the configurable parameters for the `admin` option of the admin chart and their default values.

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `nameOverride` | Override the name of the chart |`nil`|
| `fullNameOverride` | Override the fullname of the chart |`nil`|
| `domain` | NuoDB admin cluster name | `nuodb` |
| `namespace` | Namespace where admin is deployed; when peering to an existing admin cluster provide its project name | `nuodb` |
| `replicas` | Number of NuoDB Admin replicas | `1` |
| `lbPolicy` | Load balancer policy name | `nil` |
| `lbQuery` | Load balancer query | `nil` |
| `externalAccess.enabled` | Whether to deploy a Layer 4 cloud load balancer service for the admin layer | `false` |
| `externalAccess.internalIP` | Whether to use an internal (to the cloud) or external (public) IP address for the load balancer | `nil` |
| `resources` | Labels to apply to all resources | `{}` |
| `affinity` | Affinity rules for NuoDB Admin | `{}` |
| `nodeSelector` | Node selector rules for NuoDB Admin | `{}` |
| `tolerations` | Tolerations for NuoDB Admin | `[]` |
| `configFilesPath` | Directory path where `configFiles.*` are found | `/etc/nuodb/` |
| `configFiles.*` | See below. | `{}` |
| `persistence.accessModes` | Volume access modes enabled (must match capabilities of the storage class) | `ReadWriteMany` |
| `persistence.size` | Amount of disk space allocated for admin RAFT state | `10Gi` |
| `persistence.storageClass` | Storage class for volume backing admin RAFT state | `-` |
| `logPersistence.enabled` | Whether to enable persistent storage for logs | `false` |
| `logPersistence.overwriteBackoff.copies` | How many copies of the crash directory to keep within windowMinutes | `3` |
| `logPersistence.overwriteBackoff.windowMinutes` | The window within which to keep the number of crash copies | `120` |
| `logPersistence.accessModes` | Volume access modes enabled (must match capabilities of the storage class) | `ReadWriteOnce` |
| `logPersistence.size` | Amount of disk space allocated for log storage | `60Gi` |
| `logPersistence.storageClass` | Storage class for volume backing log storage | `-` |
| `envFrom` | Import ENV vars from one or more configMaps | `[]` |
| `options` | Set optons to be passed to nuoadmin as arguments | `{}` |
| `securityContext.capabilities` | add capabilities to the container | `[]` |
| `tlsCACert.secret` | TLS CA certificate secret name | `nil` |
| `tlsCACert.key` | TLS CA certificate secret key | `nil` |
| `tlsKeyStore.secret` | TLS keystore secret name | `nil` |
| `tlsKeyStore.key` | TLS keystore secret key | `nil` |
| `tlsKeyStore.password` | TLS keystore secret password | `nil` |
| `tlsTrustStore.secret` | TLS truststore secret name | `nil` |
| `tlsTrustStore.key` | TLS truststore secret key | `nil` |
| `tlsTrustStore.password` | TLS truststore secret password | `nil` |
| `tlsClientPEM.secret` | TLS client PEM secret name | `nil` |
| `tlsClientPEM.key` | TLS client PEM secret key | `nil` |
| `serviceSuffix.balancer` | The suffix to use for the LoadBalancer service name | `balancer` |
| `serviceSuffix.clusterip` | The suffix to use for the ClusterIP service name | `clusterip` |

For example, when using GlusterFS storage class, you would supply the following parameter:

```bash
  ...
  --set admin.persistence.storageClass=glusterfs
  ...
```

#### admin.configFiles.*

The purpose of this section is to detail how to provide alternate Admin configuration files for NuoDB. 

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


### Running

Verify the Helm chart:

```bash
helm install nuodb/admin --name admin \
    --debug --dry-run
```

Deploy the administration tier:

```bash
helm install nuodb/admin --name admin
```

**Tip**: It will take approximately 1 minute to deploy.

The command deploys NuoDB on the Kubernetes cluster in the default configuration. The configuration section lists the parameters that can be configured during installation.

**Tip**: List all releases using `helm list`

Wait until the deployment completes:

```bash
helm status admin
```

Verify the pods are running:

```bash
$ kubectl get pods
NAME                           READY     STATUS    RESTARTS   AGE
admin-nuodb-cluster0-0         1/1       Running   0          29m
tiller-86c4495fcc-lczdp        1/1       Running   0          5h
```

The command displays the NuoDB Pods running on the Kubernetes cluster. When completed, both the TE and the storage containers should show a **STATUS** of `Running`, and with 0 **RESTARTS**.

Verify the connected states of the database domain:

```bash
$ kubectl exec -it admin-nuodb-cluster0-0 -- nuocmd show domain

server version: 3.4.1-1-ccb6be381c, server license: Community
server time: 2019-04-10T00:25:53.054, client token: 370d671ff18dd57a4b4bb0d146c72c8f2f256e7f
Servers:
  [east-0] east-0.domain.nuodb.svc:48005 [last_ack = 6.92] [member = ADDED] [raft_state = ACTIVE] (LEADER, Leader=east-0, log=0/15/15) Connected *
Databases:
```

The command displays the status of NuoDB processes. The Servers section lists admin processes; they should all be **Connected**, one will be the **LEADER** and other designated as a **FOLLOWER**.

  **Tip**: Wait until all processes are be in a **RUNNING** state.

### Scaling

To scale the admin to 3 replicas, e.g., run the following command:

```bash
kubectl scale sts admin-nuodb-cluster0 --replicas=3
```

## Uninstalling the Chart

To uninstall/delete the deployment:

```bash
helm del --purge admin
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

[3]: http://doc.nuodb.com/Latest/Content/Nuoadmin-Obtaining-and-Installing-NuoDB-Licenses.htm
[4]: http://doc.nuodb.com/Latest/Content/Nuoadmin-Host-Properties.htm
[5]: http://doc.nuodb.com/Latest/Content/Data-Type-Mappings.htm
