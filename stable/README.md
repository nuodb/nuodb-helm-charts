# Getting Started with NuoDB/Helm on OpenShift

This post will walk you through getting both the Tiller server and Helm client up and running on OpenShift, and then installing your first NuoDB Helm Chart. It assumes that you already have the OpenShift `oc client` installed locally and that you are logged into your OpenShift instance.

The instructions are in two parts:

1. **[Getting Started with Helm on OpenShift][4]** covers how to install and configure Helm on a bastion host. It will walk you through deploying a canary application to make sure Helm is properly configured.
2. **[Deploying NuoDB using Helm Charts][5]** covers how to configure hosts to permit running NuoDB, and covers deploying your first NuoDB database using the provided Helm charts.

Bear in mind there are sub-charts in subdirectories included in this distribution. Instructions provided in this specific README are more geared towards the prerequisite setup of projects, Helm and Tiller, security settings, etc. Sub-charts have details for each of the deployed components.

## Getting Started with Helm on OpenShift

### Install Helm and Tiller

There are two parts to Helm: The Helm client (`helm`) and the Helm server (Tiller). This guide shows how to [install the client][1], and then proceeds to show two ways to install the server.

Every [release][2] of Helm provides binary releases for a variety of OSes. These binary versions can be manually downloaded and installed.

1. Download your [desired version][2]
2. Unpack it (`tar -zxvf helm-${helm-version}-linux-amd64.tgz`)
3. Find the helm binary in the unpacked directory, and move it to its desired destination (`mv linux-amd64/helm /usr/local/bin/helm`)

From there, you should be able to run the client: `helm help`.

We’ll use Helm version 2.14.1, which can be downloaded via <https://github.com/kubernetes/helm/releases/tag/v2.14.1.>

Run the following commands to install the Helm client:

```bash
$ curl -s https://storage.googleapis.com/kubernetes-helm/helm-v2.14.1-linux-amd64.tar.gz | tar xz
$ cd linux-amd64
$ mv helm /usr/local/bin
$ mv tiller /usr/local/bin
```

If you're running on Mac, the curl commands above would be:

```bash
$ curl -s https://storage.googleapis.com/kubernetes-helm/helm-v2.14.1-darwin-amd64.tar.gz | tar xz
$ cd darwin-amd64
$ mv helm /usr/local/bin
$ mv tiller /usr/local/bin
```

The following command will configure the current users environment for Helm; it will create a `.helm` directory in `${HOME}`.

```bash
$ helm init --client-only
$HELM_HOME has been configured at /home/clusteradmin/.helm.
Not installing Tiller due to 'client-only' flag having been set
Happy Helming!
```

Alternatively, to keep isolated between individuals sharing a login, you can have separate Kube, Helm, and Tiller state (so run these before the above init command):

```bash
export TILLER_NAMESPACE=nuodb
export HELM_HOME=`pwd`/.helm
export KUBECONFIG=`pwd`/.kube/config
```

### Create a Project

We're limiting the Tiller server to manage infrastructure in the project (namespace) where it is installed. So the TILLER_NAMESPACE equals the namespace where NuoDB will be installed.

When we install the Helm client, it will need to know the name of the namespace (project) where Tiller is installed. This can be indicated by locally setting the `TILLER_NAMESPACE` environment variable as follows:

```bash
export TILLER_NAMESPACE=nuodb
```

We’ll be imaginative and call the project `nuodb`, but you can call it anything you like. However, following best practice, every project should be completely isolated and have **its own tiller installation**.

```bash
oc new-project ${TILLER_NAMESPACE}
```

If you already have an OpenShift project you want to use, select it as follows:

```bash
oc project ${TILLER_NAMESPACE}
```

### Install the Tiller Server

In principle this can be done using `helm init`, but currently the helm client doesn’t fully set up the service account rolebindings that OpenShift expects. To try to keep things simple, we’ll use a pre-prepared OpenShift template instead. The template sets up a dedicated service account for the Tiller server, gives it the necessary permissions, then deploys a Tiller pod that runs under the newly created SA.

