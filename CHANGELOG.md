# Changelog

## [3.0.0](https://github.com/nuodb/nuodb-helm-charts/tree/3.0.0) (2020-11-09)

[Full Changelog](https://github.com/nuodb/nuodb-helm-charts/compare/v2.4.1...3.0.0)

**Implemented enhancements:**

- Added NuoDB Collector support for database statistics collection and visual monitoring [\#161](https://github.com/nuodb/nuodb-helm-charts/pull/161) ([sivanov-nuodb](https://github.com/sivanov-nuodb))
- Moved custom NuoDB admin podAnnotations from StatefulSet metadata to the admin pod itself [\#156](https://github.com/nuodb/nuodb-helm-charts/pull/156) ([acabrele](https://github.com/acabrele))
- Simplified the required configuration changes for NuoDB admin domains not utilizing TLS network encryption [\#155](https://github.com/nuodb/nuodb-helm-charts/pull/155) ([sivanov-nuodb](https://github.com/sivanov-nuodb))
- To support Red Hat OpenShift - Added "nuodb" service account for the admin's Load Balancer job [\#153](https://github.com/nuodb/nuodb-helm-charts/pull/153) ([kmabda](https://github.com/kmabda))
- Added the ability to pass through custom annotations \(`podAnnotations`\) to be applied to pods to enable 3rd party integrations \(Vault, CNIs, ...\) [\#149](https://github.com/nuodb/nuodb-helm-charts/pull/149) ([acabrele](https://github.com/acabrele))
- Changed all database container names to "engine" inside of the StatefulSet and Deployment pods [\#135](https://github.com/nuodb/nuodb-helm-charts/pull/135) ([adriansuarez](https://github.com/adriansuarez))
- Replaced Transaction Engine \(TE\) Load Balancer job \(lbConfig\) with Load Balancer Specification via Kubernetes object annotations \(depends on NuoDB Kubernetes Aware Admin\) [\#133](https://github.com/nuodb/nuodb-helm-charts/pull/133) ([sivanov-nuodb](https://github.com/sivanov-nuodb))
- Switched to Helm 3 testing by default [\#129](https://github.com/nuodb/nuodb-helm-charts/pull/129) ([vegichan](https://github.com/vegichan))
- Made StorageClass Persistent Volume reclaim policy configurable [\#128](https://github.com/nuodb/nuodb-helm-charts/pull/128) ([vegichan](https://github.com/vegichan))

**Fixed bugs:**

- Fixed an intermittent timing issue in NuoDB Storage Manager restore-from-backup procedure by adding a timeout-and-retry [\#160](https://github.com/nuodb/nuodb-helm-charts/pull/160) ([NikTJ777](https://github.com/NikTJ777))
- Fixed NuoDB database Storage Manager StatefulSet database.autoImport credentials by removing single-quotes from curl\_creds [\#157](https://github.com/nuodb/nuodb-helm-charts/pull/157) ([NikTJ777](https://github.com/NikTJ777))
- Added option to use embedded HTTP path security credentials in the database.autoImport feature [\#154](https://github.com/nuodb/nuodb-helm-charts/pull/154) ([NikTJ777](https://github.com/NikTJ777))
- Improved retention and handling of multiple post-restart diagnostics files or folders in $NUODB\_CRASHDIR [\#147](https://github.com/nuodb/nuodb-helm-charts/pull/147) ([acabrele](https://github.com/acabrele))

**Merged pull requests:**

- Bumped NuoDB version to 4.0.7 [\#144](https://github.com/nuodb/nuodb-helm-charts/pull/144) ([vegichan](https://github.com/vegichan))

## [v2.4.1](https://github.com/nuodb/nuodb-helm-charts/tree/v2.4.1) (2020-08-30)

[Full Changelog](https://github.com/nuodb/nuodb-helm-charts/compare/v2.4.0...v2.4.1)

**Implemented enhancements:**

- Added DaemonSets to list of NuoDB Kubernetes Aware Admin permissions [\#141](https://github.com/nuodb/nuodb-helm-charts/pull/141) ([vegichan](https://github.com/vegichan))

**Fixed bugs:**

- Forbidden!Configured service account doesn't have access. Service account may have been revoked. daemonsets.apps is forbidden [\#140](https://github.com/nuodb/nuodb-helm-charts/issues/140)

**Merged pull requests:**

- Bump NuoDB version to 4.0.7 \[2.4-dev branch\] [\#145](https://github.com/nuodb/nuodb-helm-charts/pull/145) ([vegichan](https://github.com/vegichan))

## [v2.4.0](https://github.com/nuodb/nuodb-helm-charts/tree/v2.4.0) (2020-07-15)

[Full Changelog](https://github.com/nuodb/nuodb-helm-charts/compare/v2.3.1...v2.4.0)

**Implemented enhancements:**

- Added the bootstrapServers label to the Admin StatefulSet to enable future releases of NuoDB to handle catastrophic loss of the admin-0 container [\#126](https://github.com/nuodb/nuodb-helm-charts/pull/126) ([adriansuarez](https://github.com/adriansuarez))
- Added nodeSelector, affinity and toleration parameters to the LoadBalancer Policy Job in the Admin chart in order to run an Job on a specific cluster node. [\#122](https://github.com/nuodb/nuodb-helm-charts/pull/122) ([kmabda](https://github.com/kmabda))
- Made readiness probe timeouts configurable in Admin and Database values.yaml files [\#119](https://github.com/nuodb/nuodb-helm-charts/pull/119) ([acabrele](https://github.com/acabrele))
- Added the sm.hotCopy.enableBackups value to the Database helm chart to make the build-in backup scheduling mechanism optional [\#118](https://github.com/nuodb/nuodb-helm-charts/pull/118) ([kmabda](https://github.com/kmabda))
- Allowed a ReadWriteMany PVC for persistent log volumes in Transaction Engine Deployments [\#114](https://github.com/nuodb/nuodb-helm-charts/pull/114) ([vegichan](https://github.com/vegichan))
- Added configurable number of replicas to YCSB [\#95](https://github.com/nuodb/nuodb-helm-charts/pull/95) ([vegichan](https://github.com/vegichan))
- Changed pullPolicy from Always to IfNotPresent for nuodb and busybox images [\#90](https://github.com/nuodb/nuodb-helm-charts/pull/90) ([adriansuarez](https://github.com/adriansuarez))
- Added the "leases" resource to the NuoDB role to coordinate updates to the NuoDB Admin tier that are generated by Kubernetes events [\#83](https://github.com/nuodb/nuodb-helm-charts/pull/83) ([adriansuarez](https://github.com/adriansuarez))

**Fixed bugs:**

- Changed hotcopy-{{ .Values.database.name }}-job-initial backup policy from Never to OnFailure to prevent accumulation of failed jobs [\#123](https://github.com/nuodb/nuodb-helm-charts/pull/123) ([vegichan](https://github.com/vegichan))
- Fixed Helm Operator Scorecard requirement of no empty blocks [\#104](https://github.com/nuodb/nuodb-helm-charts/pull/104) ([vegichan](https://github.com/vegichan))
- Removed cloud.sh autotagging of cloud environments to enable Helm Operator YCSB scaling [\#94](https://github.com/nuodb/nuodb-helm-charts/pull/94) ([vegichan](https://github.com/vegichan))

**Removed:**

- Removed monitoring-insights from Incubator [\#107](https://github.com/nuodb/nuodb-helm-charts/pull/107) ([vegichan](https://github.com/vegichan))
- Removed memoryOption and replaced it with resources.requests.memory [\#106](https://github.com/nuodb/nuodb-helm-charts/pull/106) ([vegichan](https://github.com/vegichan))
- Removed obsolete charts from Incubator \(backup, restore-copy, monitoring-influx\) [\#99](https://github.com/nuodb/nuodb-helm-charts/pull/99) ([vegichan](https://github.com/vegichan))
- Remove OpenShift specific TE DeploymentConfig [\#66](https://github.com/nuodb/nuodb-helm-charts/pull/66) ([vegichan](https://github.com/vegichan))

**Merged pull requests:**

- Bump NuoDB version to 4.0.6 [\#125](https://github.com/nuodb/nuodb-helm-charts/pull/125) ([vegichan](https://github.com/vegichan))
- Bump NuoDB version to 4.0.5 [\#92](https://github.com/nuodb/nuodb-helm-charts/pull/92) ([vegichan](https://github.com/vegichan))

## [v2.3.1](https://github.com/nuodb/nuodb-helm-charts/tree/v2.3.1) (2020-05-12)

[Full Changelog](https://github.com/nuodb/nuodb-helm-charts/compare/v2.3.0...v2.3.1)

**Fixed bugs:**

- Fixed an issue preventing Helm Chart release upgrade [\#111](https://github.com/nuodb/nuodb-helm-charts/issues/111)
- Removed NuodB Helm Chart version number from the release label. [\#110](https://github.com/nuodb/nuodb-helm-charts/pull/110) ([acabrele](https://github.com/acabrele))

## [v2.3.0](https://github.com/nuodb/nuodb-helm-charts/tree/v2.3.0) (2020-03-24)

[Full Changelog](https://github.com/nuodb/nuodb-helm-charts/compare/v2.2.0...v2.3.0)

**Implemented enhancements:**

- Add the "watch" verb to the NuoDB role to enable the Admin to register event listeners. [\#78](https://github.com/nuodb/nuodb-helm-charts/pull/78) ([adriansuarez](https://github.com/adriansuarez))
- Auto-configure role and service account in values files to enable NuoDB control plane synchronization with Kubernetes. [\#77](https://github.com/nuodb/nuodb-helm-charts/pull/77) ([adriansuarez](https://github.com/adriansuarez))
- Added the use of Security Context Constraints\(SCC\) for Red Hat OpenShift deployments. [\#75](https://github.com/nuodb/nuodb-helm-charts/pull/75) ([vegichan](https://github.com/vegichan))
- Increased engine readiness probe timeout value from 1 to 5 seconds. [\#74](https://github.com/nuodb/nuodb-helm-charts/pull/74) ([vegichan](https://github.com/vegichan))
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
- Add local-storage StorageClass to the StorageClass chart [\#28](https://github.com/nuodb/nuodb-helm-charts/pull/28) ([vegichan](https://github.com/vegichan))
- Add Cluster IP Services [\#27](https://github.com/nuodb/nuodb-helm-charts/pull/27) ([acabrele](https://github.com/acabrele))
- \[DB-29171\] Feature request for user provided environment variables for database chart [\#26](https://github.com/nuodb/nuodb-helm-charts/pull/26) ([sivanov-nuodb](https://github.com/sivanov-nuodb))
- Make Load Balancer services optional and add ability to provision internal cloud IPs [\#23](https://github.com/nuodb/nuodb-helm-charts/pull/23) ([acabrele](https://github.com/acabrele))

**Fixed bugs:**

- Small fixes for backup and restore [\#45](https://github.com/nuodb/nuodb-helm-charts/pull/45) ([acabrele](https://github.com/acabrele))
- Fix image pull policy reference [\#44](https://github.com/nuodb/nuodb-helm-charts/pull/44) ([acabrele](https://github.com/acabrele))
- Backup & restore, permissions and pull secrets fixes [\#42](https://github.com/nuodb/nuodb-helm-charts/pull/42) ([acabrele](https://github.com/acabrele))

**Merged pull requests:**

- Bump NuoDB version to 4.0.4 [\#47](https://github.com/nuodb/nuodb-helm-charts/pull/47) ([vegichan](https://github.com/vegichan))
- Bump NuoDB version to 4.0.3 [\#43](https://github.com/nuodb/nuodb-helm-charts/pull/43) ([vegichan](https://github.com/vegichan))
- Pin NuoDB version to 4.0.2 [\#35](https://github.com/nuodb/nuodb-helm-charts/pull/35) ([vegichan](https://github.com/vegichan))
- Downgrade to Helm v2.9.0 [\#25](https://github.com/nuodb/nuodb-helm-charts/pull/25) ([sivanov-nuodb](https://github.com/sivanov-nuodb))

## [v2.1](https://github.com/nuodb/nuodb-helm-charts/tree/v2.1) (2019-11-06)

[Full Changelog](https://github.com/nuodb/nuodb-helm-charts/compare/v2.0...v2.1)

**Implemented enhancements:**

- Add readiness probes for NuoDB engine StatefulSet/Deployment [\#18](https://github.com/nuodb/nuodb-helm-charts/pull/18) ([vegichan](https://github.com/vegichan))
- \[DB-28964\] Externalize config files using generalized idiomatic approaches. [\#13](https://github.com/nuodb/nuodb-helm-charts/pull/13) ([rbuck](https://github.com/rbuck))
- \[DB-28838\] Make nuoadmin configuration options available in helm values [\#12](https://github.com/nuodb/nuodb-helm-charts/pull/12) ([acabrele](https://github.com/acabrele))
- Add more state checks to Admin StatefulSet readiness probe [\#8](https://github.com/nuodb/nuodb-helm-charts/pull/8) ([vegichan](https://github.com/vegichan))
- \[DB-28313\] Enable passing of certificates directly to the engine [\#4](https://github.com/nuodb/nuodb-helm-charts/pull/4) ([vegichan](https://github.com/vegichan))

**Fixed bugs:**

- \[DB-28733\] add missing volumeMount to the THP chart [\#5](https://github.com/nuodb/nuodb-helm-charts/pull/5) ([vegichan](https://github.com/vegichan))
- \[DB-28712\] Fix Admin Resources [\#2](https://github.com/nuodb/nuodb-helm-charts/pull/2) ([vegichan](https://github.com/vegichan))

**Closed issues:**

- Error in admin helm chart? [\#1](https://github.com/nuodb/nuodb-helm-charts/issues/1)



\* *This Changelog was automatically generated by [github_changelog_generator](https://github.com/github-changelog-generator/github-changelog-generator)*
