### NuoDB Insights and Sumo Logic

---

Sumo Logic is an observability suite offering end-to-end observability and security tools within one platform. Sumo Logic is an alternative monitoring solution to the current NuoDB Insights suite.

#### Setup of database values yaml file

---

Using the NuoDB [database](https://github.com/nuodb/nuodb-helm-charts/blob/master/stable/database/README.md) helm chart, set in values file the following values for ```nuocollector.*``` section.

Refer to [Sumo Logic Collectors](#sumo-logic-collectors) to harvest the Sumo Logic endpoint to insert into the yaml snippet below.

Set nuocollector to be enabled.

```
nuocollector:
   enabled: true
```

Set the sumologic output configuration in the plugins section.

```
  plugins:
    ## NuoDB Collector compatible plugins specific for database services
    database:
      sumologic.conf: |
        [[outputs.sumologic]]
          url = "sumologic endpoint url here"
          namepass = [ "SqlListener*", "Commits", "Inserts", "Updates", "Deletes", "WriteThrottleTime", "Pending*", "Archive*Time", "Journal*Time" ]
          fieldpass = [ "rate", "value" ]
          data_format = "carbon2"
        [[outputs.sumologic]]
          url = "sumologic endpoint url here"
          namepass = [ "Summary.*" ]
          fieldpass = [ "raw" ]
          data_format = "carbon2"
```

Example:

```
  plugins:
    ## NuoDB Collector compatible plugins specific for database services
    database:
      sumologic.conf: |
        [[outputs.sumologic]]
          url = "https://endpoint1.collection.eu.sumologic.com/receiver/v1/http/ZaVnC4dhaV0laeqZ1_FuM8b_P9gPRX4DIJzrJ8eaIy0Lz0rQH5o-bIARKiKf7M09suw6VFJpKSIt_3zv0hq80kcq0VivLeZWM3yXey7EXDG5YFFBaMWqMg=="
          namepass = [ "SqlListener*", "Commits", "Inserts", "Updates", "Deletes", "WriteThrottleTime", "Pending*", "Archive*Time", "Journal*Time" ]
          fieldpass = [ "rate", "value" ]
          data_format = "carbon2"
        [[outputs.sumologic]]
          url = "https://endpoint1.collection.eu.sumologic.com/receiver/v1/http/ZaVnC4dhaV0laeqZ1_FuM8b_P9gPRX4DIJzrJ8eaIy0Lz0rQH5o-bIARKiKf7M09suw6VFJpKSIt_3zv0hq80kcq0VivLeZWM3yXey7EXDG5YFFBaMWqMg=="
          namepass = [ "Summary.*" ]
          fieldpass = [ "raw" ]
          data_format = "carbon2" 
```

### Sumo Logic Dashboards

---

1. [Sumo Logic Collectors](#sumo-logic-collectors) for NuoDB Insight data.

2. [NuoDB Insights graphs](#insights-graphs-on-sumo-logic-dashboards) in Sumo Logic dashboards.

3. [Kubernetes Graphs](#sumo-logic-kubernetes-graphs) in Sumo Logic (Optional)

#### Sumo Logic Collectors

---

Website is https://www.sumologic.com

Create a new account or use an existing account. For a new account, skip past the introductions and go straight to the dashboard. At the time of writing, the skip link was below the introduction text.

Click on <b>Manage Data</b>: 

![Manage Data](./images/sumologic0.png)

Click on **Add Collector**

![Add Collector](./images/sumologic1.png)

Select **Hosted Collector**

![Hosted Collector](./images/sumologic2.png)

Add an appropriate name for the collector.

Next is to add a source to the newly created collector. At the time of writing, the source was under "Generic" and is  "HTTP logs and metrics"

![Generic_logs and metrics](./images/sumologic4.png)

Add an appropriate name for the source for the collector. 

Next is to copy and keep safe the Sumo Logic metrics endpoint. This will be needed in the section for the helm database values file setup. It will be something like:

```
https://endpoint1.collection.eu.sumologic.com/receiver/v1/http/ZaVnC4dhaV0laeqZ1_FuM8b_P9gPRX4DIJzrJ8eaIy0Lz0rQH5o-bIARKiKf7M09suw6VFJpKSIt_3zv0hq80kcq0VivLeZWM3yXey7EXDG5YFFBaMWqMg==```
```

![sumologic endpoint](./images/sumologic6.png)

#### 

#### Insights graphs on Sumo Logic dashboards

---

Click on Personal -> three dots to the right of the side bar and create a new folder.

![Personal_Folder](./images/sumologic21.png)

Click on the 3 dots to the right of the folder name and select Import

![Import](./images/sumologic22.png)

Upload the NuoDB dashboard JSON as a copy paste operation.  The JSON files that can be copy pasted can be found in the following locations

```
files/NuoDB_AdHoc.json
files/NuoDB_DBA.json
files/NuoDB_NOC.json
```

Click "Import" at the bottom of the pane.

 ![Nuodb_Dashboard_JSON](./images/sumologic25.png)

Set the Collector name and Source name to be the same as set in section 1.

![Set_Collector_and_Source](./images/sumologic26.png)

**Note:** The Sumo Logic NuoDB graphs are a subset of the graphs that are available in Insights. Data that is sent to Sumo Logic is a subset of all the data that is sent to Insights. 

#### Sumo Logic Kubernetes Graphs

---

From the App Catalog search for classic apps "kubernetes" -> "Kubernetes" and follow the install instructions. 

**Tip:** Use a different namespace other than the namespace for the NuoDB install. 

**Tip:** Take care to make a note of the ```accessId```, ```accessKey```, ```clusterName``` and ```collectorName```.  ```accessKey``` is shown only once.
