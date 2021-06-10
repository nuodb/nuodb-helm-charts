# HOCKEY Demo Helm Chart

This chart deploys the NuoDB HOCKEY Demo on a Kubernetes cluster using the Helm package manager.

## Command

```bash
helm install [name] nuodb-incubator/demo-hockey [--generate-name] [--set parameter] [--values myvalues.yaml]
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

#### admin.*

The purpose of this section is to specify the NuoDB Admin parameters.

The following tables list the configurable parameters for the `admin` option of the admin chart and their default values.

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `domain` | NuoDB admin cluster name | `nuodb` |
| `namespace` | Namespace where admin is deployed; when peering to an existing admin cluster provide its project name | `nuodb` |

For example, to enable an OpenShift integration, and enable routes:

```yaml
admin:
  domain: nuodb
```

#### database.*

The following tables list the configurable parameters of the database chart and their default values.

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `name` | Database name | `demo` |

#### hockey.*

The following tables list the configurable parameters of the HOCKEY chart and their default values.

| Parameter | Description | Default |
| ----- | ----------- | ------ |
| `nameOverride` | Database name | `demo` |
| `fullnameOverride` | Number of threads | `1` |
| `replicas` | Number of NuoDB HOCKEY replicas | `0` |
| `loadName` | Name of the activity | `hockey-load` |
| `maxDelay` | HOCKEY maximum workload delay in milliseconds (Default is 4 minutes) | `240000` |
| `dbSchema` | HOCKEY database schema | `USER` |
| `image.registry` | Image repository where NuoDB image is stored | `docker.io` |
| `image.repository` | Name of the Docker image | `nuodb/hockey` |
| `image.tag` | Tag for the NuoDB Docker image | `latest` |
| `image.pullPolicy` | Image pull policy | `Always` |
| `image.pullSecrets` | Specify docker-registry secret names as an array | [] (does not add image pull secrets to deployed pods) |
| `affinity` | Affinity rules for NuoDB HOCKEY pods | `{}` |
| `nodeSelector` | Node selector rules for NuoDB HOCKEY pods | `{}` |
| `tolerations` | Tolerations for NuoDB HOCKEY pods | `[]` |

Verify the Helm chart:

```bash
helm install hockey nuodb-incubator/demo-hockey --debug --dry-run
```

Deploy the demo:

```bash
helm install hockey nuodb-incubator/demo-hockey
```

The command deploys NuoDB Hockey Quickstart on the Kubernetes cluster in the default configuration. The configuration section lists the parameters that can be configured during installation.

  **Tip**: List all releases using `helm list`

Wait until the deployment completes:

```bash
helm status hockey
```

Verify the pods are running:

```bash
$ kubectl get pods
NAME                     READY   STATUS      RESTARTS   AGE
hockey-load-xcl5f          1/1     Running     0          18h
```

  **Tip**: Wait until all processes are be in a **RUNNING** state.

Now to scale the Quickstart workload is simple enough:

```bash
$ kubectl scale rc hockey-load --replicas=1
replicationcontroller "demo-hockey" scaled
```

## Uninstalling the Chart

To uninstall/delete the deployment:

```bash
helm delete hockey
```

The command removes all the Kubernetes components associated with the chart and deletes the release.
