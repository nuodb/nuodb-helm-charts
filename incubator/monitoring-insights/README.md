# NuoDB Insights Helm Chart

This chart deploys NuoDB Insights on a Kubernetes cluster using the Helm package manager.

## Command

```bash
helm install nuodb/monitoring-insights [--name releaseName] [--set parameter] [--values myvalues.yaml]
```

## Software Version Prerequisites

Please visit the **[NuoDB Helm Chart main page](https://github.com/nuodb/nuodb-helm-charts/#software-release-requirements)** for software version prerequisites.

## Installing the Chart

### Configuration

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

#### insights.*

The purpose of this section is to specify the Insights parameters.

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `enabled` | Use to control Insights Opt In.  Insights provides database monitoring. | `true` |
| `affinity` | Affinity rules for Insights | `{}` |
| `nodeSelector` | Node selector rules for Insights | `{}` |
| `tolerations` | Tolerations for Insights | `[]` |

Verify the Helm chart:

```bash
helm install nuodb/monitoring-insights --name insights --debug --dry-run
```

Deploy Insights:

```bash
helm install nuodb/monitoring-insights --name insights
```

The command deploys NuoDB Insights on the Kubernetes cluster in the default configuration. The configuration section lists the parameters that can be configured during installation.

  **Tip**: List all releases using `helm list`

Wait until the deployment completes:

```bash
helm status monitoring-insights
```

Verify the pods are running:

```bash
$ kubectl get pods
NAME                           READY     STATUS             RESTARTS   AGE
nuodb-insights                 2/2       Running            0          33m
```

  **Tip**: Wait until all processes are be in a **RUNNING** state.

The command displays the NuoDB Pods running on the Kubernetes cluster. When completed, the nuodb-insights containers should show a **STATUS** of `Running`, and with 0 **RESTARTS**.

You need to determine the generated NuoDB Insights URL, there are several ways to do this. The easiest being to look at the top of the Insights log for the `sub_id`:

```bash
$ kubectl exec -it nuodb-insights -c insights -- bash
# then once inside the container...
> grep sub_id /var/log/nuodb/nuoca.log
2019-04-10 21:10:18,011 NuoCA INFO Output Key: 'sub_id' set to Value: 'REDACTED_SUB_ID'
```

Another way to get your subscription ID is as follows:

```bash
$ kubectl exec -it nuodb-insights -c insights -- nuoca check insights
{"insights.sub.dashboard_url": "https://insights.nuodb.com/REDACTED_SUB_ID/", "insights.sub.id": "REDACTED_SUB_ID", "insights.sub.token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWJzY3JpYmVyX2lkIjoiOTdaRDRWUjUzMSJ9.ZecY9ZcXOSdQLIMBzkbJtfANJM4x8aNDCEMLxm09tCPf_xE9abuXS1k6EY_TGBZZMb04iEDjn1y2VPxTFiiDS7ZlP1-w_hExbhAGW-WY3oI9lnpUjMdSPGdipQEfnndZtt9wPzJLmFMzCWfpi0HUskB6sR6ywD2FgMnQy6wyWgRsIQaBp3VtaJ7e3kqqJbUJSAmgU11RMzg12u879R5ESFSl4vMRPvdvRQDyUGVsCkDml9cmDix6ZfmHKVUc4rZk4Z4FOAGuhjdsXJ-Rw_nd_6CV1M9gA8sw5TuUpZzUWm6IN35G1rsYZYZ9RYLcthI7YvbFxG5XYp4Zhgwn7zUzkg", "insights.enabled": "True", "insights.sub.ingest_url": "https://insights.nuodb.com/ingest", "insights.sub.api_url": "https://insights.nuodb.com/api/1"}
```

The command above will output the `sub_id`, which becomes part of your Insights URL:

  <https://insights.nuodb.com/REDACTED_SUB_ID_HERE/>

Open your browser to the URL, substituting in your sub_id into the path, to show
live performance metrics in its Grafana dashboards.

## Uninstalling the Chart

To uninstall/delete the deployment:

```bash
helm del --purge insights
```

The command removes all the Kubernetes components associated with the chart and deletes the release.
