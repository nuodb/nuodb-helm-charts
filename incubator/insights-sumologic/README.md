# NuoDB Insights and Sumo Logic

---

Sumo Logic is an observability suite offering end-to-end observability and security tools within one platform. Sumo Logic is an alternative monitoring solution to the current NuoDB Insights suite.

### demo-insight-sumologic Helm Chart

---

This chart deploys the Insights Sumologic Demo on a Kubernetes cluster using the helm chart package manager.

### Software Version Prerequisites

---

Please visit the **[NuoDB Helm Chart main page](https://github.com/nuodb/nuodb-helm-charts/#software-release-requirements)** for software version prerequisites.

### demo-insight-sumologic Helm Chart prerequisites

---

Please review our [guide](./README-insights-sumologic.md) on setting up the Sumo Logic dashboards for NuoDB Insights and ensure that the ```sumologic.endpoint``` url is harvested from the Sumo Logic dashboard. When installing NuoDB from the helm charts, ensure that the option ```nuocollector.enabled=true```

### Command

---

```
helm install [name] nuodb-incubator/demo-insights-sumologic [--generate-name] [--set parameter] [--values myvalues.yaml]```
```

### Installing the Chart

---

All configurable parameters for each top-level scope are detailed below, organised by scope.

**sumologic.***

The purpose of this section is to specify **sumologic** settings

| Parameter | Description         | Default |
| --------- | ------------------- | ------- |
| endpoint  | Sumo Logic Endpoint | null    |

For example, to enable the Sumo Logic Endpoint collection url, 

```
sumologic:
  endpoint: https://endpoint1.collection.eu.sumologic.com/receiver/v1/http/ZaVnC4dhaV0laeqZ1_FuM8b_P9gPRX4DIJzrJ8eaIy0Lz0rQH5o-bIARKiKf7M09suw6VFJpKSIt_3zv0hq80kcq0VivLeZWM3yXey7EXDG5YFFBaMWqMg==
```

Verify the Helm Chart:

```
helm install insight-sumologic nuodb-incubator/demo-insight-sumologic --debug --dry-run
```

Deploy the demo:

```
helm install insight-sumologic nuodb-incubator/demo-insight-sumologic
```

**Tip:** Wait until the nuocollector sidecar is re-started with the additional Sumo Logic configuration. This should take a minute or less.

Verify traces in the Sumo Logic [graphs](./README-insights-sumologic.md).

Example of the DBA graph set:
![sumologic graph sample](./images/sumologic26.png)

### Uninstalling the Chart

To uninstall/delete the deployment

```
helm delete insight-sumologic
```

The command removes all the Kubernetes components associated with the chart and deletes the release
