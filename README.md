# Helm Charts

[![Build Status](https://travis-ci.org/nuodb/nuodb-helm-charts.svg?branch=master)](https://travis-ci.org/nuodb/nuodb-helm-charts)

Use this repository to submit official Charts for NuoDB. Charts are curated application definitions for Helm. For more information about installing and using Helm, see its
[README.md](https://github.com/helm/helm/tree/master/README.md). To get a quick introduction to Charts see this [chart document](https://github.com/helm/helm/blob/master/docs/charts.md).

For more information on using Helm, refer to the [Helm's documentation](https://github.com/kubernetes/helm#docs).

## Software Release requirements

| Software   | Release Requirements                           | 
|------------|------------------------------------------------|
| Kubernetes |  The latest and previous minor versions of Kubernetes. For example, if the latest minor release of Kubernetes is 1.15 then 1.15 and 1.14 are supported. Charts may still work on previous versions of Kubernertes even though they are outside the target support window.|
| Helm       |  Version 2.x, 2.9 or greater   |
| NuoDB      |  Version [4.0](https://hub.docker.com/r/nuodb/nuodb-ce/tags) and onwards. |
| NuoDB Helm Charts      |  For a list of supported NuoDB Helm Chart releases and where to download, click the `Releases` tab above. To enable automated notification of new releases, click the `Watch` button above and subscribe to the `Releases Only` selection. |

## NuoDB Helm Chart Installation

The default repository for NuoDB is located at https://nuodb-charts.storage.googleapis.com/ and must be enabled.

To add the charts for your local client, run the `helm repo add` command below:

```
helm repo add nuodb https://nuodb-charts.storage.googleapis.com/
"nuodb" has been added to your repositories
```

To list the NuoDB charts added to your repository, run `helm search nuodb/`

To install a chart into your Kubernetes cluster, run 

```
helm init
helm install nuodb/<chart>
```

## NuoDB Helm Chart Incubator Repository

The Incubator repository contains enhancements not yet available in the supported releases. To add the Incubator charts for your local client, run the `helm repo add` command below:

```
helm repo add nuodb-incubator https://nuodb-charts-incubator.storage.googleapis.com/
"nuodb-incubator" has been added to your repositories
```

To list the NuoDB incubator charts added to your repository, run `helm search nuodb-incubator/`

To install an incubator chart into your Kubernetes cluster, run 

```
helm init
helm install nuodb-incubator/<chart>
```

## Repository Structure

This GitHub repository contains the source for the packaged and versioned charts released in the [`gs://nuodb-charts` Google Storage bucket](https://console.cloud.google.com/storage/browser/nuodb-charts/) (the Chart Repository).

The Charts in the `stable/` directory in the master branch of this repository match the latest packaged Chart in the Chart Repository, though there may be previous versions of a Chart available in that Chart Repository.

The purpose of this repository is to provide a place for maintaining and contributing official Charts, with CI processes in place for managing the releasing of Charts into the Chart Repository.

The Charts in this repository are organized into two folders:

* stable
* incubator

Stable Charts meet the criteria in the [technical requirements](CONTRIBUTING.md#technical-requirements).

Incubator Charts are those that do not meet these criteria. Having the incubator folder allows charts to be shared and improved on until they are ready to be moved into the stable folder. The charts in the `incubator/` directory can be found in the [`gs://nuodb-charts-incubator` Google Storage Bucket](https://console.cloud.google.com/storage/browser/nuodb-charts-incubator).

In order to get a Chart from incubator to stable, Chart maintainers should open a pull request that moves the chart folder.


To provide that support the API versions of objects should be those that work for both the latest minor release and the previous one.


These chart repository supports 
## Status of the Project

This project is still under active development, so you might run into [issues](https://github.com/nuodb/nuodb-helm-charts/issues). If you do, please don't be shy about letting us know, or better yet, contribute a fix or feature.
