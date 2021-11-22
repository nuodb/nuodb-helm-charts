# Changelog

## [v3.4.0](https://github.com/nuodb/nuodb-helm-charts/tree/v3.4.0) (2021-12-03)

[Full Changelog](https://github.com/nuodb/nuodb-helm-charts/compare/v3.3.0...v3.4.0)

**Implemented enhancements:**

- Added functionality to enable the database protocol to be upgraded automatically, which is supported in version 4.2.3 and above of the NuoDB image [\#243](https://github.com/nuodb/nuodb-helm-charts/pull/243) ([sivanov-nuodb](https://github.com/sivanov-nuodb))
- Updated the process readiness probe to be more scalable as the number of database processes grows [\#252](https://github.com/nuodb/nuodb-helm-charts/pull/252) ([adriansuarez](https://github.com/adriansuarez))
- Updated init containers to perform recursive chmod only if files are encountered without expected permissions [\#260](https://github.com/nuodb/nuodb-helm-charts/pull/260) ([adriansuarez](https://github.com/adriansuarez))

**Fixed bugs:**

- Updated nuosm to resurrect archive object only if it will be used [\#256](https://github.com/nuodb/nuodb-helm-charts/pull/256) ([sivanov-nuodb](https://github.com/sivanov-nuodb))

**Merged pull requests:**

- Increased leaderAssignmentTimeout to account for SMs going into CrashLoopBackOff state [\#242](https://github.com/nuodb/nuodb-helm-charts/pull/242) ([adriansuarez](https://github.com/adriansuarez))
- Added support for external SQL clients, which requires version 4.2.4 and above of the NuoDB image [\#254](https://github.com/nuodb/nuodb-helm-charts/pull/254) ([sivanov-nuodb](https://github.com/sivanov-nuodb))

## [v3.3.0](https://github.com/nuodb/nuodb-helm-charts/tree/v3.3.0) (2021-08-18)

[Full Changelog](https://github.com/nuodb/nuodb-helm-charts/compare/v3.2.0...v3.3.0)

**Implemented enhancements:**

- Added warning when requested backup type does not match target backup set, e.g. stream vs hotcopy [\#240](https://github.com/nuodb/nuodb-helm-charts/pull/240) ([sivanov-nuodb](https://github.com/sivanov-nuodb))
- Replaced deprecated failure-domain.beta.kubernetes.io/zone label with topology.kubernetes.io/zone in the Transparent Huge Page Chart [\#226](https://github.com/nuodb/nuodb-helm-charts/pull/226) ([mkysel](https://github.com/mkysel))
- Added a facility to send additional storage manager logs to the admin pod to improve runtime debugging [\#225](https://github.com/nuodb/nuodb-helm-charts/pull/225) ([sivanov-nuodb](https://github.com/sivanov-nuodb))
- Set default admin setting evictUnknownProcesses=true in admin/values.yaml to handle partial network disconnects between admin and engines gracefully [\#219](https://github.com/nuodb/nuodb-helm-charts/pull/219) ([adriansuarez](https://github.com/adriansuarez))
- Added the option to mount the Storage Manager journal directory on a different persistent storage volume [\#218](https://github.com/nuodb/nuodb-helm-charts/pull/218) ([mkysel](https://github.com/mkysel))
- Added functionality to recursively change permissions in admin init container to prevent runtime write permissions errors [\#214](https://github.com/nuodb/nuodb-helm-charts/pull/214) ([adriansuarez](https://github.com/adriansuarez))
- Override default admin setting thrift.message.max in admin/values.yaml to enable reading of large messages [\#213](https://github.com/nuodb/nuodb-helm-charts/pull/213) ([adriansuarez](https://github.com/adriansuarez))

**Fixed bugs:**

- Added validation of boolean helm variables in the values.yaml files. The new functionality will emit a warning if an invalid value is passed [\#223](https://github.com/nuodb/nuodb-helm-charts/pull/223) ([sivanov-nuodb](https://github.com/sivanov-nuodb))
- Removed an incorrect admin operability check from nuosm script which prevented storage manager container startup when some admin containers were down [\#217](https://github.com/nuodb/nuodb-helm-charts/pull/217) ([sivanov-nuodb](https://github.com/sivanov-nuodb))

**Removed:**

- Removed option to start Storage Managers as DaemonSets [\#222](https://github.com/nuodb/nuodb-helm-charts/pull/222) ([mkysel](https://github.com/mkysel))

**Closed issues:**

- Longer readiness timeout in admin chart [\#198](https://github.com/nuodb/nuodb-helm-charts/issues/198)

## [v3.2.0](https://github.com/nuodb/nuodb-helm-charts/tree/v3.2.0) (2021-05-14)

[Full Changelog](https://github.com/nuodb/nuodb-helm-charts/compare/v3.1.0...v3.2.0)

**Implemented enhancements:**

- Added helm template value validation for database.sm.hotCopy.enableBackups [\#204](https://github.com/nuodb/nuodb-helm-charts/pull/204) ([sivanov-nuodb](https://github.com/sivanov-nuodb))
- Added helm template value validation for restore.source. Accepted values are ':latest', ':group-latest' or any valid URL [\#202](https://github.com/nuodb/nuodb-helm-charts/pull/202) ([sivanov-nuodb](https://github.com/sivanov-nuodb))
- Added helm template value validation for database.autoRestore.type. Accepted types are either 'stream' or 'backupset' [\#199](https://github.com/nuodb/nuodb-helm-charts/pull/199) ([mkysel](https://github.com/mkysel))
- Changed the default timeout for all backup jobs from 30 min to infinite [\#197](https://github.com/nuodb/nuodb-helm-charts/pull/197) ([sivanov-nuodb](https://github.com/sivanov-nuodb))
- Upgraded kiwigrid/sidecar NuoDB Collector dependency to 1.10.8 [\#193](https://github.com/nuodb/nuodb-helm-charts/pull/193) ([butson](https://github.com/butson))
- Added database in-place restore support based on new NuoDB 4.2+ functionality. Fine-graned selection of archive ids to restore. Database in-place restore with storage group. Manual restore option and many more... [\#184](https://github.com/nuodb/nuodb-helm-charts/pull/184) ([sivanov-nuodb](https://github.com/sivanov-nuodb))

**Fixed bugs:**

- Updated Chart.yaml icon image references from nuodb.com to GitHub [\#203](https://github.com/nuodb/nuodb-helm-charts/pull/203) ([mkysel](https://github.com/mkysel))
- Set levels to strip off path names when unpacking a TAR file of an archive or backup set via the restore.stripLevels option [\#195](https://github.com/nuodb/nuodb-helm-charts/pull/195) ([sivanov-nuodb](https://github.com/sivanov-nuodb))

**Merged pull requests:**

- Updated default version to use the NuoDB image 4.2.1 [\#207](https://github.com/nuodb/nuodb-helm-charts/pull/207) ([mkysel](https://github.com/mkysel))

## [v3.1.0](https://github.com/nuodb/nuodb-helm-charts/tree/v3.1.0) (2021-02-08)

[Full Changelog](https://github.com/nuodb/nuodb-helm-charts/compare/v3.0.0...v3.1.0)

**Implemented enhancements:**

- Made database backup job/cronjob restart policy configurable and switched from Never to OnFailure [\#181](https://github.com/nuodb/nuodb-helm-charts/pull/181) ([mkysel](https://github.com/mkysel))
- Added evicted-servers configuration to admin StatefulSet used to restore majority due to a catastrophic loss of admin servers [\#180](https://github.com/nuodb/nuodb-helm-charts/pull/180) ([mkysel](https://github.com/mkysel))
- Replaced initial backup and post-restore jobs with prerequisites in backup cron jobs to streamline the database helm upgrade process [\#179](https://github.com/nuodb/nuodb-helm-charts/pull/179) ([sivanov-nuodb](https://github.com/sivanov-nuodb))
- Added Transparent Data Encryption \(TDE\) support for NuoDB Storage Manager database pods [\#168](https://github.com/nuodb/nuodb-helm-charts/pull/168) ([sivanov-nuodb](https://github.com/sivanov-nuodb))
- Improved admin readiness probes by using a single-admin "nuocmd check server" command available in NuoDB 4.1.2+ [\#166](https://github.com/nuodb/nuodb-helm-charts/pull/166) ([adriansuarez](https://github.com/adriansuarez))

**Fixed bugs:**

- Fixed failure in incremental hotcopy due to missing full backup element by re-scheduling failed full hotcopy [\#182](https://github.com/nuodb/nuodb-helm-charts/pull/182) ([sivanov-nuodb](https://github.com/sivanov-nuodb))
- Fixed an intermittent timing issue during concurrent restore of multiple Storage Managers [\#176](https://github.com/nuodb/nuodb-helm-charts/pull/176) ([sivanov-nuodb](https://github.com/sivanov-nuodb))
- Improved database backup and restore behavior. This improvement guarantees that the newest created backup set will be used for journal hotcopy instead of the latest successful one [\#173](https://github.com/nuodb/nuodb-helm-charts/pull/173) ([sivanov-nuodb](https://github.com/sivanov-nuodb))
- The nuobackup script has been enhanced so that it can wait for a certain number of SMs with requested labels to become RUNNING before performing a backup [\#172](https://github.com/nuodb/nuodb-helm-charts/pull/172) ([sivanov-nuodb](https://github.com/sivanov-nuodb))

**Deprecated:**

- Made the admin LoadBalancer job optional. This streamlines the helm upgrade process. This legacy feature has been superseded by Kubernetes Aware Admin [\#177](https://github.com/nuodb/nuodb-helm-charts/pull/177) ([mkysel](https://github.com/mkysel))

**Removed:**

- Removed obsolete and unused Red Hat OpenShift flag from the database chart [\#175](https://github.com/nuodb/nuodb-helm-charts/pull/175) ([mkysel](https://github.com/mkysel))

**Merged pull requests:**

- Bumped NuoDB Version to 4.0.8 [\#169](https://github.com/nuodb/nuodb-helm-charts/pull/169) ([mkysel](https://github.com/mkysel))

## [v3.0.0](https://github.com/nuodb/nuodb-helm-charts/tree/v3.0.0) (2020-11-09)

[Full Changelog](https://github.com/nuodb/nuodb-helm-charts/compare/v2.4.1...v3.0.0)

**Implemented enhancements:**

- Added NuoDB Collector support for database statistics collection and visual monitoring [\#161](https://github.com/nuodb/nuodb-helm-charts/pull/161) ([sivanov-nuodb](https://github.com/sivanov-nuodb))
- Moved custom NuoDB admin podAnnotations from StatefulSet metadata to the admin pod itself [\#156](https://github.com/nuodb/nuodb-helm-charts/pull/156) ([acabrele](https://github.com/acabrele))
- Simplified the required configuration changes for NuoDB admin domains not utilizing TLS network encryption [\#155](https://github.com/nuodb/nuodb-helm-charts/pull/155) ([sivanov-nuodb](https://github.com/sivanov-nuodb))
- To support Red Hat OpenShift - Added "nuodb" service account for the admin's Load Balancer job [\#153](https://github.com/nuodb/nuodb-helm-charts/pull/153) ([kmabda](https://github.com/kmabda))
- Added the ability to pass through custom annotations \(`podAnnotations`\) to be applied to pods to enable 3rd party integrations \(Vault, CNIs, ...\) [\#149](https://github.com/nuodb/nuodb-helm-charts/pull/149) ([acabrele](https://github.com/acabrele))
- Changed all database container names to "engine" inside of the StatefulSet and Deployment pods [\#135](https://github.com/nuodb/nuodb-helm-charts/pull/135) ([adriansuarez](https://github.com/adriansuarez))
- Replaced Transaction Engine \(TE\) Load Balancer job \(lbConfig\) with Load Balancer Specification via Kubernetes object annotations \(depends on NuoDB Kubernetes Aware Admin\) [\#133](https://github.com/nuodb/nuodb-helm-charts/pull/133) ([sivanov-nuodb](https://github.com/sivanov-nuodb))
- Switched to Helm 3 testing by default [\#129](https://github.com/nuodb/nuodb-helm-charts/pull/129) ([mkysel](https://github.com/mkysel))
- Made StorageClass Persistent Volume reclaim policy configurable [\#128](https://github.com/nuodb/nuodb-helm-charts/pull/128) ([mkysel](https://github.com/mkysel))

**Fixed bugs:**

- Fixed an intermittent timing issue in NuoDB Storage Manager restore-from-backup procedure by adding a timeout-and-retry [\#160](https://github.com/nuodb/nuodb-helm-charts/pull/160) ([NikTJ777](https://github.com/NikTJ777))
- Fixed NuoDB database Storage Manager StatefulSet database.autoImport credentials by removing single-quotes from curl\_creds [\#157](https://github.com/nuodb/nuodb-helm-charts/pull/157) ([NikTJ777](https://github.com/NikTJ777))
- Added option to use embedded HTTP path security credentials in the database.autoImport feature [\#154](https://github.com/nuodb/nuodb-helm-charts/pull/154) ([NikTJ777](https://github.com/NikTJ777))
- Improved retention and handling of multiple post-restart diagnostics files or folders in $NUODB\_CRASHDIR [\#147](https://github.com/nuodb/nuodb-helm-charts/pull/147) ([acabrele](https://github.com/acabrele))

**Merged pull requests:**

- Bumped NuoDB version to 4.0.7 [\#144](https://github.com/nuodb/nuodb-helm-charts/pull/144) ([mkysel](https://github.com/mkysel))

## [v2.4.1](https://github.com/nuodb/nuodb-helm-charts/tree/v2.4.1) (2020-08-30)

[Full Changelog](https://github.com/nuodb/nuodb-helm-charts/compare/v2.4.0...v2.4.1)

**Fixed bugs:**

- Forbidden!Configured service account doesn't have access. Service account may have been revoked. daemonsets.apps is forbidden [\#140](https://github.com/nuodb/nuodb-helm-charts/issues/140)

## [v2.4.0](https://github.com/nuodb/nuodb-helm-charts/tree/v2.4.0) (2020-07-15)

[Full Changelog](https://github.com/nuodb/nuodb-helm-charts/compare/v2.3.1...v2.4.0)

**Implemented enhancements:**

- Added the bootstrapServers label to the Admin StatefulSet to enable future releases of NuoDB to handle catastrophic loss of the admin-0 container [\#126](https://github.com/nuodb/nuodb-helm-charts/pull/126) ([adriansuarez](https://github.com/adriansuarez))
- Added nodeSelector, affinity and toleration parameters to the LoadBalancer Policy Job in the Admin chart in order to run an Job on a specific cluster node. [\#122](https://github.com/nuodb/nuodb-helm-charts/pull/122) ([kmabda](https://github.com/kmabda))
- Made readiness probe timeouts configurable in Admin and Database values.yaml files [\#119](https://github.com/nuodb/nuodb-helm-charts/pull/119) ([acabrele](https://github.com/acabrele))
- Added the sm.hotCopy.enableBackups value to the Database helm chart to make the build-in backup scheduling mechanism optional [\#118](https://github.com/nuodb/nuodb-helm-charts/pull/118) ([kmabda](https://github.com/kmabda))
- Allowed a ReadWriteMany PVC for persistent log volumes in Transaction Engine Deployments [\#114](https://github.com/nuodb/nuodb-helm-charts/pull/114) ([mkysel](https://github.com/mkysel))
- Added configurable number of replicas to YCSB [\#95](https://github.com/nuodb/nuodb-helm-charts/pull/95) ([mkysel](https://github.com/mkysel))
- Changed pullPolicy from Always to IfNotPresent for nuodb and busybox images [\#90](https://github.com/nuodb/nuodb-helm-charts/pull/90) ([adriansuarez](https://github.com/adriansuarez))
- Added the "leases" resource to the NuoDB role to coordinate updates to the NuoDB Admin tier that are generated by Kubernetes events [\#83](https://github.com/nuodb/nuodb-helm-charts/pull/83) ([adriansuarez](https://github.com/adriansuarez))

**Fixed bugs:**

- Changed hotcopy-{{ .Values.database.name }}-job-initial backup policy from Never to OnFailure to prevent accumulation of failed jobs [\#123](https://github.com/nuodb/nuodb-helm-charts/pull/123) ([mkysel](https://github.com/mkysel))
- Fixed Helm Operator Scorecard requirement of no empty blocks [\#104](https://github.com/nuodb/nuodb-helm-charts/pull/104) ([mkysel](https://github.com/mkysel))
- Removed cloud.sh autotagging of cloud environments to enable Helm Operator YCSB scaling [\#94](https://github.com/nuodb/nuodb-helm-charts/pull/94) ([mkysel](https://github.com/mkysel))

**Removed:**

- Removed monitoring-insights from Incubator [\#107](https://github.com/nuodb/nuodb-helm-charts/pull/107) ([mkysel](https://github.com/mkysel))
- Removed memoryOption and replaced it with resources.requests.memory [\#106](https://github.com/nuodb/nuodb-helm-charts/pull/106) ([mkysel](https://github.com/mkysel))
- Removed obsolete charts from Incubator \(backup, restore-copy, monitoring-influx\) [\#99](https://github.com/nuodb/nuodb-helm-charts/pull/99) ([mkysel](https://github.com/mkysel))
- Remove OpenShift specific TE DeploymentConfig [\#66](https://github.com/nuodb/nuodb-helm-charts/pull/66) ([mkysel](https://github.com/mkysel))

**Merged pull requests:**

- Bump NuoDB version to 4.0.6 [\#125](https://github.com/nuodb/nuodb-helm-charts/pull/125) ([mkysel](https://github.com/mkysel))
- Bump NuoDB version to 4.0.5 [\#92](https://github.com/nuodb/nuodb-helm-charts/pull/92) ([mkysel](https://github.com/mkysel))

## [v2.3.1](https://github.com/nuodb/nuodb-helm-charts/tree/v2.3.1) (2020-05-12)

[Full Changelog](https://github.com/nuodb/nuodb-helm-charts/compare/v2.3.0...v2.3.1)

**Fixed bugs:**

- Fixed an issue preventing Helm Chart release upgrade [\#111](https://github.com/nuodb/nuodb-helm-charts/issues/111)

## [v2.3.0](https://github.com/nuodb/nuodb-helm-charts/tree/v2.3.0) (2020-03-24)

[Full Changelog](https://github.com/nuodb/nuodb-helm-charts/compare/v2.2.0...v2.3.0)

**Implemented enhancements:**

- Add the "watch" verb to the NuoDB role to enable the Admin to register event listeners. [\#78](https://github.com/nuodb/nuodb-helm-charts/pull/78) ([adriansuarez](https://github.com/adriansuarez))
- Auto-configure role and service account in values files to enable NuoDB control plane synchronization with Kubernetes. [\#77](https://github.com/nuodb/nuodb-helm-charts/pull/77) ([adriansuarez](https://github.com/adriansuarez))
- Added the use of Security Context Constraints\(SCC\) for Red Hat OpenShift deployments. [\#75](https://github.com/nuodb/nuodb-helm-charts/pull/75) ([mkysel](https://github.com/mkysel))
- Increased engine readiness probe timeout value from 1 to 5 seconds. [\#74](https://github.com/nuodb/nuodb-helm-charts/pull/74) ([mkysel](https://github.com/mkysel))
- Add customization to database service names to correspond to admin services. [\#62](https://github.com/nuodb/nuodb-helm-charts/pull/62) ([NikTJ777](https://github.com/NikTJ777))
- Allow a user to override admin services suffixes to customize ClusterIP and LoadBalancer names. [\#57](https://github.com/nuodb/nuodb-helm-charts/pull/57) ([kmabda](https://github.com/kmabda))
- Added optional persistent volumes for log collection. [\#52](https://github.com/nuodb/nuodb-helm-charts/pull/52) ([acabrele](https://github.com/acabrele))
- Add K8S cluster domain support, enabling CNI based multi cluster configurations [\#51](https://github.com/nuodb/nuodb-helm-charts/pull/51) ([acabrele](https://github.com/acabrele))

## [v2.2.0](https://github.com/nuodb/nuodb-helm-charts/tree/v2.2.0) (2020-01-27)

[Full Changelog](https://github.com/nuodb/nuodb-helm-charts/compare/v2.2...v2.2.0)

## [v2.2](https://github.com/nuodb/nuodb-helm-charts/tree/v2.2) (2020-01-15)

[Full Changelog](https://github.com/nuodb/nuodb-helm-charts/compare/v2.1...v2.2)

**Implemented enhancements:**

- Add option to use SFTP source for autoImport & autoRestore and define the type of the source [\#46](https://github.com/nuodb/nuodb-helm-charts/pull/46) ([acabrele](https://github.com/acabrele))
- Enhanced backup and restore including the addition of PiT backup [\#38](https://github.com/nuodb/nuodb-helm-charts/pull/38) ([NikTJ777](https://github.com/NikTJ777))
- Provide an option to disable DB Services for direct TE connections [\#33](https://github.com/nuodb/nuodb-helm-charts/pull/33) ([acabrele](https://github.com/acabrele))
- Move database monitoring to incubator and enhanced it to work with a VPN based multi-cluster  [\#31](https://github.com/nuodb/nuodb-helm-charts/pull/31) ([acabrele](https://github.com/acabrele))
- Move YCSB to incubator and enhanced it to allow use with a VPN based multi-cluster [\#29](https://github.com/nuodb/nuodb-helm-charts/pull/29) ([acabrele](https://github.com/acabrele))
- Add local-storage StorageClass to the StorageClass chart [\#28](https://github.com/nuodb/nuodb-helm-charts/pull/28) ([mkysel](https://github.com/mkysel))
- Add Cluster IP Services [\#27](https://github.com/nuodb/nuodb-helm-charts/pull/27) ([acabrele](https://github.com/acabrele))
- \[DB-29171\] Feature request for user provided environment variables for database chart [\#26](https://github.com/nuodb/nuodb-helm-charts/pull/26) ([sivanov-nuodb](https://github.com/sivanov-nuodb))
- Make Load Balancer services optional and add ability to provision internal cloud IPs [\#23](https://github.com/nuodb/nuodb-helm-charts/pull/23) ([acabrele](https://github.com/acabrele))

**Fixed bugs:**

- Small fixes for backup and restore [\#45](https://github.com/nuodb/nuodb-helm-charts/pull/45) ([acabrele](https://github.com/acabrele))
- Fix image pull policy reference [\#44](https://github.com/nuodb/nuodb-helm-charts/pull/44) ([acabrele](https://github.com/acabrele))
- Backup & restore, permissions and pull secrets fixes [\#42](https://github.com/nuodb/nuodb-helm-charts/pull/42) ([acabrele](https://github.com/acabrele))

**Merged pull requests:**

- Bump NuoDB version to 4.0.4 [\#47](https://github.com/nuodb/nuodb-helm-charts/pull/47) ([mkysel](https://github.com/mkysel))
- Bump NuoDB version to 4.0.3 [\#43](https://github.com/nuodb/nuodb-helm-charts/pull/43) ([mkysel](https://github.com/mkysel))
- Pin NuoDB version to 4.0.2 [\#35](https://github.com/nuodb/nuodb-helm-charts/pull/35) ([mkysel](https://github.com/mkysel))
- Downgrade to Helm v2.9.0 [\#25](https://github.com/nuodb/nuodb-helm-charts/pull/25) ([sivanov-nuodb](https://github.com/sivanov-nuodb))

## [v2.1](https://github.com/nuodb/nuodb-helm-charts/tree/v2.1) (2019-11-06)

[Full Changelog](https://github.com/nuodb/nuodb-helm-charts/compare/v2.0...v2.1)

**Implemented enhancements:**

- Add readiness probes for NuoDB engine StatefulSet/Deployment [\#18](https://github.com/nuodb/nuodb-helm-charts/pull/18) ([mkysel](https://github.com/mkysel))
- \[DB-28964\] Externalize config files using generalized idiomatic approaches. [\#13](https://github.com/nuodb/nuodb-helm-charts/pull/13) ([rbuck](https://github.com/rbuck))
- \[DB-28838\] Make nuoadmin configuration options available in helm values [\#12](https://github.com/nuodb/nuodb-helm-charts/pull/12) ([acabrele](https://github.com/acabrele))
- Add more state checks to Admin StatefulSet readiness probe [\#8](https://github.com/nuodb/nuodb-helm-charts/pull/8) ([mkysel](https://github.com/mkysel))
- \[DB-28313\] Enable passing of certificates directly to the engine [\#4](https://github.com/nuodb/nuodb-helm-charts/pull/4) ([mkysel](https://github.com/mkysel))

**Fixed bugs:**

- \[DB-28733\] add missing volumeMount to the THP chart [\#5](https://github.com/nuodb/nuodb-helm-charts/pull/5) ([mkysel](https://github.com/mkysel))
- \[DB-28712\] Fix Admin Resources [\#2](https://github.com/nuodb/nuodb-helm-charts/pull/2) ([mkysel](https://github.com/mkysel))

**Closed issues:**

- Error in admin helm chart? [\#1](https://github.com/nuodb/nuodb-helm-charts/issues/1)



\* *This Changelog was automatically generated by [github_changelog_generator](https://github.com/github-changelog-generator/github-changelog-generator)*
