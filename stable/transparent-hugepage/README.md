# Transparent Hugepage Helm Chart

This chart starts a Daemonset to disable transparent huge pages on the Linux hosts.

## Command

```bash
helm install nuodb/transparent-hugepage
```

## Prerequisites

- Kubernetes 1.9+

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

Busybox is used as the daemonset image that disables THP.

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

#### thp.*

The following tables list the configurable parameters of the transparent-hugepages chart and their default values.

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `nameOverride` | Overrides only the release name portion of the application name |`""`|
| `fullNameOverride` | Overrides the full application name |`""`|
| `affinity` | Container affinity | `{}` |
| `nodeSelector` | Container nodeSelector | `{}` |
| `tolerations` | Container tolerations | `[]` |

For example:

```yaml
thp:
  affinity: |
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: failure-domain.beta.kubernetes.io/zone
            operator: In
            values:
            {{- range .Values.cloud.zones }}
              - {{ . }}
            {{- end }}
```

### Running

Verify the Helm chart:

```bash
helm install nuodb/transparent-hugepage \
    --name transparent-hugepage \
    --debug --dry-run
```

Deploy the administration tier using volumes of the specified storage class:

```bash
helm install nuodb/transparent-hugepage \
    --name transparent-hugepage
```

Wait until the deployment completes:

```bash
helm status transparent-hugepage
```

Verify the pods are running:

```bash
$ kubectl get pods -l app=disable-thp-transparent-hugepage
NAME                      READY     STATUS    RESTARTS   AGE
disable-thp-transparent-hugepage-96tzq   1/1     Running   0          3m45s
disable-thp-transparent-hugepage-cbc8d   1/1     Running   0          3m45s
disable-thp-transparent-hugepage-cr8km   1/1     Running   0          3m45s
disable-thp-transparent-hugepage-g7bjb   1/1     Running   0          3m45s
```

## Manually Disabling Transparent Huge Pages

Run the `files/tuned.sh` script on each node to disable THP and set kernel parameters for NuoDB.
The file MUST be run as root/sudo.

## Uninstalling the Chart

To uninstall/delete the deployment:

```bash
helm del --purge transparent-hugepage
```

The command removes all the Kubernetes components associated with the chart and deletes the release.
