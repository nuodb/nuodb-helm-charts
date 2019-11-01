# NuoDB Admin Helm Chart

This chart starts a NuoDB admin deployment on a Kubernetes cluster using the Helm package manager.

## TL;DR;

```bash
helm install nuodb/admin
```

## Prerequisites

- Kubernetes 1.9+
- PV provisioner support in the underlying infrastructure (see `{provider}-storage.yaml`)

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

#### openshift.*

The purpose of this section is to specify parameters specific to OpenShift, e.g. enable features only present in OpenShift.

The following tables list the configurable parameters for the `openshift` option:

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `enabled` | Enable OpenShift features | `false` |
| `enableDeploymentConfigs` | Prefer DeploymentConfig over Deployment |`false`|
| `enableRoutes` | Enable OpenShift routes | `true` |

For example, to enable an OpenShift integration, and enable routes:

```yaml
openshift:
  enabled: true
  enableDeploymentConfigs: false
  enableRoutes: true
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
| `resources` | Labels to apply to all resources | `{}` |
| `affinity` | Affinity rules for NuoDB Admin | `{}` |
| `nodeSelector` | Node selector rules for NuoDB Admin | `{}` |
| `tolerations` | Tolerations for NuoDB Admin | `[]` |
| `configFilesPath` | Directory path where `configFiles.*` are found | `/etc/nuodb/` |
| `configFiles.*` | See below. | `{}` |
| `persistence.enabled` | Whether or not persistent storage is enabled for admin RAFT state | `false` |
| `persistence.accessModes` | Volume access modes enabled (must match capabilities of the storage class) | `ReadWriteMany` |
| `persistence.size` | Amount of disk space allocated for admin RAFT state | `10Gi` |
| `persistence.storageClass` | Storage class for volume backing admin RAFT state | `-` |
| `envFrom` | Import ENV vars from one or more configMaps | `[]` |
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

For example, when using GlusterFS storage class, you would supply the following parameter:

```bash
  ...
  --set admin.persistence.storageClass=glusterfs
  ...
```

#### admin.configFiles.*

The purpose of this section is to detail how to provide alternate configuration files for NuoDB. NuoDB has several configuration files that may be modified to suit.

There are two sets of configuration files documented:

- [Admin Configuration for a Particular Host][1]
- [Database Configuration for a Particular Host][2]

Any file located in `admin.configFilesPath` can be replaced; the YAML key corresponds to the file name being created or replaced.

The following tables list the configurable parameters for the `admin` option of the admin chart and their default values.

| Key | Description | Default |
| ----- | ----------- | ------ |
| `nuodb.lic` | [NuoDB license file content; defaults to NuoDB CE Edition][3] | `nil` |
| `nuoadmin.conf` | [NuoDB Admin host properties][4] | `nil` |
| `nuodb-types.config`| [Type mappings for the NuoDB Migrator tool][5] | `nil` |
| `nuoadmin.logback.xml` | Logging configuration. NuoDB recommends using the default settings. | `nil` |

### Running

Deploy storage classes and volumes (or suitable replacement):

```bash
kubectl create -f stable/admin/${cloud_provider}-storage.yaml
```

Verify the Helm chart:

```bash
helm install nuodb/admin -n admin \
    --set persistence.enabled=true \
    --set persistence.storageClass=nuodb-admin \
    --debug --dry-run
```

Deploy the administration tier using volumes of the specified storage class:

```bash
helm install nuodb/admin -n admin \
    --set persistence.enabled=true \
    --set persistence.storageClass=nuodb-admin
```

  **Tip**: Be patient, it will take approximately 55 seconds to deploy.

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
admin-nuodb-0                         1/1       Running   0          29m
tiller-86c4495fcc-lczdp        1/1       Running   0          5h
```

The command displays the NuoDB Pods running on the Kubernetes cluster. When completed, both the TE and the storage containers should show a **STATUS** of `Running`, and with 0 **RESTARTS**.

Verify the connected states of the database domain:

```bash
$ kubectl exec -it admin-nuodb-0 -- nuocmd show domain

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
kubectl scale sts admin-nuodb --replicas=3
```

## Uninstalling the Chart

To uninstall/delete the deployment:

```bash
helm del --purge admin
kubectl delete -f admin/${cloud_provider}-storage.yaml
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

[1]: #adminconfigfiles
[2]: /stable/database/README.md#databaseconfigfiles
[3]: http://doc.nuodb.com/Latest/Content/Nuoadmin-Obtaining-and-Installing-NuoDB-Licenses.htm
[4]: http://doc.nuodb.com/Latest/Content/Nuoadmin-Host-Properties.htm
[5]: http://doc.nuodb.com/Latest/Content/Data-Type-Mappings.htm
