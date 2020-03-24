# Changelog

## [2.3.0](https://github.com/nuodb/nuodb-helm-charts/tree/2.3.0) (2020-03-24)

[Full Changelog](https://github.com/nuodb/nuodb-helm-charts/compare/v2.2.0...2.3.0)

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
