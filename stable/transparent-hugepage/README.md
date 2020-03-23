# Transparent Hugepage Helm Chart

This chart starts a Daemonset to disable transparent huge pages (THP) on the Linux hosts.

## Command

```bash
helm install nuodb/transparent-hugepage [--name releaseName] [--set parameter] [--values myvalues.yaml]
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
    --name thp \
    --debug --dry-run
```

Run the Helm Chart:

```bash
helm install nuodb/transparent-hugepage \
    --name thp
```

Check the deployment status:

```bash
helm status thp
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
helm del --purge thp
```

The command removes all the Kubernetes components associated with the chart and deletes the release.