```bash
$ oc process -f https://github.com/openshift/origin/raw/master/examples/helm/tiller-template.yaml -p TILLER_NAMESPACE="${TILLER_NAMESPACE}" -p HELM_VERSION=v2.14.1 | oc create -f -
```

At this point, you’ll need to wait for a moment until the Tiller server is up and running:

```bash
$ oc rollout status deployment tiller
deployment "tiller" successfully rolled out
```

We’ll check that the Helm client and Tiller server are able to communicate correctly by running helm version. The results should be as follows:

```bash
$ helm version
Client: &version.Version{SemVer:"v2.14.1", GitCommit:"618447cbf203d147601b4b9bd7f8c37a5d39fbb4", GitTreeState:"clean"}
Server: &version.Version{SemVer:"v2.14.1", GitCommit:"618447cbf203d147601b4b9bd7f8c37a5d39fbb4", GitTreeState:"clean"}
```

### Grant the Tiller Server Edit Access to the Current Project

The Tiller server will probably need at least `edit` access to each project where it will manage applications. In the case that Tiller will be handling Charts containing Role objects, `admin` access will be needed.

```bash
oc policy add-role-to-user edit "system:serviceaccount:${TILLER_NAMESPACE}:tiller"
```

The [blog article written by Red Hat][0] suggests having `Tiller` in its own project, separate from application projects. However, this is **NOT** how banks deploy Tiller, it is **NOT** best practice from a security perspective. Banks typically deploy Tiller in each application's project. That way:

- the tiller `edit` role does not have to have `admin` access to the whole cluster
- each project can be individually managed w.r.t. the version of `tiller` managing its deployments; e.g., no **"big bang"** during an upgrade.

So instead lets only run our applications alongside `tiller`, in the same project. To deploy a canary application, such as Node.js:

```bash
$ helm install https://github.com/jim-minter/nodejs-ex/raw/helm/helm/nodejs-0.1.tgz -n nodejs-ex
```

Verify that `tiller` properly launched Node.js.

```bash
$ oc get pods
NAME                      READY     STATUS    RESTARTS   AGE
nodejs-example-1-build    1/1       Running   0          1m
tiller-86c4495fcc-rhq97   1/1       Running   0          15m
```

Now that we've seen Helm deploy an application successfully, we're confident in its installation and configuration, lets now move on to deploying NuoDB. But first, lets clean up Node.js:

```bash
$ helm del --purge nodejs-ex
```

## Deploying NuoDB using Helm Charts

### Deploy NuoDB (et al) via Helm Charts

    IMPORTANT:
    
    You MUST first disable THP on nodes where NuoDB will run. Run the `transparent-hugepage` chart first.

In a nutshell the order of installation is:

- **transparent-hugepage** ([documentation](transparent-hugepage/README.md))
- **admin** ([documentation](admin/README.md))
- **monitoring-influx** ([documentation](monitoring-influx/README.md))
- **database** ([documentation](database/README.md))
- **backup** ([documentation](backup/README.md))
- **restore** ([documentation](restore/README.md))
- **demo-ycsb** ([documentation](demo-ycsb/README.md))
- **demo-quickstart** ([documentation](demo-quickstart/README.md))

See the instructions for the individual charts for deploying the applications.

## Cleanup

See the instructions for the individual charts for deleting the applications.

An alternative cleanup strategy is to delete the entire project:

```bash
oc delete project <project-name>
```

## References

1. Materials herein unscrupulously stolen from [an online article written by Jim Minter, Red Hat, September 21, 2017][0].

[0]: https://blog.openshift.com/getting-started-helm-openshift/
[1]: https://helm.sh/docs/using_helm/
[2]: https://github.com/helm/helm/releases
[3]: https://docs.google.com/document/d/1G1Ljwe0c97KsH881QPUZK6ZtIShCk8jkxskXehuLpKw/edit#
[4]: #getting-started-with-helm-on-openshift
[5]: #deploying-nuodb-using-helm-charts
