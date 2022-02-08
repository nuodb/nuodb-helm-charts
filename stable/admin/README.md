# NuoDB Admin Helm Chart

This chart starts a NuoDB Admin deployment on a Kubernetes cluster using the Helm package manager and must be running before attempting to start a NuoDB database.

## Command

```bash
helm install [name] nuodb/admin [--generate-name] [--set parameter] [--values myvalues.yaml]
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
| `addServiceAccount` | Whether to create a new service account for NuoDB containers | `true` |
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
| `lbConfig.prefilter` | Global load balancer prefilter expression | `nil` |
| `lbConfig.default` | Global load balancer default query | `nil` |
| `lbConfig.policies` | Load balancer named policies | `{ nearest: ... }` |
| `lbConfig.fullSync` | Remove any manually created load balancer configuration | `false` |
| `lbPolicy` | Load balancer policy name | `nearest` |
| `lbQuery` | Load balancer query | `random(first(label(pod ${pod:-}) label(node ${node:-}) label(zone ${zone:-}) any))` |
| `externalAccess.enabled` | Whether to deploy a Layer 4 cloud load balancer service for the admin layer | `false` |
| `externalAccess.internalIP` | Whether to use an internal (to the cloud) or external (public) IP address for the load balancer | `nil` |
| `externalAccess.type` | The service type used to enable external access for NuoDB Admin. The supported types are `NodePort` and `LoadBalancer` (defaults to `LoadBalancer`) | `nil` |
| `externalAccess.annotations` | Annotations to pass through to the Service of type `LoadBalancer` | `{}` |
| `ingress.enabled` | Whether to deploy an Ingress resources for the NuoDB Admin | `false` |
| `ingress.api.hostname` | The fully qualified domain name (FQDN) of the network host for the NuoDB Admin REST API | `""` |
| `ingress.api.path` | Path that is matched against the path of the incoming HTTP request | `/` |
| `ingress.api.className` | The associated IngressClass name defines which Ingress controller will implement the resource | `""` |
| `ingress.api.annotations` | Custom annotations that are set on the Ingress resource | `{ ingress.kubernetes.io/ssl-passthrough: "true" }` |
| `ingress.api.tls` | Enable TLS termination in the Ingress controller. It is recommended to use SSL Passthrough feature instead | `false` |
| `ingress.api.secretName` | The name of the secret used by Ingress controller to terminate the TLS traffic | `""` |
| `ingress.sql.hostname` | The fully qualified domain name (FQDN) of the network host used by SQL clients | `""` |
| `ingress.sql.className` | The associated IngressClass name defines which Ingress controller will implement the resource | `""` |
| `ingress.sql.annotations` | Custom annotations that are set on the Ingress resource | `{ ingress.kubernetes.io/ssl-passthrough: "true" }` |
| `ingress.sql.tls` | Enable TLS termination in the Ingress controller. It is recommended to use SSL Passthrough feature instead | `false` |
| `ingress.sql.secretName` | The name of the secret used by Ingress controller to terminate the TLS traffic | `""` |
| `resources` | Kubernetes resource requests and limits used for the NuoDB Admin containers | `{}` |
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
| `options` | Set options to be passed to nuoadmin as arguments | `{}` |
| `initContainers.runInitDisk` | Whether to run the `init-disk` init container to set volume permissions | `true` |
| `initContainers.runInitDiskAsRoot` | Whether to run the `init-disk` init container as root | `true` |
| `securityContext.fsGroupOnly` | Creates a security context for Pods containing only the `securityContext.fsGroup` value | `false` |
| `securityContext.runAsNonRootGroup` | Creates a security context for Pods containing a non-root user and group (1000:1000) along with the `securityContext.fsGroup` value | `false` |
| `securityContext.enabled` | Creates a security context for Pods containing the `securityContext.runAsUser` and `securityContext.fsGroup` values | `false` |
| `securityContext.runAsUser` | The user ID for the Pod security context created if `securityContext.enabled` is `true`. | `1000` |
| `securityContext.fsGroup` | The `fsGroup` for the Pod security context created if any of `securityContext.fsGroupOnly`, `securityContext.runAsNonRootGroup`, or `securityContext.enabled` are `true`. | `1000` |
| `securityContext.enabledOnContainer` | Whether to create SecurityContext for containers | `false` |
| `securityContext.capabilities` | Capabilities for to admin container security context | `{ add: [], drop: [] }` |
| `securityContext.privileged` | Run the NuoDB Admin containers in privileged mode. Processes in privileged containers are essentially equivalent to root on the host | `false` |
| `securityContext.allowPrivilegeEscalation` | Whether a process can gain more privileges than its parent process. This boolean directly controls if the `no_new_privs` flag will be set on the container process | `false` |
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
| `serviceSuffix.nodeport` | The suffix to use for the NodePort service name | `nodeport` |
| `readinessTimeoutSeconds` | Admin readiness probe timeout, sometimes needs adjusting depending on environment and pod resources | `1` |
| `podAnnotations` | Annotations to pass through to the Admin pod | `nil` |
| `tde.secrets` | Transparent Data Encryption secret names used for different databases | `{}` |
| `tde.storagePasswordsDir` | Transparent Data Encryption storage passwords mount path | `/etc/nuodb/tde` |
| `evicted.servers` | A list of evicted servers excluded from RAFT consensus. Used during disaster recovery. | `[]` |


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
| `nodeport` | suffix for the NodePort service | "nodeport" |

#### admin.legacy

Features in this section have been deprecated but not yet removed.

| Key | Description | Default |
| ----- | ----------- | ------ |
| `loadBalancerJob.enabled` | Create a job that sets the default load balancer policy for the admin tier. Replaced by Kubernetes Aware Admin. | false |


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
| `plugins.admin` | NuoDB Collector additional plugins for admin services |`{}`|
| `resources` | Kubernetes resource requests and limits used for the nuocollector sidecar |`{}`|

### Running

Verify the Helm chart:

```bash
helm install admin nuodb/admin --debug --dry-run
```

Deploy the administration tier:

**Tip**: If you plan to deploy NuoDB Insights visual monitoring, add the `--set nuocollector.enabled=true` switch as below.


```bash
helm install admin nuodb/admin --set nuocollector.enabled=true
```

The command deploys NuoDB on the Kubernetes cluster using the default configuration. The configuration section lists the parameters that can be configured during installation.

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
```

