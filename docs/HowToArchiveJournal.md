# NuoDB Storage Manager (SM) Journal in Kubernetes

## Introduction

The storage manager, SM for short, is responsible for maintaining a complete copy of the database.
The atoms, database elements, are stored to either local disk, or in a separate volume.
These atoms are written by a specific module called the `archive`.
However, the archive doesn’t have any special mechanisms to satisfy the durability requirements.
In addition to the archive, NuoDB storage managers also maintain a write-ahead-log called the `journal`.
The journal has the onerous task of ensuring durability in the face of an unexpected process termination.
E.g., a power loss, machine meltdown, or an unexpected cloud outage.
The journal will make sure that the archive is reconstructed in a consistent state when coming back online.

For more info on NuoDB database journaling, please consult the official [NuoDB docs](https://doc.nuodb.com/nuodb/latest/database-administration/about-database-journaling/).

## Separating the Journal from the Archive

Since the `journal` has to write all commits durably to disk, the speed of the disk directly influences the commit latency.
To achieve the best performance, NuoDB recommends placing the `journal` on the disk with the best latency available.
In contrast, the `archive` should be using the disk with the largest throughput.

By default, the `journal` is located in a subdirectory of the `archive`.
To achieve the best cost vs speed tradeoff, you can separate the journal from the archive.
This can be done by setting the database helm template values `database.sm.noHotCopy.journalPath.enabled` and `database.sm.hotCopy.journalPath.enabled` to `true` and configure it with the desired persistence settings.

Kubernetes stateful sets volume mounts are immutable and as such, the setting can not be changed easily on an existing database.

## Upgrading existing domains

### Overview

StatefulSet volume mounts are immutable and any attempts to change this value will result in the following error:
```
Error: UPGRADE FAILED: cannot patch "sm-database-nuodb-cluster0-demo-hotcopy" with kind StatefulSet: \
StatefulSet.apps "sm-database-nuodb-cluster0-demo-hotcopy" is invalid: spec: Forbidden: \
updates to statefulset spec for fields other than 'replicas', 'template', and 'updateStrategy' are forbidden
```

To change the domain from using `an in-archive journal` to an `a separated journal` or vice versa, Kubernetes requires the StatefulSet to be deleted and recreated with the new settings.

NuoDB Helm Charts offer two distinct Storage Manager StatefulSets (SS).
One StatefulSet controls Storage Managers which have a `backup` directory and are referred to as `HotCopy SMs`.
The second StatefulSet which does not have such a directory, and we call Storage Managers in this SS `NoHotCopy SMs`.
Since these StatefulSets can be upgraded independently, NuoDB will not suffer downtime or data loss when performing this upgrade operation.

### Pre-requisites

1) Your domain must contain an Enterprise License, and the ability to run two or more Storage Managers.
2) You must have at least 1 Storage Manager in both `HotCopy SM` SS and the `NoHotCopy SM` SS.
3) Your database option `max-lost-archives` is set to allow the loss of a whole StatefulSet.

### Upgrade Process
In this example, we will assume that the domain is running 1 Storage Manager in each StatefulSet and the name of the database is `demo` as follows:
```shell
kubectl get statefulsets.apps -n nuodb
NAME                                      READY
admin-nuodb-cluster0                      1/1
sm-database-nuodb-cluster0-demo           1/1
sm-database-nuodb-cluster0-demo-hotcopy   1/1
```

NOTE: NuoDB also requires the NuoDB admin layer, which is also a StatefulSet but won't be directly involved in this migration process.

There will be a number of PersistentVolumeClaims.
- 1 PVC for the NoHotCopy Storage Manager
- 2 PVCs for the HotCopy Storage Manager
- 1 PVC for the admin layer

```shell
kubectl get persistentvolumeclaims -n nuodb
NAME                                                       STATUS   VOLUME                                     CAPACITY   ACCESS MODES
archive-volume-sm-database-nuodb-cluster0-demo-0           Bound    pvc-e8227428-6b9f-47b1-acc0-41b70a21d043   20Gi       RWO
archive-volume-sm-database-nuodb-cluster0-demo-hotcopy-0   Bound    pvc-19540f97-8407-4fe1-a940-3b3f83538f0f   20Gi       RWO
backup-volume-sm-database-nuodb-cluster0-demo-hotcopy-0    Bound    pvc-daebd030-a4eb-46e1-8451-6d65e4fa3061   20Gi       RWO
raftlog-admin-nuodb-cluster0-0                             Bound    pvc-4930b80b-a622-487a-90a6-60a2bd3f0548   1Gi        RWO

```


#### Upgrade the NoHotCopy StatefulSet
First, scale the StatefulSet to 0.
```shell
kubectl scale statefulset -n nuodb sm-database-nuodb-cluster0-demo --replicas=0
statefulset.apps/sm-database-nuodb-cluster0-demo scaled
```

Second, delete all PVCs
```shell
kubectl delete pvc -n nuodb <PVC_NAME>
```

In this case, the `<PVC_NAME>` will be `archive-volume-sm-database-nuodb-cluster0-demo-0`.

Third, delete the StatefulSet:
```shell
kubectl delete statefulset -n nuodb sm-database-nuodb-cluster0-demo
```

Fourth, reinstall the StatefulSet using Helm.
To enable the journal on the deleted StatefulSet, use the following value:
```
database.sm.noHotCopy.journalPath.enabled=true
```

```shell
helm upgrade -n nuodb database nuodb/database \
--set database.sm.noHotCopy.journalPath.enabled=true \
--set database.sm.noHotCopy.replicas=1
```

Wait for the Storage Manager Pod to become READY before proceeding.

#### Upgrade the HotCopy StatefulSet
First, scale the StatefulSet to 0.
```shell
kubectl scale statefulset -n nuodb sm-database-nuodb-cluster0-demo-hotcopy --replicas=0
statefulset.apps/sm-database-nuodb-cluster0-demo-hotcopy scaled
```

Second, delete all PVCs
```shell
kubectl delete pvc -n nuodb <PVC_NAME>
```

In this case, the `<PVC_NAME>` will be `archive-volume-sm-database-nuodb-cluster0-demo-hotcopy-0`.

Third, delete the StatefulSet:
```shell
kubectl delete statefulset -n nuodb sm-database-nuodb-cluster0-demo-hotcopy
```

Fourth, reinstall the StatefulSet using Helm.
To enable the journal on the deleted StatefulSet, use the following value:
```
database.sm.hotCopy.journalPath.enabled=true
```

```shell
helm upgrade -n nuodb database nuodb/database \
--set database.sm.noHotCopy.journalPath.enabled=true \
--set database.sm.hotCopy.journalPath.enabled=true \
--set database.sm.noHotCopy.replicas=1
```

Wait for the Storage Manager Pod to become READY.

Clean up the remaining `backup-` PersistentVolumeClaims.
```shell
kubectl delete pvc -n nuodb backup-volume-sm-database-nuodb-cluster0-demo-hotcopy-0
```