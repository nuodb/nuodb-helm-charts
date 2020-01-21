# Helm Charts

[![Build Status](https://travis-ci.org/nuodb/nuodb-helm-charts.svg?branch=master)](https://travis-ci.org/nuodb/nuodb-helm-charts)

Use this repository to submit official Charts for NuoDB. Charts are curated application definitions for Helm. For more information about installing and using Helm, see its
[README.md](https://github.com/helm/helm/tree/master/README.md). To get a quick introduction to Charts see this [chart document](https://github.com/helm/helm/blob/master/docs/charts.md).

For more information on using Helm, refer to the [Helm's documentation](https://github.com/kubernetes/helm#docs).

## How to use this repository?

For a list of supported NuoDB Helm Chart releases and where to download, click the `Releases` tab above.

To enable automated notification of new releases, click the `Watch` button above and subscribe to the `Releases Only` selection.

## How do I install these charts?

To install, run `helm install nuodb/<chart>`. This is the default repository for NuoDB which is located at
 https://nuodb-charts.storage.googleapis.com/ and must be enabled to use.

To add the charts for your local client, run the `helm repo add` command below:

```bash
$ helm repo add nuodb https://nuodb-charts.storage.googleapis.com/
"nuodb" has been added to your repositories
```

To list the installed NUoDB charts, run `helm search repo nuodb/`

## How do I enable the Incubator repository?

The Incubator repository contains enhancements not yet available in the supported releases. To add the Incubator charts for your local client, run the `helm repo add` command below:

```bash
$ helm repo add nuodb-incubator https://nuodb-charts-incubator.storage.googleapis.com/
"nuodb-incubator" has been added to your repositories
```

To list the installed NuoDB incubator charts, run `helm search repo nuodb-incubator/`

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

## Supported Kubernetes Versions

This chart repository supports the latest and previous minor versions of Kubernetes. For example, if the latest minor release of Kubernetes is 1.8 then 1.7 and 1.8 are supported. Charts may still work on previous versions of Kubernertes even though they are outside the target supported window.

To provide that support the API versions of objects should be those that work for both the latest minor release and the previous one.

## Supported NuoDB Versions

These chart repository supports NuoDB version [4.0](https://hub.docker.com/layers/nuodb/nuodb-ce/4.0/images/sha256-aaa558ef71795f15d5b3a1ef07b6be4890925dbd023c59b1f9a674ca20614763) and onwards.

## Status of the Project

This project is still under active development, so you might run into [issues](https://github.com/nuodb/nuodb-helm-charts/issues). If you do, please don't be shy about letting us know, or better yet, contribute a fix or feature.
