# Getting Started with NuoDB Helm Charts

This section will walk you through getting both the Tiller server and Helm client up and running in your Kubernetes cluster. If using Red Hat OpenShift, this page assumes that you already have the OpenShift `oc` client program installed locally and that you are logged into your OpenShift instance.

The instructions are in two parts:

1. **[Getting Started with Helm][4]** covers how to install and configure Helm on a client host. It will walk you through deploying a canary application to make sure Helm is properly configured.
2. **[Deploying NuoDB using Helm Charts][5]** covers how to configure hosts to permit running NuoDB, and covers deploying your first NuoDB database using the provided Helm charts.

There are sub-charts in subdirectories included in this distribution. Instructions provided on this page are for initial configuration of Helm and Tiller, in some cases, required security settings. Sub-charts pages include instructions for deploying each required NuoDB component.

## Version Prerequisites

Please visit the **[NuoDB Helm Chart main page][6]** for software version pre-requisites.


## Install Helm and Tiller on your Client

There are two parts to Helm: The Helm client (`helm`) and the Helm server (Tiller). This guide shows how to [install the client][1], and then proceeds to show two ways to install the server.

Every [release][2] of Helm provides binary releases for a variety of OSes. These binary versions can be manually downloaded and installed.

1. Download your [desired version][2]
2. Unpack it (`tar -zxvf helm-${helm-version}-linux-amd64.tgz`)
3. Find the helm binary in the unpacked directory, and move it to its desired destination (`mv linux-amd64/helm /usr/local/bin/helm`)

From there, you should be able to run the client: `helm help`.

We’ll use Helm version 2.16.1, which can be downloaded via <https://github.com/kubernetes/helm/releases/tag/v2.16.1.>

Run the following commands to install the Helm client:

```bash
$ curl -s https://storage.googleapis.com/kubernetes-helm/helm-v2.16.1-linux-amd64.tar.gz | tar xz
$ cd linux-amd64
$ mv helm /usr/local/bin
$ mv tiller /usr/local/bin
```

If you're running on Mac, the curl commands above would be:

```bash
$ curl -s https://storage.googleapis.com/kubernetes-helm/helm-v2.16.1-darwin-amd64.tar.gz | tar xz
$ cd darwin-amd64
$ mv helm /usr/local/bin
$ mv tiller /usr/local/bin
```

To keep isolated between individuals sharing a login, you can have separate Kubernetes, Helm, and Tiller state. Run these before proceeding.

```bash
export TILLER_NAMESPACE=nuodb
export HELM_HOME=`pwd`/.helm
export KUBECONFIG=`pwd`/.kube/config
```

The following commands will configure the current users environment for Helm; it will create a `.helm` directory in `${HOME}`.

```bash
$ helm init --client-only
$HELM_HOME has been configured at /home/clusteradmin/.helm.
Not installing Tiller due to 'client-only' flag having been set
Happy Helming!
```

### Install the Tiller Server in Kubernetes

We will be creating the Tiller server in the `kube-system` namespace so that it is available to all projects.

Create a new service account for tiller.
```bash
kubectl -n kube-system create serviceaccount tiller-system
```

Give `cluster-admin` permissions to the newly created service account.
```bash
kubectl create clusterrolebinding tiller-system \
--clusterrole cluster-admin \
--serviceaccount=kube-system:tiller-system
```

### Start the Tiller server
```bash
helm init --service-account tiller-system --tiller-namespace kube-system
```
Check that the Helm client and Tiller server are able to communicate correctly by running helm version. The results should be as follows:

```bash
$ helm version
Client: &version.Version{SemVer:"v2.14.1", GitCommit:"618447cbf203d147601b4b9bd7f8c37a5d39fbb4", GitTreeState:"clean"}
Server: &version.Version{SemVer:"v2.14.1", GitCommit:"618447cbf203d147601b4b9bd7f8c37a5d39fbb4", GitTreeState:"clean"}
```

You can also verify that the Tiller server is running via `kubectl`.
```bash
kubectl get pods -n kube-system
NAME                                READY   STATUS    RESTARTS   AGE
...
tiller-deploy-8c5679674-k9c7m       1/1     Running   0          47m
...
```

### Create a NuoDB namespace to install NuoDB

`kubectl create namespace nuodb`

### For Red Hat OpenShift only - Grant OpenShift privileges

Create a new service account `nuodb`.

`kubectl -n nuodb create serviceaccount nuodb`

Next, you will want to give your new service account the correct SecurityContextConstraints to run NuoDB.
You can find the recommended SecurityContextConstraints in ([deploy/nuodb-scc.yaml](deploy/nuodb-scc.yaml)

```bash
oc apply -f deploy/nuodb-scc.yaml -n nuodb
oc adm policy add-scc-to-user nuodb-scc system:serviceaccount:nuodb:nuodb -n nuodb
oc adm policy add-scc-to-user nuodb-scc system:serviceaccount:nuodb:default -n nuodb
```

NuoDB charts and their privilege requirements.

| Charts | Privilege | Short Explanation |
| ----- | ----------- | ------ |
| transparent-hugepage| allowHostDirVolumePlugin: true | To mount hostPath and disable THP on host|
| transparent-hugepage| volumes.hostPath | To mount hostPath and disable THP on host|
| transparent-hugepage| seLinuxContext.* | To mount hostPath and disable THP on host|
| admin, database| allowedCapabilities.FOWNER | To change directory ownership in PV to the nuodb process|
| admin, database| defaultAddCapabilities.FOWNER | To change directory ownership in PV to the nuodb process|


## Deploying NuoDB using Helm Charts

**Note:** You MUST first disable Linux Transparent Huge Pages(THP) on all cluster nodes that will host NuoDB pods. Run the `transparent-hugepage` chart first.

- **transparent-hugepage** ([documentation](transparent-hugepage/README.md))

Optionally, consider configuring storage classes for persistent storage use. You can still use persistent storage without using the Storage Class Helm Chart and the NuoDB Admin and Database charts will choose a default based on what's available in your Kubernetes  environment.

- **Storage Classes** ([documentation](storage-class/README.md)) 

Deploy the NuoDB Components in this order : 

- **NuoDB Admin** ([documentation](admin/README.md)) 
- **NuoDB Database** ([documentation](database/README.md)) 


## Cleanup

See the instructions for the individual charts for deleting the applications.
An alternative cleanup strategy is to delete the entire project:

```bash
kubectl delete namespace nuodb
```

[1]: https://helm.sh/docs/using_helm/
[2]: https://github.com/helm/helm/releases
[4]: #getting-started-with-helm-on-openshift
[5]: #deploying-nuodb-using-helm-charts
[6]: https://github.com/nuodb/nuodb-helm-charts

