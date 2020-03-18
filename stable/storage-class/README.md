# NuoDB Storage Class Helm Chart

This chart installs NuoDB storage classes in a Kubernetes cluster using the Helm package manager.

## Command

```bash
helm install nuodb/storage-class [--name releaseName] [--set parameter] [--values myvalues.yaml]
```

## Installing the Chart

### Configuration

All configurable parameters for each top-level scope are detailed below, organized by scope.

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

### Permissions

This chart installs storage classes required for the operation of NuoDB.
Since storage classes are cluster-scoped objects, in order to install the
chart, the user installing the chart must have cluster-role permissions.

There are two approaches to providing cluster-role permissions to Helm
in order to install the chart: more secure, and less secure approaches.
In one a separate Tiller server running in the kube-system namespace is
configured with cluster role permissions, and that role and Tiller server
is used to install the chart. In the second approach, the Tiller server
within the NuoDB namespace is configured with cluster role permissions.
The latter is a less secure approach as it grants a namespace-scoped
Tiller permissions to jailbreak out of the current namespace and affect
objects in other namespaces.

Both approaches are documented below.

#### Kube System Administrative Tiller Role

The service account and role below may be used to configure an administrative
role at cluster scope within the kube-system namespace to install the
chart.

```bash
kubectl -n kube-system create serviceaccount tiller-system
kubectl create clusterrolebinding tiller-system --clusterrole cluster-admin --serviceaccount=kube-system:tiller-system

helm init --service-account tiller-system --tiller-namespace kube-system
...
```

#### NuoDB Administrative Tiller Role

The service account and role below may be used to configure the NuoDB
namespace scoped Tiller server permissions to jailbreak out of the current
namespace and affect cluster-wide objects. Less secure, but indeed an
approach.

For example:

```yaml
---
kind: ServiceAccount
apiVersion: v1
metadata:
  name: tiller
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: tiller-storage-class
rules:
  - apiGroups: ["storage.k8s.io"]
    resources: ["storageclasses"]
    verbs: ["*"]
---
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: tiller-manager
rules:
  - apiGroups: ["", "batch", "extensions", "apps"]
    resources: ["*"]
    verbs: ["*"]
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: tiller-storage-class-binding
subjects:
- kind: ServiceAccount
  name: tiller
  namespace: nuodb
roleRef:
  kind: ClusterRole
  name: tiller-storage-class
  apiGroup: rbac.authorization.k8s.io
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: tiller-manager-binding
subjects:
- kind: ServiceAccount
  name: tiller
roleRef:
  kind: Role
  name: tiller-manager
  apiGroup: rbac.authorization.k8s.io
```

```bash
# Create resources in above template
kubectl -n nuodb create -f tiller.yaml

helm init --service-account tiller --tiller-namespace nuodb
...
```

### Running

Verify the Helm chart:

```bash
helm install nuodb/storage-class \
    --name storage-class \
    --debug --dry-run
```

Deploy the storage classes:

```bash
helm install nuodb/storage-class \
    --name storage-class
```

## Uninstalling the Chart

To uninstall/delete the deployment:

```bash
helm del --purge storage-class
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

[0]: #permissions
