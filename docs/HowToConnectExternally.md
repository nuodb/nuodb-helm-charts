# Connect to NuoDB Database Externally

## Introduction

NuoDB supports SQL clients via a NuoDB-provided client driver such as JDBC.
NuoDB provides a configurable internal load balancer that can direct clients to Transaction Engines (TEs) according to user-specified load balancing rules.
This enables the administrator to direct different client workloads to different pools of TEs in order to maximize cache locality and execute certain workloads on a TE with the best configuration for that workload.

When a SQL client connects to a database, the driver first establishes a connection with the NuoDB Admin, which will perform load balancing and redirect the client to an appropriate TE in that database.
Load balancing can be configured using the powerful `LBQuery` language to define the database connection strategy for routing SQL clients to the most appropriate TEs.

For more information on NuoDB client connections, see [Client Development](https://doc.nuodb.com/nuodb/latest/client-development/).

For more information on the LBQuery language syntax, see [Load Balancer Policies](https://doc.nuodb.com/nuodb/latest/client-development/load-balancer-policies/).

A _direct_ connection to TEs is also supported by NuoDB, however, in this case, the load balancing must be performed by the client application.

SQL clients and applications running in the same Kubernetes cluster as the NuoDB domain can connect to the database with the default NuoDB Helm Charts configuration using the NuoDB Admin `ClusterIP` service or directly using the TE `ClusterIP` service. This document focuses on external client applications running outside of the Kubernetes cluster where the NuoDB database is hosted. Allowing external access to the NuoDB database is not enabled by default and requires additional configuration.

## Transaction Engine Groups

A TE group consists of TEs with the same configuration which is part of the same database most often used to serve specific SQL workload.
NuoDB Helm Charts 3.4.0+ supports the deployment of one or more TE groups per database.

Multiple Helm releases of the `database` chart for the same NuoDB database can be installed in the same Kubernetes namespace where only one of them is _primary_.
One or more _secondary_ Helm releases are used to deploy additional TE groups for the same database with different configuration options.
This allows flexible configuration of each TE group including but not limited to the number of TEs in a group, their resource requirements, process labels, and scheduling rules.
SQL clients can be configured to target each TE group separately using NuoDB Admin load balancer rules.
Specifying the type of the Helm `database` release is controlled by the  `database.primaryRelease` option (`true` by default).

## External Access for TE Groups

External access for the NuoDB database is _not_ enabled by default.
To enable external access to the NuoDB domain and database, set the `admin.externalAccess.enabled=true` and `database.te.externalAccess.enabled=true` options.
Configuring external database access via the NuoDB Admin with TE groups is supported with NuoDB Helm Charts 3.4.0+ and NuoDB 4.2.4+.
In earlier versions, it is only possible to connect to the provisioned TEs directly and bypass the NuoDB Admin.

A Kubernetes service of type `LoadBalancer` or `NodePort` is created per TE group.
This should allow the external SQL clients to connect to the TEs backing the service by uniquely targeting each TE group.
Most of the cloud vendors provide Kubernetes Load Balancer controllers that support different service annotations used to control the properties and configuration of the provisioned cloud load balancer.
The Kubernetes cluster should be properly configured so that the external network (Layer4) cloud load balancer is provisioned automatically.
For more information on additional cluster configuration, see [Cloud Provider Specifics](#cloud-provider-specifics).

> **NOTE**: When external access is enabled, the NuoDB Helm Charts will create Internet-facing load balancers by default.  Be sure to understand the difference between _internet facing_ load balancers (allowing connectivity external to the cloud and Kubernetes cluster) and  _internal_ load balancers (allowing connectivity external to the Kubernetes cluster but within the cloud virtual network).

If the SQL clients are located outside of the Kubernetes cluster where NuoDB is deployed but in the same cloud provider or virtual network, then _internal_ load balancer can be used.
This is configured by setting the `admin.externalAccess.internalIP=true` and `database.te.externalAccess.internalIP=true` or further customized by explicitly setting custom annotations for the Kubernetes services using `admin.externalAccess.annotations` and `database.te.externalAccess.annotations` options.
The user-provided custom annotations will overwrite the default annotations for the services.

> **IMPORTANT**: The customer is responsible for correctly configuring security rules and restricting access to the cloud load balancers.

Services of type `NodePort` can be created as well by setting the `admin.externalAccess.type=NodePort` and `database.te.externalAccess.type=NodePort` options.
For this type of deployment, a Layer4 load balancer must be manually provisioned. It should also be configured to load balance traffic across all worker nodes in the Kubernetes cluster.
The different TE groups will be reachable on different node ports.

The external address and port for TEs are configured using the `external-address` and `external-port` process labels.
If supplied, they will be advertised by NuoDB Admin to the SQL clients during the second stage of the client connection protocol.
For more information, see [Use External Address](https://doc.nuodb.com/nuodb/latest/client-development/load-balancer-policies/#_use_external_address).

Obtaining and configuring the address of the L4 load balancers can be tedious and error-prone as they are provisioned asynchronously by the Kubernetes controllers.
NuoDB can inspect the Kubernetes services and configure the TE database processes with the `external-address` and `external-port` process labels automatically if the `--enable-external-access` process option is provided.
This simplifies the deployment and ensures the correct configuration of the TE database processes.

- For services of type `LoadBalancer`, NuoDB will configure the `external-address` process label with the service ingress IP or hostname as a value.
- For services of type `NodePort` are provisioned, the customer must configure the  `external-address` and NuoDB will configure the `external-port` process label with the service node port as a value.
If the `--enable-external-access` process option is supplied but the `external-address` process label is no defined, then external access won’t be enabled and a warning message will be logged.

> **NOTE**: If either the hostname or the IP address value(s) of the provisioned cloud load balancer change, the TE database process must be restarted for the new value(s) to take effect.

## Examples

To demonstrate the external access using TE groups, let's consider a working example with the following requirements:

- A NuoDB database is deployed in a single Kubernetes cluster.
- Transactional SQL workload should be processed by TEs in _group 1_.
- Reporting SQL workload should be processed by TEs in _group 2_.
- A set of the applications are installed in the same Kubernetes cluster as NuoDB.
- Another set of the applications are installed in a different Kubernetes cluster, on bare metal, or in a different cloud.

The resource requirements for the different TE groups may be different as there is a direct dependency on the type of SQL workload that TEs will serve.
There will be 2 _smaller_ TEs dedicated for the _transactional_ and 1 _bigger_ TE dedicated for the _reporting_ workload deployed in _nuodb_ namespace.

To satisfy the requirement of having several workloads targeting a different set of TEs, the `LBQuery` connection property will be used.
Alternatively, load balancer policies can be configured by an administrator and SQL clients can reference them using the `LBPolicy` connection property.
For more information, see [Registering Named Load Balancer Policies](https://doc.nuodb.com/nuodb/latest/client-development/load-balancer-policies/#_registering_named_load_balancer_policies_lbpolicy).

For simplicity, the `tx-type` database process label will be used to filter which workload is served by a set of TEs using the `LBQuery` connection property.
The label value is either `transactional` or `reporting`.

NuoDB supports multi-tenant, multi-cluster, and multi-cloud database deployments using TE groups, however, for simplicity a single-cluster single-tenant deployment will be demonstrated here.
The below diagram illustrates the deployed resources and SQL clients along with the `LBQuery` syntax.

![External Access with TE Groups](../images/database-groups.png)

### Deployment

The steps below will deploy a NuoDB database with 2 TE groups in Google Kubernetes Engine (GKE).
If you are deploying in a different environment, make sure that the correct cloud provider is set in the `cloud.provider` option.

> **NOTE**: Use the `nuodb.image.tag` option to specify the NuoDB product version.
NuoDB 4.2.4+ docker image should be used.

Install the [admin](../stable/admin/README.md) chart and enable external access.
Service of type `LoadBalancer` will be provisioned by default.

```shell
helm install admin nuodb/admin \
    --namespace nuodb \
    --cloud.provider=google \
    --set admin.externalAccess.enabled=true
```

Install the [database](../stable/database/README.md) chart for the primary Helm release which deploys TE _group 1_.
Configure the `tx-type=transactional` label for the TEs in this group.

```shell
helm install database-group1 nuodb/database \
    --namespace nuodb \
    --set cloud.provider=google \
    --set database.name=demo \
    --set database.te.externalAccess.enabled=true \
    --set database.te.otherOptions.enable-external-access=true \
    --set database.te.replicas=2 \
    --set database.te.resources.limits.cpu=4 \
    --set database.te.resources.limits.memory=8Gi \
    --set database.te.resources.requests.cpu=4 \
    --set database.te.resources.requests.memory=8Gi \
    --set database.te.labels.tx-type=transactional
```

Install the [database](../stable/database/README.md) chart for the secondary Helm release which deploys TE _group 2_.
Configure the `tx-type=reporting` label for the TEs in this group.

```shell
helm install database-group2 nuodb/database \
    --namespace nuodb \
    --set cloud.provider=google \
    --set database.name=demo \
    --set database.primaryRelease=false \
    --set database.te.externalAccess.enabled=true \
    --set database.te.otherOptions.enable-external-access=true \
    --set database.te.replicas=1 \
    --set database.te.resources.limits.cpu=4 \
    --set database.te.resources.limits.memory=16Gi \
    --set database.te.resources.requests.cpu=4 \
    --set database.te.resources.requests.memory=16Gi \
    --set database.te.labels.tx-type=reporting
```

Wait for the NuoDB database to become ready:

```shell
kubectl exec -ti admin-nuodb-cluster0-0 -n nuodb -- nuocmd check database \
    --db-name demo \
    --check-running \
    --num-processes 4 \
    --wait-forever
```

### Verification

Obtain the external address for the NuoDB Admin service:

```shell
kubectl config set-context --current --namespace=nuodb
DOMAIN_ADDRESS=$(kubectl get services nuodb-balancer -o jsonpath='{.status.loadBalancer.ingress[].ip}')
```

Check database processes node IDs:

```shell
kubectl exec -ti admin-nuodb-cluster0-0 -- nuocmd show database \
    --db-name demo \
    --process-format '{type} {hostname} {node_id}'
```

Use `nuosql`, which can be found in the [NuoDB Client-only Package](https://github.com/nuodb/nuodb-client), to connect to the NuoDB database from the local machine.
Repeat the command several times to ensure that each time the expected node ID is printed.

```shell
echo 'select GETNODEID() from dual;' |  nuosql demo@${DOMAIN_ADDRESS} \
    --user dba --password secret \
    --connection-property 'LBQuery=round_robin(first(label(tx-type transactional) any))'
```

Repeat the steps for the _reporting_ workload.

```shell
echo 'select GETNODEID() from dual;' |  nuosql demo@${DOMAIN_ADDRESS} \
    --user dba --password secret \
    --connection-property 'LBQuery=round_robin(first(label(tx-type reporting) any))'
```

Verify that internal SQL clients can still connect to the database by setting `PreferInternalAddress=true` connection property.

```shell
kubectl exec -ti admin-nuodb-cluster0-0 -- bash -c \
    "echo 'select GETNODEID() from dual;' |  nuosql demo@nuodb-clusterip \
        --user dba --password secret \
        --connection-property 'LBQuery=round_robin(first(label(tx-type transactional) any))' \
        --connection-property 'PreferInternalAddress=true'"
```

## Cloud Provider Specifics

### Native CNI

All of the managed Kubernetes offerings from cloud providers supported by the NuoDB Helm Charts (EKS, AKS, GKE) have the ability to deploy Kubernetes with a native Container Network Interface (CNI) plugin.
This causes Kubernetes pods to be assigned IP addresses from the underlying virtual network CIDR, making them addressable by SQL clients running in the same VPC (or peered VPCs / VPN connections with the correct routing configured).
As such, the TE IP address returned by the second stage of the client connection protocol is already addressable by the client, and it is enough to make only the Admin layer available external to the Kubernetes cluster using the `admin.externalAccess.*` values, to receive client connections.  

> **NOTE**: The node security configuration may also need to explicitly allow the NuoDB TE port (48006).

### GCP

No additional configuration is needed to enable NuoDB database external access in GKE.
For more information, check [Configuring TCP/UDP load balancing](https://cloud.google.com/kubernetes-engine/docs/how-to/service-parameters).

### AWS

By default, the NuoDB Helm Charts 3.4.0 excepts that the [AWS Load Balancer Controller](https://docs.aws.amazon.com/eks/latest/userguide/aws-load-balancer-controller.html) is deployed in the Amazon Elastic Kubernetes Service (EKS) and the native CNI plugin is used when provisioning internet-facing load balancer.
It will automatically provision AWS Network Load Balancer (NLB) for Kubernetes services of type `LoadBalancer`, as described in the [Network load balancing on Amazon EKS](https://docs.aws.amazon.com/eks/latest/userguide/network-load-balancing.html) guide.
For more information on how to customize the provisioned NLB, check [Network Load Balancer](https://kubernetes-sigs.github.io/aws-load-balancer-controller/latest/guide/service/nlb/).
When the `internalIP` option is set to “true”, no additional configuration is needed.

If the [legacy AWS cloud provider load balancer controller](https://kubernetes.io/docs/concepts/services-networking/service/#loadbalancer) is used the Kubernetes service annotations must be changed. For example:

```yaml
cloud:
  provider: aws

admin:
  externalAccess:
    enabled: true
    internalIP: false
    annotations:
      service.beta.kubernetes.io/aws-load-balancer-scheme: internet-facing

database:
  te:
    externalAccess:
      enabled: true
      internalIP: false
      annotations:
        service.beta.kubernetes.io/aws-load-balancer-scheme: internet-facing
```

### Azure

No additional configuration is needed to enable NuoDB database external access in Azure Kubernetes Service (AKS).
For more information, see [Use a public Standard Load Balancer](https://docs.microsoft.com/en-us/azure/aks/load-balancer-standard).

## Troubleshooting

### Unable to connect

There may be different reasons for client connectivity problems such as:

- not _Ready_ NuoDB Admin Pods
- not _Ready_ TE pods
- incorrect external access configuration
- incorrect NuoDB Load Balancer configuration
- incorrect NLB configuration
- network connectivity problems including lack of routing, firewall configuration, and many more

> **ACTION**: You can rule out any of the points above one by one.
Start by checking the overall Pod status for the NuoDB domain and database.
Some of the common troubleshooting steps are listed below, however, there might be additional verifications specific to your deployment.

1. Verify that all AP, TE, and SM pods are reported _Ready_ using `kubectl get pods` command.
2. Check the NuoDB domain and database using `nuocmd show domain` command.
3. Verify the database availability inside the cluster is using the same connection properties as the application uses.
To verify this, use the  `nuosql` tool inside an AP Pod.
4. Make sure that `external-address` and/or `external-port` process labels are configured correctly on the TE database processes using `nuocmd --show-json-fields hostname,labels get processes` command.
If you are using the `--enable-external-access` process option, verify that the configured values are the same as the `EXTERNAL-IP` shown for the Kubernetes service in `kubectl get services` output.
If the value is not the same, restart the TE database process and verify again.
5. Verify that the `LBQuery` or `LBPolicy` syntax is correct. If using `LBPolicy`, verify that the expected policies are configured in the NuoDB Admin using `nuocmd get load-balancers` and `nuocmd get load-balancer-config`.
6. Check Kubernetes events for issues deploying cloud load balancers and external addresses.
For more information on this, see the [TE does not start](#te-does-not-start) section.
7. Verify that the cloud load balancer is provisioned and forwards traffic to the correct Kubernetes cluster.
Check that its configuration is correct and modify the Kubernetes service annotations if needed.
8. Verify that the configured security rules allow external access.

### TE does not start

TEs started with the `--enable-external-access` process option will wait for the `LoadBalancer` service IP address or hostname to be available before they start.
In a case of a problem during NLB provisioning, the IP address will never be populated and the _engine_ container will fail.
The following errors can be seen in the TE container logs:

```text
2021-12-15T07:43:05.015+0000 INFO  [admin-nuodb-cluster0-0:te-database-nuodb-cluster0-demo-55d664c8d7-bn57p] CustomAdminCommands Found service name=database-nuodb-cluster0-demo-balancer, type=LoadBalancer, selector={u'app': u'database-nuodb-cluster0-demo', u'component': u'te'}
2021-12-15T07:43:05.015+0000 INFO  [admin-nuodb-cluster0-0:te-database-nuodb-cluster0-demo-55d664c8d7-bn57p] CustomAdminCommands Waiting for load balancer service database-nuodb-cluster0-demo-balancer ingress address...
'start te' failed: Timeout after 120.0 sec waiting for ingress hostname in service database-nuodb-cluster0-demo-balancer
```

> **ACTION**: Verify that the `EXTERNAL-IP` for the service is available using `kubectl get services` command.
Check the events for the service for this TE group using `kubectl describe service database-nuodb-cluster0-demo-balancer` command.
Look into the cloud provider documentation on how to troubleshoot the load balancer controller.
