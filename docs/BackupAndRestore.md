# NuoDB Backup and Restore in Kubernetes

<details>
<summary>Table of Contents</summary>
<!-- TOC -->

- [NuoDB Backup and Restore in Kubernetes](#nuodb-backup-and-restore-in-kubernetes)
    - [Introduction](#introduction)
        - [Terminology](#terminology)
    - [Compatibility Matrix](#compatibility-matrix)
    - [Backup](#backup)
        - [Scheduled online database backups](#scheduled-online-database-backups)
            - [Backup retention management](#backup-retention-management)
    - [Restore](#restore)
        - [Restore/Import source and type](#restoreimport-source-and-type)
        - [Fine-grained archive selection](#fine-grained-archive-selection)
        - [Distributed restore](#distributed-restore)
            - [Distributed database restore from source](#distributed-database-restore-from-source)
            - [Distributed database restore with storage groups](#distributed-database-restore-with-storage-groups)
            - [Manual database restore](#manual-database-restore)
        - [Automatic archive initial import](#automatic-archive-initial-import)
        - [Automatic archive restore](#automatic-archive-restore)
        - [Archive seed restore](#archive-seed-restore)
    - [Troubleshooting](#troubleshooting)
        - [Backup failure](#backup-failure)
            - [No SM matching backup labels](#no-sm-matching-backup-labels)
            - [Backup not completed in the specified timeout](#backup-not-completed-in-the-specified-timeout)
            - [Overlapping backups](#overlapping-backups)
            - [Backup storage full](#backup-storage-full)
        - [Restore failure](#restore-failure)
            - [Invalid archive ID](#invalid-archive-id)
            - [Invalid restore source](#invalid-restore-source)

<!-- /TOC -->
</details>

## Introduction

NuoDB provides several automated mechanisms for backing up and restoring a database in a Kubernetes environment when deploying NuoDB using the NuoDB Helm Charts.

After database installation, the available backup and restore mechanisms are: 

1. [Scheduled online database backups](#scheduled-online-database-backups)
2. [Fine-grained archive selection](#fine-grained-archive-selection)
3. [Distributed database restore from source](#distributed-database-restore-from-source)
4. [Distributed database restore with storage groups](#distributed-database-restore-with-storage-groups)
5. [Manual database restore](#manual-database-restore)
6. [Automatic archive initial import](#automatic-archive-initial-import)
7. [Automatic archive restore](#automatic-archive-restore)
8. [Archive seed restore](#archive-seed-restore)

Several of these mechanisms require additional configuration before using.

> **NOTE**: For information on NuoDB backup and restore operations in non-Kubernetes deployments of NuoDB, see [Backing Up and Restoring Databases](https://doc.nuodb.com/nuodb/latest/deployment-models/physical-or-vmware-environments-with-nuodb-admin/database-operations/backing-up-and-restoring-databases/). 
The documentation in this repository expands on the product documentation and is specific to this Helm Chart repository.

### Terminology

- `SM` = A Kubernetes container that runs NuoDB Storage Manager. 
In Kubernetes deployments, all database processes are started externally by a helper script called `nuodocker` which allows for greater control including the execution of prerequisite actions like database or archive creation in the NuoDB Admin tier that must occur prior to actual database process start up.
- `Archive` =  The persistent storage on disk of a NuoDB database served by a Storage Manager process (SM).
A single NuoDB archive is represented by an archive directory and journal directory located on persistent storage.
In Kubernetes, the journal directory is internal to the archive and normally located inside the archive directory.
Each archive on disk has metadata representation in the admin layer uniquely identified by its `archiveId`.
- `Hot copy` = Online backup for a database archive. This means taking a backup while the database is running.
- `Backup set` = Hot copies are organized into backup sets. Backup sets provide a convenient mechanism for identifying hot copies and maintaining metadata information. A backup set is a directory created when initiating a full hot copy.  
- `Offline backup` = A copy of a database archive taken while the Storage Manager (SM) process serving the archive is stopped. The copy/snapshot should include both the NuoDB archive and the NuoDB journal in a transactionally-consistent manner.
- `HC SM` = Storage Manager engine pod controlled by NuoDB hot copy statefulset.
- `nonHC SM` = Storage Manager engine pod controlled by NuoDB non-hotcopy statefulset.
- `Database in-place restore` = A database restore that allows the data stored in the current database to go back in time without the need to provision additional resources (like more database processes, storage, etc).
In most cases, this type of restore can be used in catastrophic failure situations, where restore is performed from the last good backup.
- `Coordinator SM` = Storage Manager that is selected to prepare the domain state for database in-place restore.
The coordinator SM will delete not restored archives and wait for all other archives to be restored before starting. 
All other database processes wait for the coordinator SM to start before continuing their start up.
- `Backup group` = A subset of Storage Managers in a database to which a backup request is sent by a scheduled backup job.
Typically all HC SMs deployed in a single cluster are part of one backup group.
- `Named backup group`= Backup group that has a custom name configured by the user.
- `Seed restore` = A restore that allows a single archive to be restored using a backup. The archive is then updated with the latest transactions and brought current through the database _SYNCing_ process from a RUNNING SM.
Normally used to speed up the SYNC process for the so-called "accelerated" sync use case.

## Compatibility Matrix

The table below shows the compatibility matrix for different backup and restore mechanisms in combination with the NuoDB product version and NuoDB Helm Charts.

|                                       | NuoDB 4.0.x+ Helm Charts 2.x.x | NuoDB 4.0.x+ Helm Charts 3.0.x | NuoDB 4.0.x+ Helm Charts 3.1.x | NuoDB 4.2.x+ Helm Charts 3.2.x+ |
|---------------------------------------|--------------------------------|--------------------------------|---------------------------------|---------------------------------|
| Scheduled online database backups     | ✓<sup>[1]</sup>                | ✓<sup>[1]</sup>                | ✓                               | ✓                               |
| Fine-grained archive selection        | -                              | -                              | -                               | ✓                               |
| Distributed database restore from source | ✓<sup>[2]</sup><sup>[3]</sup>  | ✓<sup>[2]</sup><sup>[3]</sup>  | ✓<sup>[3]</sup>                 | ✓<sup>[3]</sup>                 |
| Distributed database restore with storage groups  | -                              | -                              | -                               | ✓                               |
| Manual database restore               | -                              | -                              | -                               | ✓                               |
| Automatic archive initial import      | ✓                              | ✓                              | ✓                               | ✓                               |
| Automatic archive restore             | ✓                              | ✓                              | ✓                               | ✓                               |
| Archive seed restore                  | ✓<sup>[4]</sup>                | ✓<sup>[4]</sup>                | ✓<sup>[4]</sup>                 | ✓                               |

[1] - Scheduled online backups workflow is supported by _initial_ database backup job and _post-restore_ cron job.
Additional capability was added into the `nuobackup` script allowing the removal of these jobs in [Helm Charts v3.1.0](https://github.com/nuodb/nuodb-helm-charts/releases/tag/v3.1.0).

[2] - For deployments that contain nonHC SMs, it's recommended that database in-place restore is done by manually starting database engines in sequential order. First start the HC SMs, then nonHC SMs and last the TEs.

[3] - _Automatic archive initial import_ is recommended when restoring a backup set taken from a different environment.

[4] - _Archive seed restore_ is not deterministic. Once the restore request is created, the first SM that starts or is restarted will perform the restore of its archive.

For detailed information about features, enhancements, and fixed problems, please check [NuoDB Release Notes](https://doc.nuodb.com/nuodb/latest/release-notes/) and [NuoDB Helm Charts Release Notes](https://github.com/nuodb/nuodb-helm-charts/releases)

> **NOTE**: The examples and the sample output shown in this document are specific to NuoDB 4.2.x+ and Helm Charts 3.2.x+.

## Backup

### Scheduled online database backups

Multiple backup cron jobs are scheduled as part of the [database](../stable/database) chart deployment by default.
The backup schedules can be customized to meet different backup and recovery objectives.
They are exposed as Helm chart values in particular:
- `database.hotCopy.fullSchedule` - configures the schedule of _full_ backup Kubernetes job
- `database.hotCopy.incrementalSchedule` - configures the schedule of _incremental_ backup Kubernetes job
- `database.hotCopy.journalBackup.intervalMinutes` - configures the schedule of _journal_ backup Kubernetes job. _journal_ hot copy must be explicitly enabled using `database.hotCopy.journalBackup.enabled`.

[Cron expression format](https://en.wikipedia.org/wiki/Cron) documents the format of the first two schedules listed above.

An HC SM has a backup volume attached and is selected for hot copy operation during a scheduled backup.
Each _full_ hot copy will create a backup set with the same name in the backup root directory (by default `/var/opt/nuodb/backup`) on all HC SMs in a single database and a single backup group.
By default, all HC SMs in a single Kubernetes cluster are part of one backup group.
This can be adjusted by changing the `database.hotCopy.backupGroup` value for multiple Kubernetes clusters in a multi-cluster deployment.

_Incremental_ and _journal_ hot copies are always stored in the current backup set if one exists.
If there is no current backup set for some of the archives served by a HC SM when _incremental_ or _journal_ backup is requested, a _full_ hot copy is triggered automatically.
The current backup set is the backup set created by the latest issued hotcopy and can be seen using `nuodocker get current-backup --db-name <db>`.

The `nuobackup` script is made available as a configMap and is used by all backup cron jobs to send a hot copy request of the specified type to the HC SMs.
The _Full_ hot copy result is recorded into the NuoDB Admin domain key-value (KV) store upon successful finish.
Use `nuobackup --type report-latest --db-name <db> --group <backup group>` to show the latest successful _full_ hot copy backup set for a specific backup group.

> **IMPORTANT**: It is important to take periodic journal hot copies when journal hot copy is enabled. Otherwise, the journal will keep growing which can lead to filling the archive volume space.

It is possible to disable backup jobs creation by setting `hotCopy.enableBackups` to `false`.
Custom backup jobs that use the `nuobackup` script can be configured to support more complex workflows and backup strategies.

The logs for each backup job can be seen using `kubectl get logs <job pod name>`. For example:

```bash
kubectl logs full-hotcopy-demo-cronjob-1621415040-vxlq8 --namespace nuodb
```

A successful hot copy will mark the job pod status as `Completed` and will report hot copy status as `completed`.
An example logs from _full_ hot copy can be seen below.

```
Starting full backup for database demo on processes with labels 'backup cluster0 ' ...
{
  "beginTimestamp": "2021-05-19 09:04:06",
  "coordinatorStartId": "1",
  "destinationDirectory": "/var/opt/nuodb/backup/20210519T090406/tmp/full/data",
  "endTimestamp": "2021-05-19 09:04:07",
  "hotCopyId": "e5b8b9ba-f29f-440e-8c67-1c71cef34292",
  "message": "Hot copy successfully completed",
  "stage": 3,
  "status": "completed",
  "uri": "http://nuodb.nuodb.svc:8888/api/1/databases/hotCopyStatus?coordinatorStartId=1&hotCopyId=e5b8b9ba-f29f-440e-8c67-1c71cef34292"
}
/nuodb/nuobackup/demo/cluster0/latest = 2
/nuodb/nuobackup/demo/cluster0/2 = 20210519T090406
/nuodb/nuobackup/demo/latest = cluster0
```

The above _full_ backup will be available in all _RUNNING_ HC SMs at `/var/opt/nuodb/backup/20210519T090406` with the following content:

```
$ kubectl exec -ti --namespace nuodb \
    sm-database-nuodb-cluster0-demo-hotcopy-0 -- \
    ls /var/opt/nuodb/backup/20210519T090406
full  state.xml  tmp
```

#### Backup retention management

There are no backup retention policies currently available.
The admin user is responsible to make sure that there is enough space in the backup volume of each SM for new hot copies to be created.
Customer provided `ReadWriteMany` backup volumes, such as Elastic Block Store (EBS) and Azure Files, can be used  to enable retention management using cloud tooling.

## Restore

### Restore/Import source and type

The restore/import source location can be one of the following:
- remote _URL_ in a form of `protocol://authority/path`.
The restore logic uses [cURL](https://curl.se/) to download a stream of data into an SM pod and restore the archive from it.
Standard protocols such as _HTTP/s_, _FTP_ and _SFTP_ can be used as a source.
Any other protocols that have built-in support in the shipped `cURL` binary in the NuoDB image can be used as well.
Basic authentication or an authentication token that is embedded in the URL are supported.
The remote URL should point to the `tar.gz` file which will be automatically extracted using the configured `stripLevels` setting (by default `1`). 
These are the number of leading path elements removed by the `tar` command during archive extraction so that the root data (e.g for a _stream_, the .atm files) can be accessed.
Typically, the `tar.gz` only contains one parent level folder, being the database name, for example _demo_. 

- backupset available in the backup root directory.
HC SMs have a backup volume attached and are selected for hot copy operation during the scheduled backups. Only these SMs have access to local hot copies and their archives can be restored using a local backup set.

The most recent successful backup set name can be referenced in one of the following ways:
- `:latest` - The latest backup set name among all backup groups.
- `:group-latest` - The latest backup set name for the backup group in current database deployment.

Two restore/import source types are available:
- _stream_ - Used as a synonym for an offline backup.
A remote _stream_ source is downloaded and extracted directly into the archive directory which doesn't require extra space for holding the restore/import source file.
- _backupset_ - Hot copy backup set which is either fetched from a remote source or available in the backup directory.
A remote _backupset_ requires additional space as it is downloaded and extracted temporarily in the archive volume to be used during archive restore operation.

### Fine-grained archive selection

NuoDB _restore_ chart provides several ways to select which archives should be restored during the [distributed database restore](#distributed-database-restore-from-source) or [Archive seed restore](#archive-seed-restore).

- explicitly selecting specific archive IDs - `restore.archiveIds` variable can be set to specific archive IDs in the domain state for the target database.
To list all archives for a database, the `nuocmd show archives --db-name <db>` command is used.
The process for each archive can be seen under the archive info which makes it easier for the user to mark the archives for restore.
- selecting archives served by database processes with specific labels - `restore.labels` can be used to configure process labels, which then define the archives which will be selected for a restore.
Any configured label and value that matches an SM process will add its archive to the list of selected archives for a restore.
For instance, this can be used to easily select all SMs in a specific backup group - `--set restore.labels.backup="<backup group>"`.

NuoDB will automatically assign the following process labels:

- archive-pvc - the name of the archive PVC associated to this pod (available only for SMs)
- container-id - the ID of the container
- pod-name - Kubernetes pod name
- pod-uid - Kubernetes pod UID
- backup - The name of the backup group (available only for HC SMs)
- host - Kubernetes node hostname on which the pod is running (available only for TEs)

### Distributed restore

#### Distributed database restore from source

Database in-place restore can recover from a complete loss or data corruption of all database archives by reverting the database state to a previous restore point.
NuoDB _restore_ chart is used to **overwrite** the existing database state using a configured restore source.
To perform a database in-place restore of the entire database, the `restore.type` setting should be set to `database` and the database must be shut down and restarted.
NuoDB will ensure that on the restart, the SMs serving restored archives are started first, and all other database processes will wait for the database restore to complete.
Archives that are not restored will then be synced from the running SM processes.

> **NOTE**: If the restore source was copied from another environment, then all archives need to be restored. Otherwise, the Storage Manager process will fail to start with error _Archive "/var/opt/nuodb/archive/nuodb/demo" doesn't match the database.  Expected UUID \*\*\*, got \*\*\*._
An alternative method is to create a new database deployment using [Automatic archive initial import](#automatic-archive-initial-import).
If part of doing this involves destroying an old environment, be sure to clean up the old storage volumes to avoid additional storage costs.

The database restart during the restore process is controlled by the `restore.autoRestart` setting value (by default _"true"_).
Users can retain control to manually stop and restart the database pods by disabling the auto-restart during the _restore_ chart installation.
All SMs serving restored archives will need to be started together so that the database restore process can complete.

The high-level steps to perform database in-place restore are the following:

1. Identify the restore source which can be either a backup set created by scheduled backup jobs or a remote source.
2. Identify which archives need to be restored.
By default, all archives served by HC SMs in a single Kubernetes cluster will be selected for a restore as they can define the database state by restoring from a backup set available in the backup volumes.
3. Ensure that the restore source is available on all SMs selected for restore.
4. Ensure that there is enough free disk space in the archive volume of each SM selected for restore so that it can accommodate a backup of the existing archive contents and the restored archive.
If a backup set using URL is selected as a restore source, it will be downloaded temporarily in the archive directory.
5. Invoke the database restore request by installing NuoDB _restore_ chart.
By default, the database set as `restore.target` will be **shutdown** and restarted.

Distributed database restore starts by installing the `restore` chart which will invoke the database restore request.
In this example, a specific backup set that is available on all HC SMs will be used as a restore source.

> **NOTE**: If the restore source is set to `:latest` or `:group-latest`, it's recommended that the incremental and journal backups are temporarily suspended. This is to ensure that another backup is not accidentally taken.

```bash
helm install -n nuodb restore nuodb/restore \
  --namespace nuodb \
  --set cloud.cluster.name="cluster0" \
  --set admin.domain="nuodb" \
  --set restore.target=demo \
  --set restore.source="20210219T123205"
```

Wait for the Kubernetes restore job to finish configuring the restore request for the target database.

```bash
kubectl wait \
  --for=condition=complete \
  --namespace nuodb \
  job/restore-demo
```

The log for the pod created for this job should indicate a successful restore request placement similar to the one below.

```
2021-02-19T13:06:33.110+0000 restore_type=database; restore_source=20210219T123205; arguments= --labels backup cluster0
2021-02-19T13:06:33.997+0000 restore.autoRestart=true - initiating full database restart for database demo
2021-02-19T13:06:34.735+0000 Restore job completed
```

Restore for the selected archives will be performed after the container is automatically restarted by Kubernetes and before the SM process is started.
This involves the creation of new archive metadata with a new archive ID in the NuoDB admin tier for the restored archive.
The SM pod log file will indicate a successful archive restore similar to the one below.

```
2021-02-19T13:06:36.775+0000 ===========================================
2021-02-19T13:06:36.784+0000 logsize=5085; maxlog=5000000
2021-02-19T13:06:36.837+0000 Directory /var/opt/nuodb/archive/nuodb/demo exists
2021-02-19T13:06:40.071+0000 archiveId=4; DB=demo; hostname=sm-database-nuodb-cluster0-demo-hotcopy-0
2021-02-19T13:06:41.938+0000 path=/var/opt/nuodb/archive/nuodb/demo; atoms=71; catalogs=198
2021-02-19T13:06:43.943+0000 Archive with archiveId=4 has been requested for a restore
2021-02-19T13:06:43.951+0000 Archive restore will be performed for archiveId=4, source=20210219T123205, type=backupset, strip=1
2021-02-19T13:06:43.965+0000 Restoring 20210219T123205; existing archive directores: total 16
drwxr-xr-x 33 nuodb root 4096 Feb 19 13:06 demo
drwxr-xr-x 33 nuodb root 4096 Feb 19 11:45 demo_moved
drwxr-xr-x  2 nuodb root 4096 Feb 19 10:32 demo-save-20210219T103238
drwxr-xr-x  2 nuodb root 4096 Feb 19 11:45 demo-save-20210219T114543
2021-02-19T13:06:44.043+0000 (restore) recreated /var/opt/nuodb/archive/nuodb/demo; atoms=0
2021-02-19T13:06:44.048+0000 Calling nuodocker to restore 20210219T123205 into /var/opt/nuodb/archive/nuodb/demo
2021-02-19T13:06:45.956+0000 restore: Finished restoring /var/opt/nuodb/backup/20210219T123205 to /var/opt/nuodb/archive/nuodb/demo. Created archive with archive ID 5
...
```

The archive data before restore will be saved for house-keeping in a `/var/opt/nuodb/archive/<db>-save-<timestamp>` directory.
It is recommended to move this directory after a successful database restore.

One of the SMs serving a restored archive will be selected as the database restore coordinator, which will be responsible for preparing the domain state for the database in-place restore operation.
All TEs will wait for the database restore to complete before starting, hence we can check their log entries which will show a successful database restore similar to the log snippet below:

```
2021-02-19T13:06:38.019+0000 INFO  root Waiting for database restore to complete ...
2021-02-19T13:06:56.019+0000 INFO  root Found coordinator process startId=11, hostname=sm-database-nuodb-cluster0-demo-hotcopy-1.demo.nuodb.svc.cluster.local, archiveId=6
2021-02-19T13:06:56.019+0000 INFO  root Waiting for coordinator process startId=11 to become RUNNING
2021-02-19T13:06:57.019+0000 INFO  root Database demo restored by process startId=11, address=sm-database-nuodb-cluster0-demo-hotcopy-1/172.17.0.12:48006, archiveId=6
...
```

The database restore request will be cleared after a successful database in-place restore and can be viewed using `nuodocker get restore-requests --db-name <db>`.
The database state should be manually verified by `nuocmd show domain` and using SQL queries to ensure that it's in the desired state after a successful restore.

If there is a need to modify the database restore request, a new database shutdown and restart cycle is needed so that the request can be completed.
The same applies if the reason for the previously failed database restore is corrected.
When troubleshooting distributed database restore failures, it's important to collect logs from all database processes and their corresponding pods.

> **NOTE**: In multi-cluster deployments, different backup cron jobs execute with different schedules unless they are configured to be part of the same backup group.

Backup sets produced by different backup groups can't be used together to perform _distributed database in-place restore_.
To perform a database restore in multi-cluster deployment, proceed with one of the approaches described below.

1. Select the hot copy SMs in one of the clusters, see section [Fine grained archive restore selection](#fine-grained-archive-selection), and provide a backup set available in that cluster or `:group-latest` as `restore.source`.
2. Make the selected backup set available as a remote URL in both clusters. It can be used as a restore source for all the SMs in the database only if no custom storage groups are configured.
3. Configure HC SMs in all clusters to be part of a single backup group and leave one set of backup jobs enabled in only one of the clusters to control this backup group.
This will ensure that all HC SMs backups will be coordinated during hot copy requests.

#### Distributed database restore with storage groups

Database restore using user-defined storage groups (TP/SG) is a special case of a database restore operation.

The process is documented in the [Distributed database restore from source](#distributed-database-restore-from-source) section. Several considerations need to be taken into account:

- a complete set of archives serving all storage groups must be restored when performing database in-place restore
- to ensure that backup coverage is complete, each storage group must be served by at least one HC SM

A complete set of archives can be selected using several of the methods described in [Fine-grained archive selection](#fine-grained-archive-selection) section. NuoDB won't perform any special checks during database restore to ensure that the archives selected for a restore are a complete set of archives.
If some of the storage groups are missing from the selection, their state won't be restored.

> **NOTE**: For more information about storage groups, check [Using Table Partitions and Storage Groups](https://doc.nuodb.com/nuodb/latest/database-administration/using-table-partitions-and-storage-groups/)

#### Manual database restore

_Manual database restore_ is a special case of database restore operation and allows complex restore operations to be executed easier in Kubernetes deployments.
The _Manual database restore_ operation is initiated by installing the _restore_ chart with `restore.manual="true"` which creates a manual restore request.
This mode blocks all database pods before allowing them to form a database, in order to give the user access to the persistent volumes used by an SM.
After the archive restore is marked as completed, NuoDB will unblock the processes, a restore coordinator will be selected and the restore process will continue as documented in the [Distributed database restore from source](#distributed-database-restore-from-source) section.

As an example, a Point-in-Time (PiT) restore in a new environment will be demonstrated to fix a "fat-finger" error in production.
Currently, the automatic initial archive restore doesn't support restore to a specific point in time, hence an _Manual database restore_ is used.

> **NOTE**: For more information on PiT restore and `nuoarchive`, please check [here](https://doc.nuodb.com/nuodb/latest/reference-information/command-line-tools/nuodb-archive/nuodb-archive---restoring/).

We have selected to restore from a backup set that has several journal backup elements which can be seen using `nuoarchive restore --report-backups <backup set>`:

```xml
<BackupSet id="323ba60e-f8ff-6040-438c-c16e397f52ff" database="demo" databaseId="84a3b51b-3a17-444c-f188-e4836769c12a" collectionId="772d059c-7fba-40f9-8f1c-2ec4c820a152" archiveId="b18a0698-c630-3e4a-e59c-739d86752ecd">
    <BackupElements>
        <BackupElement id="b3e78c2b-45a7-4a4c-841f-fca454d324e5" type="journal" startDate="2021-02-19 17:01:05" endDate="2021-02-19 17:01:05"/>
        <BackupElement id="7733fdd4-0a00-416b-a819-1dbfe442481d" type="journal" startDate="2021-02-19 17:00:06" endDate="2021-02-19 17:00:06"/>
        <BackupElement id="cba2d8ae-8553-49f6-8d4b-6cedddf4b18a" type="incremental" startDate="2021-02-19 17:00:08" endDate="2021-02-19 17:00:08"/>
        <BackupElement id="772d059c-7fba-40f9-8f1c-2ec4c820a152" type="full" startDate="2021-02-19 17:00:05" endDate="2021-02-19 17:00:05"/>
    </BackupElements>
</BackupSet>
```

The target for the restore will be a pre-production database that is already running.
Since the backup set used here is taken from a different environment, all database archives will need to be restored.
We are using a database with two SMs for simplicity and request both database archives for a restore.

```bash
helm install -n nuodb restore nuodb/restore \
  --namespace nuodb \
  --set cloud.cluster.name="cluster0" \
  --set admin.domain="nuodb" \
  --set restore.target=demo \
  --set restore.type=database \
  --set restore.source="http://nginx.web.svc.cluster.local/20210219T170005.tar.gz" \
  --set restore.archiveIds="{0,1}" \
  --set restore.manual=true
```

The SM process serving archive IDs 0 and 1 will block and wait for their archives to be restored which is visible in the log below.
All TEs will wait for the database restore to complete before they attempt to start.

```
...
2021-02-19T17:10:13.019+0000 INFO  root Manual restore has been requested for archiveId=0, database=demo. Waiting for archive restore to complete ...
```

Once connected to the corresponding SM pod, the archive restore is done manually by using `nuoarchive restore --report-timestamps` and `nuoarchive restore --restore-snapshot` commands.
The first command will report timestamp and transaction ID mappings available for PiT restore, and the second command will restore the snapshot identified by a transaction ID.

```
2021-02-19T17:00:30 4484
2021-02-19T17:00:31 4868
```

In this example, transaction ID 4484 will be used during the restore.

```bash
cd /var/opt/nuodb/archive/nuodb
mv demo demo-save-20210219T183055
mkdir download && cd download
curl -k  http://nginx.web.svc.cluster.local/20210219T170005.tar.gz | tar xzf - -C .

nuoarchive restore --report-timestamps $PWD/20210219T170005

nuoarchive restore --restore-snapshot 4484 --restore-dir ../demo $PWD/20210219T170005
```

After successful archive restore, its original archive ID should be marked as complete.
This will delete the archive metadata from the NuoDB Admin tier and will cause the SM to proceed with its startup operations.

```
nuodocker complete restore --db-name demo --archive-ids 0
```

Repeat the above steps for archive ID 1.

> **NOTE**: If the database has only one archive, you will need to delete the database by using `nuocmd delete database --db-name <db>` before the last archive can be removed.
`nuodocker start sm` will automatically recreate the database and the archive once the archive restore is marked complete.

The database restore request will be cleared after a successful database in-place restore and can be viewed using `nuodocker get restore-requests --db-name <db>`.
The database state should be manually verified by `nuocmd show domain` and using SQL queries to ensure that it's in the desired state after a successful restore.

### Automatic archive initial import

The _automatic archive initial import_ is a special case of a database restore operation.
This operation allows the user to specify an external URL that contains a `tar.gz` archive copy or a backup set which will be automatically downloaded and used to define the database initial state, avoiding the need to do a post-deployment database load manually.
Initial import is used for importing a backup taken from a different environment. It is configured in the `database.autoImport` section and will execute for every Storage Manager in the database deployment when all of the below conditions are met:

- `database.autoImport.source` and `database.autoImport.type` are configured
- the archive hasn't been initialized yet - `1.atm` is missing from the archive directory
- no archive metadata is found for the archive in the admin layer or on disk

To use automatic archive import in a pre-existing environment, it should be destroyed first so that the existing domain is ready to host the new database with the same name.

1. Delete the database chart deployment.
2. Remove the associated PVC, PV, and cloud storage volume that remains after database deletion.
3. Remove the database from the domain by using `nuocmd delete database --db-name <db>`.
4. List removed archives associated with the database by using `nuocmd show archives --db-name demo --removed`.
5. Remove all archives associated with the database by using `nuocmd delete archive --archive-id <id> --purge`.

NuoDB Admin tear raft PVs and PVCs can be removed instead of steps 3 to 5.

The high-level steps to configure the automatic archive import operation to bootstrap a new database are the following:

1. Prepare the import source and upload it to a remote location available to NuoDB Storage Manager pods.
2. Set `database.autoImport` variables during NuoDB database chart installation.

> **NOTE**: For simplicity _nginx_ deployment will be used in this example so that our remote source is available via HTTP in the Kubernetes cluster.

Create nginx deployment and service:

```bash
kubectl create namespace web
kubectl create deployment nginx --image=nginx --namespace web
kubectl create service clusterip nginx --tcp=80:80 --namespace web
```

Prepare and upload the NuoDB import source.
In this case, the `tag.gz` file contains NuoDB database archive copy taken from another environment.
Typically, ensure it only contains one parent level folder, being the database name, for example, _demo_. Otherwise `database.autoImport.stripLevels` should be set during database deployment.

```bash
file=demo-20201223T132332-stream.tar.gz
nginx_pod=$(kubectl get pod \
  --selector "app=nginx" \
  --namespace web \
  --output=jsonpath={.items..metadata.name})
kubectl cp "${file}" --namespace web "${nginx_pod}:/usr/share/nginx/html/${file}"
```

To avoid connectivity issues, it's recommended to verify that the remote source is available to pods in the NuoDB namespace.
We'll start a sample busybox pod to validate the URL.

```bash
kubectl run -ti busybox --rm --image=busybox --restart=Never --namespace nuodb -- \
  wget --spider http://nginx.web.svc.cluster.local/$file
```

Automatic archive import is typically used to create a new database deployment.
After installing the NuoDB admin chart, deploy the NuoDB database chart and enable automatic archive import.

```bash
helm install database nuodb/database \
    --namespace nuodb \
    -f values.yaml \
    --set database.autoImport.source="http://nginx.web.svc.cluster.local/$file" \
    --set database.autoImport.type="stream" \
    --set database.autoImport.stripLevels="1"
```

If remote authentication is needed, `database.autoImport.credentials` should also be set in a form of `user:password`.

The process of automatically downloading the remote source will require some time depending on the network speed and the size of the file.
Import progress can be monitored by requesting the log of the desired SM pod.

An example log indicating successful automatic archive import can be seen below.

```
2021-02-19T10:32:27.748+0000 ===========================================
2021-02-19T10:32:27.754+0000 logsize=73; maxlog=5000000
2021-02-19T10:32:27.765+0000 Created new dir /var/opt/nuodb/archive/nuodb/demo
2021-02-19T10:32:33.185+0000 archiveId=-1; DB=demo; hostname=sm-database-nuodb-cluster0-demo-0
2021-02-19T10:32:34.742+0000 path=/var/opt/nuodb/archive/nuodb/demo; atoms=0; catalogs=0
2021-02-19T10:32:36.395+0000 Automatic archive import will be performed for archiveId=-1, source=http://nginx.web.svc.cluster.local/demo-20201223T132332-stream.tar.gz, type=stream, strip=1
2021-02-19T10:32:36.441+0000 Restoring http://nginx.web.svc.cluster.local/demo-20201223T132332-stream.tar.gz; existing archive directores: total 4
drwxr-xr-x 2 nuodb root 4096 Feb 19 10:32 demo
2021-02-19T10:32:36.470+0000 (restore) recreated /var/opt/nuodb/archive/nuodb/demo; atoms=0
2021-02-19T10:32:36.538+0000 curl -k  http://nginx.web.svc.cluster.local/demo-20201223T132332-stream.tar.gz | tar xzf - --strip-components 1 -C /var/opt/nuodb/archive/nuodb/demo
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100 74801  100 74801    0     0  8116k      0 --:--:-- --:--:-- --:--:-- 8116k
2021-02-19T10:32:36.650+0000 restoring archive and/or clearing restored archive physical metadata
2021-02-19T10:32:38.387+0000 Finished restoring /var/opt/nuodb/archive/nuodb/demo to /var/opt/nuodb/archive/nuodb/demo. Created archive with archive ID 1
...
```

### Automatic archive restore

_Automatic archive restore_ is a special case of the seed restore operation.
It allows the SM to automatically restore its archive if it is corrupted or missing from a predefined restore source so that the accelerated sync use case is supported automatically.
Although this mechanism can be used with a remote URL, it is typically used with `:latest` or `:group-latest` configured as a restore source.

_Automatic archive restore_ is configured in the `database.autoRestore` section and will execute for every Storage Manager in the database deployment when all of the below conditions are met:

- `database.autoRestore.source` and `database.autoRestore.type` are configured
- archive metadata is found for the archive in the admin layer or on disk
- configured restore source is available in the SM pod
- the archive is detected as corrupted or missing - if the atom files count is less than 20 and catalog atom files count is less than 2

> **NOTE**: _Automatic archive restore_ with a local hot copy is used only for HC SMs.
Alternatively, a remote _URL_ containing recent backup can be configured to perform a seed restore for all SMs.

To demonstrate, we can simulate a situation where archive data in a running database for one of the HC SMs is lost.
This will cause the Storage Manager process to assert once it hits the disk and Kubernetes will automatically restart the _engine_ container.
After, the restart _automatic archive restore_ will be performed, reducing the archive sync time and returning the database back to its original capacity as soon as possible.

Configure _automatic archive restore_ during initial installation of the database chart or later on:

```bash
helm upgrade --install database nuodb/database \
    --namespace nuodb \
    -f values.yaml \
    --set database.autoRestore.source=":group-latest" \
    --set database.autoRestore.type="backupset"
```

> **NOTE**: For multi-cluster deployments `:latest` should not be used as an _automatic archive restore_ source unless the database deployments in all clusters are configured to be part of the same backup group.

Simulate a total loss of an archive for an HC SM by moving the archive directory to some other location.

```bash
kubectl exec -ti sm-database-nuodb-cluster0-demo-hotcopy-0 --namespace nuodb -- \
  mv /var/opt/nuodb/archive/nuodb/demo /var/opt/nuodb/archive/nuodb/demo_moved
```

Wait for the engine container to assert and restart. Then, check that the _automatic archive restore_ operation is performed for its archive.
If there is no SQL workload running against the database, a sample SQL query can be executed to accelerate the engine exit.
The SM pod log will indicate a successful archive restore similar to the log snippet below.

```
2021-02-19T11:45:33.182+0000 ===========================================
2021-02-19T11:45:33.239+0000 logsize=2749; maxlog=5000000
2021-02-19T11:45:33.244+0000 Created new dir /var/opt/nuodb/archive/nuodb/demo
2021-02-19T11:45:37.059+0000 archiveId=2; DB=demo; hostname=sm-database-nuodb-cluster0-demo-hotcopy-0
2021-02-19T11:45:38.444+0000 path=/var/opt/nuodb/archive/nuodb/demo; atoms=0; catalogs=0
2021-02-19T11:45:43.686+0000 Latest restore for cluster0 resolved to 20210219T114404
2021-02-19T11:45:43.692+0000 Automatic archive repair will be performed for archiveId=2, source=20210219T114404, type=backupset, strip=1
2021-02-19T11:45:43.749+0000 Restoring 20210219T114404; existing archive directores: total 12
drwxr-xr-x  2 nuodb root 4096 Feb 19 11:45 demo
drwxr-xr-x 33 nuodb root 4096 Feb 19 11:45 demo_moved
drwxr-xr-x  2 nuodb root 4096 Feb 19 10:32 demo-save-20210219T103238
2021-02-19T11:45:43.778+0000 (restore) recreated /var/opt/nuodb/archive/nuodb/demo; atoms=0
2021-02-19T11:45:43.783+0000 Calling nuodocker to restore 20210219T114404 into /var/opt/nuodb/archive/nuodb/demo
2021-02-19T11:45:45.470+0000 Finished restoring /var/opt/nuodb/backup/20210219T114404 to /var/opt/nuodb/archive/nuodb/demo. Created archive with archive ID 4
2021-02-19T11:45:46.853+0000 Purging archive metadata archiveId=2
...
```

### Archive seed restore
The _archive seed restore_ operation restores a corrupted or lost archive while the database is running.
This operation will reset the existing archive state, and once it joins the database, the new SM will sync to the current state of the running database.
For this reason, you cannot restore the entire database by restoring a single SM in a running database.
The _seed restore_ has the potential to reduce the _SYNCing_ time from other running SMs as only the changed atoms will be transferred.
The database will remain running throughout this operation.

The high-level steps to perform _archive seed restore_ are the following:

1. Identify which archives will be restored.
2. Ensure that the restore source is available either locally in the backup volume or as a remote URL.
3. Ensure that there is enough free disk space in the archive volume of each SM selected for restore so that it can accommodate a backup of the existing archive contents and the restored archive.
If a backup set using URL is selected as a restore source, it will be downloaded temporarily in the archive directory.
4. Invoke the database restore request by installing the NuoDB _restore_ chart and selecting archives for restore.
5. Start or restart the selected SM

In this example, there is an SM exited because it experienced archive data corruption.
One way to fix the issue is to completely delete the corrupted archive and let the SM sync the whole archive from one of the running processes in the database.
To reduce the sync time, we can upload one of the recent online database backups to a remote location and restore the corrupted archive.

> **NOTE**: A copy of an archive or backup set from another database can't be used to perform _archive seed restore_. Otherwise the Storage Manager process will fail to start with error _Archive "/var/opt/nuodb/archive/nuodb/demo" doesn't match database.  Expected UUID \*\*\*, got \*\*\*._

Select the archive for restore by specifying either `restore.archiveIds` or `restore.labels` and install the _restore_ chart with `restore.type="archive"`.
We are using the `pod-name` process label to select the desired SM.

```bash
helm install restore nuodb/restore \
  --namespace nuodb \
  --set cloud.cluster.name="cluster0" \
  --set admin.domain="nuodb" \
  --set restore.target=demo \
  --set restore.type=archive \
  --set restore.source="http://nginx.web.svc.cluster.local/20210219T160002.tar.gz" \
  --set restore.labels.pod-name="sm-database-nuodb-cluster0-demo-1"
```

By default the database process selected for a restore will be restarted if it is running which can be seen from the logs of the restore pod.

```
2021-03-02T11:25:44.199+0000 restore_type=archive; restore_source=http://nginx.web.svc.cluster.local/20210219T160002.tar.gz; arguments= --labels pod-name sm-database-nuodb-cluster0-demo-1
2021-03-02T11:25:46.838+0000 restore.autoRestart=true - initiating process startId=0 restart
2021-03-02T11:25:47.509+0000 Restore job completed
```

The SM will download the remote source and perform the requested restore while the rest of the database is running.
A successful restore should be seen in the log of the SM pod:

```
2021-02-19T16:25:43.586+0000 Archive with archiveId=1 has been requested for a restore
2021-02-19T16:25:43.591+0000 Archive restore will be performed for archiveId=1, source=http://nginx.web.svc.cluster.local/20210219T160002.tar.gz, type=backupset, strip=1
2021-02-19T16:25:43.600+0000 Restoring http://nginx.web.svc.cluster.local/20210219T160002.tar.gz; existing archive directores: total 8
drwxr-xr-x 33 nuodb root 4096 Feb 19 16:25 demo
drwxr-xr-x  2 nuodb root 4096 Feb 19 10:32 demo-save-20210219T103236
2021-02-19T16:25:43.682+0000 (restore) recreated /var/opt/nuodb/archive/nuodb/demo; atoms=0
2021-02-19T16:25:43.693+0000 curl -k  http://nginx.web.svc.cluster.local/20210219T160002.tar.gz | tar xzf - --strip-components 1 -C /var/opt/nuodb/archive/nuodb/20210219T160002-downloaded
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  100k  100  100k    0     0  9164k      0 --:--:-- --:--:-- --:--:-- 9164k
2021-02-19T16:25:44.285+0000 restoring archive and/or clearing restored archive physical metadata
2021-02-19T16:25:45.876+0000 Finished restoring /var/opt/nuodb/archive/nuodb/20210219T160002-downloaded to /var/opt/nuodb/archive/nuodb/demo. Created archive with archive ID 8
2021-02-19T16:25:45.879+0000 removing /var/opt/nuodb/archive/nuodb/20210219T160002-downloaded
2021-02-19T16:25:48.585+0000 Restore request for archiveId=1 marked as completed
...
```

## Troubleshooting

### Backup failure

The backup jobs execution should be monitored regularly to ensure that they complete successfully.
By default, each backup job execution will be retried on failure which can be configured by the `database.hotCopy.restartPolicy` setting.
The number of pods kept for failed backup jobs is defined by the `database.hotCopy.failureHistory` setting.
This setting is useful when debugging failed backup job execution.
The logs of the pod executing the job should be checked for errors.

For more examples, check [Examples of Hotcopy Errors](https://doc.nuodb.com/nuodb/latest/deployment-models/physical-or-vmware-environments-with-nuodb-admin/database-operations/backing-up-and-restoring-databases/monitoring-hot-copy-execution/examples-of-hot-copy-errors/).

#### No SM matching backup labels

The output below shows a failed _full_ hot copy job:

```bash
$ kubectl get pods --selector job-name --namespace nuodb
NAME                                                READY   STATUS             RESTARTS   AGE
full-hotcopy-demo-cronjob-1613661840-rgpxv          0/1     Completed          0          32m
full-hotcopy-demo-cronjob-1613663760-r8lm8          0/1     Error              1          14s
```

The logs from the pod indicate that there are no running SMs that match the provided labels.

```bash
$ kubectl logs full-hotcopy-demo-cronjob-1613663760-r8lm8 --namespace nuodb
Starting full backup for database demo on processes with labels 'backup cluster0 ' ...
'backup database' failed: No SMs found matching the provided labels backup cluster0
Error running hotcopy 1
```

> **ACTION**: Check if there are _RUNNING_ HC SMs that match the configured labels using `nuocmd show domain` and `nuocmd get processes` commands.

Further investigation shows that the HC SMs statefulset has been scaled down to 0 replicas.

```bash
$ kubectl get statefulsets.apps --namespace nuodb
NAME                                      READY   AGE
admin-nuodb-cluster0                      1/1     3h41m
sm-database-nuodb-cluster0-demo           2/2     3h41m
sm-database-nuodb-cluster0-demo-hotcopy   0/0     3h41m
```

#### Backup not completed in the specified timeout

The hot copy operation timeout can be configured in `database.sm.hotCopy.timeout` and `sm.hotCopy.journalBackup.timeout` settings which defines the number of seconds to wait for the hot copy operation to finish.
In NuoDB Helm Charts 3.2.0+ the timeout has been changed to `0` by default and backup jobs will wait forever for the operation to finish.
NuoDB doesn't cancel the underlying hot copy operation when the client timeout waiting for it to finish.

The following output shows logs from a failed backup job due to timeout.

```
Starting full backup for database demo on processes with labels 'backup cluster0 ' ...
{
  "beginTimestamp": "2021-05-20 08:35:01",
  "coordinatorStartId": "4",
  "destinationDirectory": "/var/opt/nuodb/backup/20210520T083501/tmp/full/data",
  "hotCopyId": "e9552478-9b91-4ca0-bfb3-e7e8b63c2572",
  "message": "Stage1: Begin: timeout while waiting for hot copy operation to finish",
  "stage": 1,
  "status": "running",
  "uri": "http://nuodb.nuodb.svc:8888/api/1/databases/hotCopyStatus?coordinatorStartId=4&hotCopyId=e9552478-9b91-4ca0-bfb3-e7e8b63c2572"
}
Error running hotcopy 1
```

> **ACTION**: Check the current status of the hot copy operation using `nuocmd get hotcopy-status` command by supplying the `hotCopyId` and `coordinatorStartId` taken from the output above.
Ensure that the backup timeout settings are properly configured in your environment.

#### Overlapping backups

NuoDB is preventing from running concurrent hot copy operations of the same type.
The check is done at the archive level. Multiple archives can be participating in a single hot copy operation. By default single hot copy request targets all HC SMs in a single cluster.

The following error message can be found in backup job pod logs if it fails because another hot copy operation is already running.

```
Starting full backup for database demo on processes with labels 'backup cluster0 ' ...
'backup database' failed: Failure while performing hot-copy: Error while sending engine request: Only one full hot copy at a time allowed
Error running hotcopy 1
```

> **ACTION**: No immediate action is needed as Kubernetes will restart the failed backup job based on the `database.hotCopy.restartPolicy` setting.
It's recommended to check the reason why there are concurrent backups from the same type and adjust the backup schedules if needed.

#### Backup storage full

A backup job will fail if there is no more disk space available in the backup volume for all or some of the HC SMs.
NuoDB doesn't remove old backups automatically so manual action is required to extend the backup volume size or free more disk space.

Sample output can be seen below.

```
Starting full backup for database demo on processes with labels 'backup cluster0 ' ...
{
  "beginTimestamp": "2021-05-20 09:00:54",
  "coordinatorStartId": "5",
  "destinationDirectory": "/var/opt/nuodb/backup/20210520T085520/tmp/full/data",
  "hotCopyId": "95ce5791-c582-448c-a059-52719480a564",
  "message": "Hot copy failed: Node 6: /var/opt/nuodb/backup/20210520T085520/tmp/full/data/journal/9558a2fb-ac45-a94e-edba-9af0a0f1e958/njf_0/2.njf: File read failed: No space left on device",
  "stage": 1,
  "status": "failed",
  "uri": "http://nuodb.nuodb.svc:8888/api/1/databases/hotCopyStatus?coordinatorStartId=5&hotCopyId=95ce5791-c582-448c-a059-52719480a564"
}
Error running hotcopy 1
```

> **ACTION**: Increase backup volume size or free more disk space.
When using backup sets, the `data` directories may be compressed and moved to cold storage or deleted from local storage.

It may be desirable to delete old hot copies to save space, or to comply with data retention requirements.
When using backup sets, the minimum unit of deletion is the backup set. Delete backup sets in order from oldest to newest.
If a backup set contains a transaction of interest, do not delete the prior backup set.

For more information on how to work with backup sets, please see [Using Backup Sets](https://doc.nuodb.com/nuodb/latest/deployment-models/physical-or-vmware-environments-with-nuodb-admin/database-operations/backing-up-and-restoring-databases/using-online-backup/using-backup-sets/).

### Restore failure

[Distributed database in-place restore](#distributed-database-restore-from-source) begins with requesting restore operation which is executed only after the database is restarted.
At a high level the restore failures can be classified into two groups:

1. Failures during the `restore` job execution.
These are easier to investigate by looking into the job pod logs.
2. Failures during distributed database restore operations.
To investigate the reason for the failure, the container logs from all database processes should be collected.
NuoDB SMs maintain `nuosm.log` file in `/var/log/nuodb` directory which contains historical logs from SM containers startup procedures.
The log file will be available between pod restarts if `database.sm.logPersistence.enabled` setting is enabled.

#### Invalid archive ID

Database archive IDs can be provided when requesting in-place restore.
The restore request will fail if an invalid archive ID is provided.
To investigate the reason, check the logs from the restore job pod. For example:

```bash
$ kubectl logs restore-demo-tzgfw --namespace nuodb
2021-05-20T09:41:09.187+0000 restore_type=database; restore_source=20210520T085520; arguments= --archive-ids 5
'request restore' failed: Unexpected archiveIds for database demo: 5
Restore request failed
```

> **ACTION**: Check the available database archive IDs using `nuocmd show archives` command and update the necessary values. In this specific case the `restore.archiveIds` setting much be corrected.
Delete the old `restore` chart release and install new one with the updated Helm values file.

#### Invalid restore source

The selected restore source should be available to all SMs requested for a restore.
Monitor the database processes pods status after the database has been restarted to check for any failures.
In the bellow output, an error is seen on _sm-database-nuodb-cluster0-demo-hotcopy-1_ pod.

```
NAME                                                READY   STATUS             RESTARTS   AGE
sm-database-nuodb-cluster0-demo-hotcopy-1           0/1     Error              2          12m
```

Checking the pod logs shows that the archive served by this SM has been requested for a restore using restore source _20210520T085520_ which is not available in its backup volume.

```bash
kubectl logs full-hotcopy-demo-cronjob-1613663760-r8lm8
Starting full backup for database demo on processes with labels 'backup cluster0 ' ...
'backup database' failed: No SMs found matching the provided labels backup cluster0
Error running hotcopy 1
```

$ kubectl logs sm-database-nuodb-cluster0-demo-hotcopy-1
2021-05-20T09:52:15.198+0000 ===========================================
2021-05-20T09:52:15.234+0000 logsize=8154; maxlog=5000000
2021-05-20T09:52:15.236+0000 Directory /var/opt/nuodb/archive/nuodb/demo exists
2021-05-20T09:52:18.547+0000 archiveId=2; DB=demo; hostname=sm-database-nuodb-cluster0-demo-hotcopy-1
2021-05-20T09:52:20.283+0000 path=/var/opt/nuodb/archive/nuodb/demo; atoms=68; catalogs=183
2021-05-20T09:52:21.982+0000 Archive with archiveId=2 has been requested for a restore
2021-05-20T09:52:21.995+0000 Archive restore will be performed for archiveId=2, source=20210520T085520, type=backupset, strip=1
2021-05-20T09:52:22.010+0000 Error while performing restore for archiveId=2: Backupset 20210520T085520 cannot be found in /var/opt/nuodb/backup
```

The pod status will soon transition to `CrashLoopBackOff` and all other database processes will wait until the database restore completes which will never happen.

Sample output including current database pods can be seen below.

```
NAME                                                READY   STATUS             RESTARTS   AGE
restore-demo-d7jwn                                  0/1     Completed          0          10m
sm-database-nuodb-cluster0-demo-0                   0/1     Running            1          21m
sm-database-nuodb-cluster0-demo-hotcopy-0           0/1     Running            1          20m
sm-database-nuodb-cluster0-demo-hotcopy-1           0/1     CrashLoopBackOff   6          21m
te-database-nuodb-cluster0-demo-6c47c5c696-r5dg2    0/1     Running            1          31m
```

Example logs from another SM which was not selected for a restore:

```bash
$ kubectl logs sm-database-nuodb-cluster0-demo-0 --namespace nuodb
2021-05-20T09:52:15.276+0000 ===========================================
2021-05-20T09:52:15.283+0000 logsize=9821; maxlog=5000000
2021-05-20T09:52:15.288+0000 Directory /var/opt/nuodb/archive/nuodb/demo exists
2021-05-20T09:52:18.536+0000 archiveId=0; DB=demo; hostname=sm-database-nuodb-cluster0-demo-0
2021-05-20T09:52:20.272+0000 path=/var/opt/nuodb/archive/nuodb/demo; atoms=68; catalogs=183
2021-05-20T09:52:25.020+0000 INFO  root Waiting for database restore to complete ...
```

> **ACTION**: Ensure that the selected backup set is available on all SMs requested for restore.
Once the restore source is verified and updated, delete the old `restore` chart release and install a new one with the updated Helm values file.
Delete all pods for database processes using `kubectl delete pods` command so that Kubernetes restart them at once.