The command displays the NuoDB Pods running on the Kubernetes cluster. When completed, both the TE and the storage containers should show a **STATUS** of `Running`, and with 0 **RESTARTS**.

Verify the connected states of the database domain:

```bash
$ kubectl exec -it admin-nuodb-cluster0-0 -- nuocmd show domain

server version: 4.0.7-2-6526a2db74, server license: Community
server time: 2020-12-10T21:13:25.722, client token: e64322a4728c8bf35ff7f02cda62cc74aca40b66
Servers:
  [admin-nuodb-cluster0-0] admin-nuodb-cluster0-0.nuodb.nuodb-helm.svc.cluster.local:48005 
     (LEADER, Leader=admin-nuodb-cluster0-0, log=0/35/35) Connected *
Databases:
```

The Servers section lists admin processes; each admin server should transition to the **Connected** state. When multiple Admins are started, one will be the **LEADER** and other designated as a **FOLLOWER**.

### Scaling

To scale the admin to 3 replicas, e.g., run the following command:

```bash
kubectl scale sts admin-nuodb-cluster0 --replicas=3
```

## Uninstalling the Chart

To uninstall/delete the deployment:

```bash
helm delete admin
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

[3]: http://doc.nuodb.com/Latest/Content/Nuoadmin-Obtaining-and-Installing-NuoDB-Licenses.htm
[4]: http://doc.nuodb.com/Latest/Content/Nuoadmin-Host-Properties.htm
[5]: http://doc.nuodb.com/Latest/Content/Data-Type-Mappings.htm
