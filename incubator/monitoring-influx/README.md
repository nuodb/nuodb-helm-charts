# NuoDB Influx Helm Chart

This chart deploys NuoDB Influx on a Kubernetes cluster using the Helm package manager.

## TL;DR;

```bash
helm install nuodb/monitoring-influx -n influx
```

## Prerequisites

- Kubernetes 1.9+
- PV provisioner support in the underlying infrastructure (see `{provider}-storage.yaml`)
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
| `enableRoutes` | Enable OpenShift routes | `true` |

For example, to enable an OpenShift integration, and enable routes:

```yaml
openshift:
  enabled: true
  enableRoutes: true
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

#### influx.*

The purpose of this section is to specify the Influx parameters.

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `image.registry` | Image repository where InfluxDB image is stored | `docker.io` |
| `image.repository` | Name of the Docker image | `influxdb` |
| `image.tag` | Tag for the InfluxDB Docker image | `1.6.0` |
| `image.pullPolicy` | Image pull policy | `Always` |
| `image.pullSecrets` | Specify docker-registry secret names as an array | [] (does not add image pull secrets to deployed pods) |
| `persistence.enabled` | Whether or not persistent storage is enabled for InfluxDB | `false` |
| `persistence.accessModes` | Volume access modes enabled (must match capabilities of the storage class) | `ReadWriteMany` |
| `persistence.size` | Amount of disk space allocated for InfluxDB | `20Gi` |
| `persistence.storageClass` | Storage class for volume backing InfluxDB | `-` |
| `service.type` | Influx service type | `ClusterIP` |
| `service.clusterIP` | Cluster IP (if any) | `nil | None` |
| `service.nodePort` | Specify the nodePort value for the LoadBalancer and NodePort service types. | `nil` |
| `service.externalIPs` | Specify the externalIP value ClusterIP service type. | `nil | []` |
| `service.loadBalancerIP` | Specify the loadBalancerIP value for LoadBalancer service types. | `nil` |
| `service.loadBalancerSourceRanges` | Specify the loadBalancerSourceRanges value for LoadBalancer service types. | `nil | []` |
|

For example, when using an internal Artifactory server, you would supply the following parameter:

```bash
  ...
  --set influx.registry=acme-dockerv2-virtual.jfrog.io
  ...
```

#### grafana.*

The purpose of this section is to specify the Grafana parameters.

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `image.registry` | Image repository where influxdb image is stored | `docker.io` |
| `image.repository` | Name of the Docker image | `grafana/grafana` |
| `image.tag` | Tag for the influxdb Docker image | `master` |
| `image.pullPolicy` | Image pull policy | `Always` |
| `image.pullSecrets` | Specify docker-registry secret names as an array | [] (does not add image pull secrets to deployed pods) |

For example, when using an internal Artifactory server, you would supply the following parameter:

```bash
  ...
  --set grafana.registry=acme-dockerv2-virtual.jfrog.io
  ...
```

Deploy storage classes and volumes (or suitable replacement):

```bash
kubectl create -f stable/monitoring-influx/${cloud_provider}-storage.yaml
```

Verify the Helm chart:

```bash
helm install nuodb/monitoring-influx -n monitoring-influx \
    --set influx.persistence.enabled=true \
    --set influx.persistence.storageClass=influx-data \
    --debug --dry-run
```

Deploy the InfluxDB-based monitoring solution:

```bash
helm install nuodb/monitoring-influx -n monitoring-influx \
    --set influx.persistence.enabled=true \
    --set influx.persistence.storageClass=influx-data \
    --debug --dry-run
```

  **Tip**: List all releases using `helm list`

Wait until the deployment completes:

Check on deployment status:

```bash
helm status monitoring-influx
```

### Connecting to Grafana

The following instructions detail how to connect to the Grafana dashboard.

1. Identify the Pod name:

    ```bash
    $ kubectl get pods | grep nuodb-dashboard-display | awk '{print $1}'
    nuodb-dashboard-display-6c5d6dd766-7rvsc   1/1   Running   0   7m
    ```

2. Port-forward to the Pod:

    ```bash
    $ kubectl -n nuodb port-forward `kubectl get pods | grep nuodb-dashboard-display | awk '{print $1}'` 3000
    Forwarding from 127.0.0.1:3000 -> 3000
    Forwarding from [::1]:3000 -> 3000
    ```

3. Open your browser to `http://localhost:3000/login`

4. Enter the following credentials: admin:nuodb

## Uninstalling the Chart

To uninstall/delete the deployment:

```bash
helm del --purge influx
```

The command removes all the Kubernetes components associated with the chart and deletes the release.
