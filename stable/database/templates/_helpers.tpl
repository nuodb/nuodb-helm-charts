{{/*
Expand the name of the chart.
*/}}
{{- define "database.name" -}}
{{- default .Chart.Name .Values.database.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 50 chars because some Kubernetes name fields are limited to 63 chars (by the DNS naming spec)
and we have to allow for added suffixes including "-hotcopy" and "-NN" where NN is the pod number.
*/}}
{{- define "database.fullname" -}}
{{- $domain := default "domain" ( include "admin.domainName" . ) -}}
{{- $cluster := default "cluster0" .Values.cloud.cluster.name -}}
{{- if .Values.database.fullnameOverride -}}
{{- .Values.database.fullnameOverride | trunc 50 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.database.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- printf "%s-%s-%s-%s" .Release.Name $domain $cluster ( include "database.dbName" . ) | trunc 50 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s-%s-%s-%s" .Release.Name $domain $cluster ( include "database.dbName" . ) $name | trunc 50 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
The StatefulSet name is limited to 52 charts (https://github.com/kubernetes/kubernetes/issues/64023).
*/}}
{{- define "database.statefulset.name" -}}
{{ template "truncWithHash" (list . 52) }}
{{- end -}}

{{/*
Create chart name as used by the chart label.
*/}}
{{- define "database.chart" -}}
{{- printf "%s" .Chart.Name | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Return the proper NuoDB image name
*/}}
{{- define "nuodb.image" -}}
{{- $registryName := .Values.nuodb.image.registry -}}
{{- $repositoryName := .Values.nuodb.image.repository -}}
{{- $tag := .Values.nuodb.image.tag | toString -}}
{{/*
Helm 2.11 supports the assignment of a value to a variable defined in a different scope,
but Helm 2.9 and 2.10 doesn't support it, so we need to implement this if-else logic.
Also, we can't use a single if because lazy evaluation is not an option
*/}}
{{- if .Values.global }}
    {{- if .Values.global.imageRegistry }}
        {{- printf "%s/%s:%s" .Values.global.imageRegistry $repositoryName $tag -}}
    {{- else -}}
        {{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
    {{- end -}}
{{- else -}}
    {{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
{{- end -}}
{{- end -}}

{{/*
Return the backup hooks sidecar image
*/}}
{{- define "backupHooks.image" -}}
{{- $registryName := .Values.database.backupHooks.image.registry -}}
{{- $repositoryName := .Values.database.backupHooks.image.repository -}}
{{- $tag := .Values.database.backupHooks.image.tag | toString -}}
{{/*
Helm 2.11 supports the assignment of a value to a variable defined in a different scope,
but Helm 2.9 and 2.10 doesn't support it, so we need to implement this if-else logic.
Also, we can't use a single if because lazy evaluation is not an option
*/}}
{{- if .Values.global }}
    {{- if .Values.global.imageRegistry }}
        {{- printf "%s/%s:%s" .Values.global.imageRegistry $repositoryName $tag -}}
    {{- else -}}
        {{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
    {{- end -}}
{{- else -}}
    {{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
{{- end -}}
{{- end -}}

{{/*
Renders the freeze mode for backup hooks.
*/}}
{{- define "backupHooks.freezeMode" -}}
{{- if eq .Values.database.backupHooks.freezeMode "" -}}
  {{- if (eq (include "defaultfalse" .Values.database.backupHooks.useSuspend) "true") -}}
suspend
  {{- else -}}
hotsnap
  {{- end -}}
{{- else -}}
{{ .Values.database.backupHooks.freezeMode }}
{{- end -}}
{{- end -}}

{{/*
Return init image
*/}}
{{- define "init.image" -}}
{{- $registryName :=  default "docker.io" .Values.busybox.image.registry -}}
{{- $repositoryName := .Values.busybox.image.repository -}}
{{- $tag := default "latest" .Values.busybox.image.tag | toString -}}
{{/*
Helm 2.11 supports the assignment of a value to a variable defined in a different scope,
but Helm 2.9 and 2.10 doesn't support it, so we need to implement this if-else logic.
Also, we can't use a single if because lazy evaluation is not an option
*/}}
{{- if .Values.global }}
    {{- if .Values.global.imageRegistry }}
        {{- printf "%s/%s:%s" .Values.global.imageRegistry $repositoryName $tag -}}
    {{- else -}}
        {{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
    {{- end -}}
{{- else -}}
    {{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
{{- end -}}
{{- end -}}

{{/*
Create a default fully qualified admin address.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "admin.address" -}}
{{- $domain := default "nuodb" ( include "admin.domainName" . ) -}}
{{- $namespace := default .Release.Namespace .Values.admin.namespace -}}
{{- printf "%s.%s.svc" $domain $namespace  | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified admin clusterip address.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "admin.clusterip" -}}
{{- $domain := default "nuodb" ( include "admin.domainName" . ) -}}
{{- $suffix := default "clusterip" .Values.admin.serviceSuffix.clusterip -}}
{{- $namespace := default .Release.Namespace .Values.admin.namespace -}}
{{- printf "%s-%s.%s.svc" $domain $suffix $namespace  | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Return the proper Docker Image Registry Secret Names
*/}}
{{- define "nuodb.imagePullSecrets" -}}
{{/*
Helm 2.11 supports the assignment of a value to a variable defined in a different scope,
but Helm 2.9 and 2.10 does not support it, so we need to implement this if-else logic.
Also, we can not use a single if because lazy evaluation is not an option
*/}}
{{- if .Values.global }}
{{- if .Values.global.imagePullSecrets }}
imagePullSecrets:
{{- range .Values.global.imagePullSecrets }}
  - name: {{ . }}
{{- end }}
{{- else if or .Values.busybox.image.pullSecrets .Values.nuodb.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.busybox.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- range .Values.nuodb.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- end -}}
{{- else if or .Values.busybox.image.pullSecrets .Values.nuodb.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.busybox.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- range .Values.nuodb.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- end -}}
{{- end -}}

{{/*
Get Pod securityContext (core/v1/PodSecurityContext)
*/}}
{{- define "database.securityContext" -}}
{{- if or (eq (include "defaultfalse" .Values.database.securityContext.enabled) "true") (eq (include "defaultfalse" .Values.database.securityContext.runAsNonRootGroup) "true") (eq (include "defaultfalse" .Values.database.securityContext.fsGroupOnly) "true") }}
securityContext:
  fsGroup: {{ default 1000 .Values.database.securityContext.fsGroup }}
  {{- include "sc.fsGroupChangePolicy" . | indent 2 }}
  {{- include "sc.runAs" . | indent 2 }}
{{- end }}
{{- end -}}

{{/*
Get security context runAsUser and runAsGroup
*/}}
{{- define "sc.runAs" -}}
{{- if eq (include "defaultfalse" .Values.database.securityContext.enabled) "true" }}
runAsUser: {{ default 1000 .Values.database.securityContext.runAsUser }}
runAsGroup: 0
{{- if ne (toString (default 1000 .Values.database.securityContext.runAsUser)) "0" }}
{{- include "sc.runAsNonRoot" . }}
{{- end }}
{{- else if eq (include "defaultfalse" .Values.database.securityContext.runAsNonRootGroup) "true" }}
runAsUser: 1000
runAsGroup: 1000
{{- include "sc.runAsNonRoot" . }}
{{- end }}
{{- end -}}

{{/*
Get security context runAsNonRoot
*/}}
{{- define "sc.runAsNonRoot" -}}
{{- $runAsNonRoot := true -}}
{{- if and
       (eq (include "defaulttrue" .Values.database.initContainers.runInitDisk) "true")
       (eq (include "defaulttrue" .Values.database.initContainers.runInitDiskAsRoot) "true") }}
  {{- $runAsNonRoot = false -}}
{{- else if .Values.database.backupHooks -}}
  {{- if and
         (eq (include "defaultfalse" .Values.database.backupHooks.enabled) "true")
         (eq (include "backupHooks.freezeMode" .) "fsfreeze")
         (eq (include "defaultfalse" .Values.database.sm.noHotCopy.journalPath.enabled) "true") -}}
    {{- $runAsNonRoot = false -}}
  {{- end -}}
{{- end -}}
{{- if $runAsNonRoot }}
runAsNonRoot: true
{{- end }}
{{- end -}}

{{/*
Get fsGroupChangePolicy if Kubernetes version supports it
*/}}
{{- define "sc.fsGroupChangePolicy" -}}
{{- if semverCompare ">=1.20" .Capabilities.KubeVersion.Version }}
fsGroupChangePolicy: OnRootMismatch
{{- end }}
{{- end -}}

{{/*
Get the Container securityContext (core/v1/SecurityContext)
*/}}
{{- define "database.sc.containerSecurityContext" }}
  {{- if eq (include "defaultfalse" .Values.database.securityContext.enabledOnContainer) "true" }}
securityContext:
  privileged: {{ include "defaultfalse" .Values.database.securityContext.privileged }}
  allowPrivilegeEscalation: {{ include "defaultfalse" .Values.database.securityContext.allowPrivilegeEscalation }}
  readOnlyRootFilesystem: {{ include "defaultfalse" .Values.database.securityContext.readOnlyRootFilesystem }}
  {{- include "sc.capabilities" . | indent 2 }}
  {{- include "sc.runAs" . | indent 2 }}
  {{- end }}
{{- end -}}

{{/*
Get container securityContext defining capabilities
*/}}
{{- define "sc.capabilities" -}}
  {{- if .Values.database.securityContext.capabilities }}
    {{- if typeIs "[]interface {}" .Values.database.securityContext.capabilities }}
capabilities:
      {{- with .Values.database.securityContext.capabilities }}
  add: {{ . }}
      {{- end }}
    {{- else if or .Values.database.securityContext.capabilities.add .Values.database.securityContext.capabilities.drop }}
capabilities:
      {{- if .Values.database.securityContext.capabilities.add }}
  add:
        {{- toYaml .Values.database.securityContext.capabilities.add | trim | nindent 4 }}
      {{- end }}
      {{- if .Values.database.securityContext.capabilities.drop }}
  drop:
        {{- toYaml .Values.database.securityContext.capabilities.drop | trim | nindent 4 }}
      {{- end }}
    {{- end }}
  {{- end }}
{{- end -}}

{{/*
Import ENV vars from configMaps
**BEWARE!!**
   The values for envFrom are formated into a single line because some parsers
   - either in k8s or rancher - throw errors occasionally if the multi-line format is used.
   You Have Been Warned.
*/}}
{{- define "database.envFrom" }}
envFrom: [ configMapRef: { name: {{ template "database.fullname" . }}-restore } {{- range $map := .Values.database.envFrom.configMapRef }}, configMapRef: { name: {{$map}} } {{- end }} ]
{{- end -}}

{{/*
Return options as $key $value
*/}}
{{- define "opt.key-values" -}}
{{- range $opt, $val := . }} {{$opt}} {{$val}} {{- end}}
{{- end -}}

{{/*
Return the API version of the CronJob kind
*/}}
{{- define "cronjob.apiVersion" -}}
{{- if .Capabilities.APIVersions.Has "batch/v1/CronJob" -}}
batch/v1
{{- else -}}
batch/v1beta1
{{- end -}}
{{- end -}}

{{/*
Return the hotcopy group prefix
*/}}
{{- define "hotcopy.groupPrefix" -}}
{{ default .Values.cloud.cluster.name .Values.database.sm.hotCopy.backupGroupPrefix }}
{{- end -}}

{{/*
Return the hotcopy cronjob schedule by backup group and hot copy type. It will take
into account any schedule overrides configured per backup group.
*/}}
{{- define "hotcopy.group.schedule" -}}
  {{- $scheduleProp := printf "%sSchedule" .hotCopyType -}}
  {{- $schedule := index .Values.database.sm.hotCopy ( print $scheduleProp ) -}}
  {{- $override := "" -}}
  {{- if eq .hotCopyType "journal" -}}
    {{- $schedule = .Values.database.sm.hotCopy.journalBackup.journalSchedule -}}
  {{- end -}}
  {{- if .Values.database.sm.hotCopy.backupGroups -}}
    {{- $group := (index .Values.database.sm.hotCopy.backupGroups ( print .backupGroup )) -}}
    {{- if $group -}}
      {{- $override = index $group ( print $scheduleProp ) -}}
    {{- end -}}
  {{- end -}}
  {{- default $schedule $override }}
{{- end -}}

{{/*
Renders the name of the HotCopy CronJob
*/}}
{{- define "hotcopy.cronjob.name" -}}
{{- $name := printf "%s-hotcopy-%s-%s-%s" .hotCopyType ( include "admin.domainName" . ) ( include "database.dbName" . ) .backupGroup -}}
{{ template "truncWithHash" (list $name 52) }}
{{- end -}}

{{/*
Return labels for a specific backup group. If there is no user-defined backup
groups, return the pod name of a single HCSM by extracting the pod ordinal from
the automatically generated backup group name. Otherwise return the configured
backup group labels or empty value (representing all HCSMs in the database).
*/}}
{{- define "hotcopy.group.labels" -}}
  {{- if .Values.database.sm.hotCopy.backupGroups -}}
    {{- $group := index .Values.database.sm.hotCopy.backupGroups (print .backupGroup) -}}
    {{- if $group -}}
    {{- default "" $group.labels -}}
    {{- end -}}
  {{- else -}}
    {{- $groupPrefix := include "hotcopy.groupPrefix" . -}}
    {{- printf "pod-name sm-%s-hotcopy%s" (include "database.fullname" .) (trimPrefix $groupPrefix .backupGroup) }}
  {{- end -}}
{{- end -}}

{{/*
Return process filter for a specific backup if one is defined.
*/}}
{{- define "hotcopy.group.processFilter" -}}
  {{- if .Values.database.sm.hotCopy.backupGroups -}}
    {{- $group := index .Values.database.sm.hotCopy.backupGroups (print .backupGroup) -}}
    {{- if $group -}}
    {{- default "" $group.processFilter -}}
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Import user defined ENV vars
*/}}
{{- define "database.env" }}
{{- if not (empty .Values.database.env) }}
{{ toYaml .Values.database.env | trim }}
{{- end }}
{{- end -}}

{{/*
Validates parameter that supports bool value only
*/}}
{{- define "validate.boolString" -}}
{{- $valid := list "true" "false" "" nil }}
{{- if not (has . $valid) }}
{{- fail (printf "Invalid boolean value: %s" .) }}
{{- end }}
{{- end -}}

{{/*
Takes a boolean as argument return it's value if it was defined or return true otherwise
Note: Sprig's default function on an empty/not defined variable returns false, workaround
it by calling typeIs "bool" https://github.com/Masterminds/sprig/issues/111
*/}}
{{- define "defaulttrue" -}}
{{- if typeIs "bool" . -}}
{{- . -}}
{{- else -}}
{{- template "validate.boolString" . -}}
{{- default true . -}}
{{- end -}}
{{- end -}}

{{/*
Takes a boolean as argument return it's value if it was defined or return false otherwise
Note: Sprig's default function on an empty/not defined variable returns false, workaround
it by calling typeIs "bool" https://github.com/Masterminds/sprig/issues/111
*/}}
{{- define "defaultfalse" -}}
{{- if typeIs "bool" . -}}
{{- . -}}
{{- else -}}
{{- template "validate.boolString" . -}}
{{- default false . -}}
{{- end -}}
{{- end -}}

{{/*
Imports per-database load balancer configuration via annotations.
The configuration is imported only in the entrypoint cluster.
*/}}
{{- define "database.loadBalancerConfig" -}}
{{- if .Values.database.lbConfig }}
{{- if (eq (default "cluster0" .Values.cloud.cluster.name) (default "cluster0" .Values.cloud.cluster.entrypointName)) }}
{{- if eq (include "defaulttrue" .Values.database.primaryRelease) "true" }}
{{- with .Values.database.lbConfig.prefilter }}
"nuodb.com/load-balancer-prefilter": {{ . | quote }}
{{- end -}}
{{- with .Values.database.lbConfig.default }}
"nuodb.com/load-balancer-default": {{ . | quote }}
{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Enables automatic database protocol upgrade via annotations.
The configuration is imported only in the entrypoint cluster.
*/}}
{{- define "database.automaticProtocolUpgrade" -}}
  {{- if .Values.database.automaticProtocolUpgrade }}
    {{- if eq (include "defaultfalse" .Values.database.automaticProtocolUpgrade.enabled) "true" -}}
      {{- if eq (include "defaulttrue" .Values.database.primaryRelease) "true" }}
        {{- if (eq (default "cluster0" .Values.cloud.cluster.name) (default "cluster0" .Values.cloud.cluster.entrypointName)) }}
"nuodb.com/automatic-database-protocol-upgrade": "true"
          {{- with .Values.database.automaticProtocolUpgrade.tePreferenceQuery }}
"nuodb.com/automatic-database-protocol-upgrade.te-preference-query": {{ . | quote }}
          {{- end -}}
        {{- end -}}
      {{- end -}}
    {{- end -}}
  {{- end -}}
{{- end -}}

{{- define "autoRestore.type" -}}
{{- $valid := list "stream" "backupset" "" }}
{{- if not (has .Values.database.autoRestore.type $valid) }}
{{- fail (printf "Invalid autorestore type: %s" .Values.database.autoRestore.type) }}
{{- end }}
{{- with .Values.database.autoRestore.type }}
{{- . }}
{{- end -}}
{{- end -}}

{{- define "autoRestore.source" -}}
{{- if hasPrefix ":" .Values.database.autoRestore.source }}
{{- $valid := list ":latest" ":group-latest" "" }}
{{- if not (has .Values.database.autoRestore.source $valid) }}
{{- fail (printf "Invalid autorestore source: %s" .Values.database.autoRestore.source) }}
{{- end }}
{{- end -}}
{{- with .Values.database.autoRestore.source }}
{{- . }}
{{- end -}}
{{- end -}}

{{/*
Renders nuodocker options and flags. An option is rendered only if its value is
not empty. Flags can be defined by setting their value to boolean true or "true"
*/}}
{{- define "database.otherOptions" -}}
  {{- range $opt, $val := . }}
    {{- if $val }}
      {{- if (ne (toString $val) "false") }}
- "--{{$opt}}"
        {{ if ne (toString $val) "true" }}
- "{{$val}}"
        {{- end }}
      {{- end }}
    {{- end }}
  {{- end }}
{{- end}}

{{/*
Renders the database service name for external access based on the service type
*/}}
{{- define "database.externalServiceName" -}}
  {{- $serviceType := (default "LoadBalancer" .Values.database.te.externalAccess.type) -}}
  {{- if eq $serviceType "LoadBalancer" -}}
{{ template "database.fullname" . }}-{{ default .Values.admin.serviceSuffix.balancer .Values.database.serviceSuffix.balancer }}
  {{- else if eq $serviceType "NodePort" -}}
{{ template "database.fullname" . }}-{{ default .Values.admin.serviceSuffix.nodeport .Values.database.serviceSuffix.nodeport }}
  {{- else -}}
{{ template "database.fullname" . }}
  {{- end }}
{{- end }}

{{/*
Renders the annotations for the LoadBalancer database service
*/}}
{{- define "database.externalAccessAnnotations" -}}
  {{- if eq (default "LoadBalancer" .Values.database.te.externalAccess.type) "LoadBalancer" }}
    {{- if .Values.database.te.externalAccess.annotations }}
{{ toYaml .Values.database.te.externalAccess.annotations | trim }}
    {{- else -}}
      {{- if .Values.cloud.provider }}
        {{- if eq .Values.cloud.provider "amazon" }}
          {{- if .Values.database.te.externalAccess.internalIP }}
service.beta.kubernetes.io/aws-load-balancer-internal: "true"
service.beta.kubernetes.io/aws-load-balancer-scheme: "internal"
          {{- else }}
service.beta.kubernetes.io/aws-load-balancer-type: "external"
service.beta.kubernetes.io/aws-load-balancer-nlb-target-type: "ip"
service.beta.kubernetes.io/aws-load-balancer-scheme: "internet-facing"
          {{- end }}
        {{- else if eq .Values.cloud.provider "azure" }}
          {{- if .Values.database.te.externalAccess.internalIP }}
service.beta.kubernetes.io/azure-load-balancer-internal: "true"
          {{- end }}
        {{- else if eq .Values.cloud.provider "google" }}
          {{- if .Values.database.te.externalAccess.internalIP }}
cloud.google.com/load-balancer-type: "Internal"
networking.gke.io/load-balancer-type: "Internal"
          {{- end -}}
        {{- end -}}
      {{- end -}}
    {{- end }}
  {{- end }}
{{- end -}}

{{/*
Renders the Transaction engine labels and injects external address and port
if Ingress is enabled
*/}}
{{- define "database.teLabels" -}}
  {{- $extraLabels := ""  }}
  {{- if .Values.database.te.ingress }}
    {{- if eq (include "defaultfalse" .Values.database.te.ingress.enabled) "true" }}
      {{- if .Values.database.te.ingress.hostname }}
        {{- if not (index .Values.database.te.labels "external-address") }}
          {{- $extraLabels = printf "external-address %s" .Values.database.te.ingress.hostname }}
        {{- end }}
        {{- if not (index .Values.database.te.labels "external-port") }}
          {{- $extraLabels = printf "%s external-port 443" $extraLabels }}
        {{- end }}
      {{- end }}
    {{- end }}
  {{- end }}
  {{- if or .Values.database.te.labels $extraLabels }}
- "--labels"
- "{{ $extraLabels }}{{ include "opt.key-values" .Values.database.te.labels }}"
  {{- end }}
{{- end -}}

{{/*
Renders the labels for all resources deployed by this Helm chart
*/}}
{{- define "database.resourceLabels" -}}
{{- include "database.labels" (dict "Root" . "ExtraLabels" .Values.database.resourceLabels ) -}}
{{- end -}}

{{- define "database.labels" -}}
{{- if not .Root}}{{fail "<nuodb.env> Root is required"}}{{end}}
{{- $extras := .ExtraLabels | default (dict) -}}
{{- include "database.selectorLabels" .Root }}
group: nuodb
database: {{ include "database.dbName" .Root }}
domain: {{ include "admin.domainName" .Root }}
chart: {{ template "database.chart" .Root }}
release: {{ .Root.Release.Name | quote }}
{{- range $k, $v := $extras }}
"{{ $k }}": "{{ $v }}"
{{- end }}
{{- end -}}

{{/*
Selector labels
*/}}
{{- define "database.selectorLabels" -}}
app: {{ template "database.fullname" . }}
{{- end }}

{{/*
Renders the labels for SM pods deployed by this Helm chart
*/}}
{{- define "database.sm.podLabels" -}}
{{- include "database.resourceLabels" .}}
{{- end -}}

{{/*
Renders the labels for TE pods deployed by this Helm chart
*/}}
{{- define "database.te.podLabels" -}}
{{- include "database.resourceLabels" .}}
{{- end -}}

{{/*
Renders the labels for SM pod volumes deployed by this Helm chart
*/}}
{{- define "database.sm.volumeLabels" -}}
{{- include "database.resourceLabels" .}}
{{- end -}}

{{/*
Renders the labels for TE pod volumes deployed by this Helm chart
*/}}
{{- define "database.te.volumeLabels" -}}
{{- include "database.resourceLabels" .}}
{{- end -}}

{{/*
Renders the storage group labels. IMPORTANT: Adding new entries must be done
with caution as these labels are rendered in the PVC spec which is immutable.
The storage group name can not be reconfigured by definition.
*/}}
{{- define "database.storageGroupLabels" -}}
{{- if .Values.database.sm.storageGroup -}}
{{- if .Values.database.sm.storageGroup.enabled -}}
storage-group: {{ include "database.storageGroup.name" . | quote }}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Renders the name of the Secret for this database
*/}}
{{- define "database.secretName" -}}
{{ include "admin.domainName" . }}-{{ include "database.dbName" . }}
{{- end -}}

{{/*
Renders the storage group name
*/}}
{{- define "database.storageGroup.name" -}}
{{ .Values.database.sm.storageGroup.name | default .Release.Name | trim }}
{{- end -}}

{{/*
Renders the storage groups option for nuosm script
*/}}
{{- define "database.storageGroup.args" -}}
{{- if .Values.database.sm.storageGroup }}
{{- if .Values.database.sm.storageGroup.enabled }}
{{- with .Values.database.sm.storageGroup.name }}
{{- $invalid := list "ALL" "UNPARTITIONED" }}
{{- if has (upper .) $invalid }}
{{- fail (printf "Invalid storage group name: %s" .) }}
{{- end }}
{{- if contains " " (trim .) }}
{{- fail (printf "Multiple storage group names provided: %s" .) }}
{{- end }}
{{- end }}
- "--storage-groups"
- {{ include "database.storageGroup.name" . | quote }}
{{- end }}
{{- end }}
{{- end -}}

{{/*
Renders the storage group domain process label
*/}}
{{- define "database.storageGroup.label" -}}
{{- if .Values.database.sm.storageGroup -}}
{{- if .Values.database.sm.storageGroup.enabled -}}
{{ printf "sg %s" (include "database.storageGroup.name" .) }}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Checks if the TE Deployment should be enabled. If the default value is absent,
the TE Deployment is disabled if TPSG is enabled and this is a secondary release.
*/}}
{{- define "database.te.enablePod" -}}
  {{- if kindIs "invalid" .Values.database.te.enablePod -}}
    {{- if and (eq (include "defaultfalse" .Values.database.sm.storageGroup.enabled) "true") (eq (include "defaulttrue" .Values.database.primaryRelease) "false") -}}
      {{- false -}}
    {{- else -}}
      {{- true -}}
    {{- end -}}
  {{- else -}}
    {{- include "defaulttrue" .Values.database.te.enablePod -}}
  {{- end -}}
{{- end -}}

{{/*
Renders an ephemeral volume for an engine process.
*/}}
{{- define "database.ephemeralVolume" -}}
{{- $ := index . 0 -}}
{{- $engine := index . 1 -}}
{{- $engine_type := index . 2 -}}
{{- if eq (include "defaultfalse" $.Values.database.ephemeralVolume.enabled) "true" }}
ephemeral:
  volumeClaimTemplate:
    metadata:
      labels:
        {{- if (eq "te" $engine_type) }}
          {{- include "database.te.volumeLabels" $ | nindent 10 }}
        {{- else }}
          {{- include "database.sm.volumeLabels" $ | nindent 10 }}
        {{- end }}
    spec:
      accessModes:
      - ReadWriteOnce
      {{- if $.Values.database.ephemeralVolume.storageClass }}
      {{- if (eq "-" $.Values.database.ephemeralVolume.storageClass) }}
      storageClassName: ""
      {{- else }}
      storageClassName: {{ $.Values.database.ephemeralVolume.storageClass }}
      {{- end }}
      {{- end }}
      resources:
        requests:
          {{- if eq (include "defaultfalse" $.Values.database.ephemeralVolume.sizeToMemory) "true" }}
          storage: {{ $engine.resources.limits.memory }}
          {{- else }}
          storage: {{ $.Values.database.ephemeralVolume.size }}
          {{- end }}
{{- else }}
emptyDir: {}
{{- end }}
{{- end -}}

{{/*
Returns true if either a generic ephemeral volume or emptyDir volume should be
rendered, false otherwise.
*/}}
{{- define "database.enableEphemeralVolume" -}}
{{- $ := index . 0 -}}
{{- $engine := index . 1 -}}
{{- $ret := false -}}
{{- if eq (include "defaultfalse" $engine.logPersistence.enabled) "false" -}}
  {{- $ret = true -}}
{{- end -}}
{{- if eq (include "defaultfalse" $.Values.database.securityContext.enabledOnContainer) "true" -}}
{{- if eq (include "defaultfalse" $.Values.database.securityContext.readOnlyRootFilesystem) "true" -}}
  {{- $ret = true -}}
{{- end -}}
{{- end -}}
{{- if $.Values.nuocollector }}
{{- if eq (include "defaultfalse" $.Values.nuocollector.enabled) "true" }}
  {{- $ret = true -}}
{{- end -}}
{{- end -}}
{{ $ret }}
{{- end -}}

{{/*
Truncates a string to the max length of characters. If the original string
length exceeds the supplied max length, a 7 characters SHA256 hash is
calculated and appended to the truncated string to make it unique.
*/}}
{{- define "truncWithHash" -}}
{{- $s := index . 0 -}}
{{- $max := index . 1 -}}
{{- if not (typeIs "string" $s) -}}
{{- fail (printf "truncWithHash template failed: not supported type '%s'" (typeOf $s)) -}}
{{- end -}}
{{- if not (typeIs "int" $max) -}}
{{- fail (printf "truncWithHash template failed: not supported type '%s' for max length" (typeOf $max)) -}}
{{- end -}}
{{- if lt $max 15 -}}
{{- fail (printf "truncWithHash template failed: max string length %d is too small; must be at lest 15 chars" $max) -}}
{{- end -}}
{{- if gt (len $s) $max -}}
{{- $hash := sha256sum $s | trunc 7 -}}
{{- $truncated := $s | trunc (int (sub $max 8)) | trimSuffix "-" -}}
{{ printf "%s-%s" $truncated $hash }}
{{- else -}}
{{ $s }}
{{- end -}}
{{- end -}}

{{/*
Renders the TLS secrets as projected volume. Combining all secrets into one
volume allows mounting them into single directory (/etc/nuodb/keys) without
"subPath".
*/}}
{{- define "database.tlsVolume" -}}
{{- if .Values.admin.tlsClientPEM }}
- name: tls
  projected:
    defaultMode: 0440
    sources:
{{- if .Values.admin.tlsCACert }}
    - secret:
        name: {{ .Values.admin.tlsCACert.secret }}
        items:
        - key: {{ .Values.admin.tlsCACert.key }}
          path: ca.cert
{{- end }}
{{- if .Values.admin.tlsKeyStore }}
    - secret:
        name: {{ .Values.admin.tlsKeyStore.secret }}
        items:
        - key: {{ .Values.admin.tlsKeyStore.key }}
          path: nuoadmin.p12
{{- end }}
{{- if .Values.admin.tlsClientPEM }}
    - secret:
        name: {{ .Values.admin.tlsClientPEM.secret }}
        items:
        - key: {{ .Values.admin.tlsClientPEM.key }}
          path: nuocmd.pem
{{- end }}
{{- end }}
{{- end -}}

{{/*
Renders the TLS password for keystore or truststore.
*/}}
{{- define "database.keystorePassword" -}}
{{- $ := index . 0 -}}
{{- $store := index . 1 -}}
{{- if $store.password -}}
{{ $store.password }}
{{- else -}}
  {{- $secret := lookup "v1" "Secret" $.Release.Namespace $store.secret -}}
  {{- if $secret -}}
  {{- $encoded := index $secret "data" (default "password" $store.passwordKey) -}}
    {{- if $encoded -}}
{{ print $encoded | b64dec }}
    {{- end -}}
  {{- end -}}
{{- end -}}
{{- end -}}

{{/*
Renders the TLS keystore or truststore.
*/}}
{{- define "database.keystore" -}}
{{- $ := index . 0 -}}
{{- $store := index . 1 -}}
{{- $secret := lookup "v1" "Secret" $.Release.Namespace $store.secret -}}
{{- if $secret -}}
  {{- $encoded := index $secret "data" $store.key -}}
  {{- if $encoded -}}
{{ print $encoded | b64dec }}
  {{- end -}}
{{- end -}}
{{- end -}}

{{/*
Renders TE/SM Pod annotations for TLS checksum that forces engine pods restart on
configuration change. Pod restart is needed if the keystore or its password
have changed. If the AP server certificate is rotated, new process certificates
will be generated by NuoAdmin on engine startup.
*/}}
{{- define "database.tlsConfigAnnotations" -}}
{{- if .Values.admin.tlsKeyStore }}
  {{- $password := include "database.keystorePassword" (list . .Values.admin.tlsKeyStore) -}}
  {{- if $password }}
checksum/tls-passwords: {{ sha256sum $password }}
  {{- end }}
  {{- $keystore := include "database.keystore" (list . .Values.admin.tlsKeyStore) -}}
  {{- if $keystore }}
checksum/tls-keystore: {{ sha256sum $keystore }}
  {{- end }}
{{- end }}
{{- end -}}

{{/*
Set shareProcessNamespace=true if needed for an SM statefulset.
*/}}
{{- define "database.sm.shareProcessNamespace" -}}
{{- $hasSidecars := false -}}
{{- if .Values.nuocollector -}}
  {{- if eq (include "defaultfalse" .Values.nuocollector.enabled) "true" -}}
    {{- $hasSidecars = true -}}
  {{- end -}}
{{- end -}}
{{- if .Values.database.backupHooks -}}
  {{- if eq (include "defaultfalse" .Values.database.backupHooks.enabled) "true" -}}
    {{- $hasSidecars = true -}}
  {{- end -}}
{{- end -}}
{{- if $hasSidecars }}
shareProcessNamespace: true
{{- end }}
{{- end -}}

{{/*
Set shareProcessNamespace=true if needed for a TE deployment.
*/}}
{{- define "database.te.shareProcessNamespace" -}}
{{- $hasSidecars := false -}}
{{- if .Values.nuocollector -}}
  {{- if eq (include "defaultfalse" .Values.nuocollector.enabled) "true" -}}
    {{- $hasSidecars = true -}}
  {{- end -}}
{{- end -}}
{{- if $hasSidecars }}
shareProcessNamespace: true
{{- end }}
{{- end -}}

{{/*
Render dataSourceRef if dataSourceRef.name is not empty.
*/}}
{{- define "database.dataSource" -}}
{{- if . -}}
  {{- if .name -}}
dataSourceRef:
  {{- toYaml . | nindent 2 }}
  {{- end -}}
{{- end -}}
{{- end -}}

{{/*
Render dataSourceRef from snapshotRestore configuration.
*/}}
{{- define "database.snapshotRestoreDataSource" -}}
{{- $ := index . 0 -}}
{{- $volumeType := index . 1 -}}
{{- if $.Values.database.snapshotRestore -}}
  {{- if $.Values.database.snapshotRestore.backupId -}}
    {{- $context := dict "backupId" $.Values.database.snapshotRestore.backupId "volumeType" $volumeType "Template" $.Template -}}
dataSourceRef:
  name: {{ tpl $.Values.database.snapshotRestore.snapshotNameTemplate $context | quote }}
  kind: VolumeSnapshot
  apiGroup: snapshot.storage.k8s.io
  {{- end -}}
{{- end -}}
{{- end -}}

{{/*
Validate and render dataSourceRef.
*/}}
{{- define "database.validateAndRenderDataSource" -}}
{{- $ := index . 0 -}}
{{- $dataSource := index . 1 -}}
{{- if $dataSource -}}
  {{- if eq (include "defaulttrue" $.Values.database.persistence.validateDataSources) "true" -}}
    {{- $ref := (fromYaml $dataSource).dataSourceRef -}}
    {{- $apiVersion := "v1" -}}
    {{- if $ref.apiGroup -}}
      {{- $apiVersion = printf "%s/v1" $ref.apiGroup -}}
    {{- end -}}
    {{- $namespace := default $.Release.Namespace $ref.namespace -}}
    {{- if not (lookup $apiVersion $ref.kind $namespace $ref.name) -}}
      {{- if and $.Release.IsUpgrade (eq (include "defaultfalse" $.Values.database.persistence.preprovisionVolumes) "true") -}}
        {{- $dataSource = "" -}}
      {{- else -}}
        {{- fail (printf "Invalid data source: %s/%s/%s not found in namespace %s" $apiVersion $ref.kind $ref.name $namespace) -}}
      {{- end -}}
    {{- end -}}
  {{- end -}}
  {{- print $dataSource -}}
{{- end -}}
{{- end -}}

{{/*
Render archive dataSourceRef, giving precedence to
database.persistent.archiveDataSource and falling back to snapshot name
resolved from snapshotRestore configuration.
*/}}
{{- define "database.archiveDataSource" -}}
{{- $dataSource := include "database.dataSource" .Values.database.persistence.archiveDataSource -}}
{{- if not $dataSource -}}
  {{- $dataSource = include "database.snapshotRestoreDataSource" (list . "archive") -}}
{{- end -}}
{{- include "database.validateAndRenderDataSource" (list . $dataSource) -}}
{{- end -}}

{{/*
Render journal dataSourceRef, giving precedence to
database.persistent.journalDataSource and falling back to snapshot name
resolved from snapshotRestore configuration.
*/}}
{{- define "database.journalDataSource" -}}
{{- $dataSource := include "database.dataSource" .Values.database.persistence.journalDataSource -}}
{{- if not $dataSource -}}
  {{- $dataSource = include "database.snapshotRestoreDataSource" (list . "journal") -}}
{{- end -}}
{{- include "database.validateAndRenderDataSource" (list . $dataSource) -}}
{{- end -}}

{{/*
Get SNAPSHOT_RESTORED environment variable value.
*/}}
{{- define "database.snapshotRestored" -}}
{{- $dataSource := include "database.archiveDataSource" . -}}
{{- if $dataSource -}}
true
{{- else -}}
false
{{- end -}}
{{- end -}}

{{/*
Get spec of archive PVC or volumeClaimTemplate of SM statefulset.
*/}}
{{- define "database.archivePvcSpec" -}}
{{- $ := index . 0 -}}
{{- $includeDataSource := index . 1 -}}
accessModes:
{{- range $.Values.database.persistence.accessModes }}
  - {{ . }}
{{- end }}
{{- if $.Values.database.persistence.storageClass }}
{{- if eq "-" $.Values.database.persistence.storageClass }}
storageClassName: ""
{{- else }}
storageClassName: {{ $.Values.database.persistence.storageClass }}
{{- end }}
{{- end }}
{{- if eq (include "defaultfalse" $includeDataSource) "true" }}
{{ include "database.archiveDataSource" $ }}
{{- end }}
{{- if $.Values.database.isManualVolumeProvisioning }}
selector:
  matchLabels:
    database: {{ include "database.dbName" . }}
{{- end }}
resources:
  requests:
    storage: {{ $.Values.database.persistence.size }}
{{- end -}}

{{/*
Get spec of journal PVC or volumeClaimTemplate of SM statefulset.
*/}}
{{- define "database.journalPvcSpec" -}}
{{- $ := index . 0 -}}
{{- $includeDataSource := index . 1 -}}
accessModes:
{{- range $.Values.database.sm.noHotCopy.journalPath.persistence.accessModes }}
  - {{ . }}
{{- end }}
{{- if $.Values.database.sm.noHotCopy.journalPath.persistence.storageClass }}
{{- if eq "-" $.Values.database.sm.noHotCopy.journalPath.persistence.storageClass }}
storageClassName: ""
{{- else }}
storageClassName: {{ $.Values.database.sm.noHotCopy.journalPath.persistence.storageClass }}
{{- end }}
{{- end }}
{{- if eq (include "defaultfalse" $includeDataSource) "true" }}
{{ include "database.journalDataSource" $ }}
{{- end }}
{{- if $.Values.database.isManualVolumeProvisioning }}
selector:
  matchLabels:
    database: {{ include "database.dbName" . }}
{{- end }}
resources:
  requests:
    storage: {{ $.Values.database.sm.noHotCopy.journalPath.persistence.size }}
{{- end -}}

{{/*
Any additional fields that need to go into the database container spec.
*/}}
{{- define "database.podSpecExtras"}}
{{/*
Extension point that can be overriden by an embedding chart.
*/}}
{{- end }}

{{/*
Name of the domain to connect to.
*/}}
{{- define "admin.domainName" -}}
{{- .Values.admin.domain -}}
{{- end -}}

{{/*
Database name
*/}}
{{- define "database.dbName" -}}
{{- .Values.database.name -}}
{{- end -}}

{{/*
Returns true if TE replicas must be set on the Deployment. When TE autoscaling
is enabled, it is recommended that the `spec.replicas` field is removed. The HPA
will be implicitly deactivated if the desired replicas is set to 0 (zero). See
https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#migrating-deployments-and-statefulsets-to-horizontal-autoscaling
*/}}
{{- define "database.te.setReplicas" -}}
{{- if or (eq (include "database.te.replicas" .) "0") (and (eq (include "database.hpa.enabled" .) "false") (eq (include "database.keda.enabled" .) "false")) -}}
true
{{- else -}}
false
{{- end -}}
{{- end -}}

{{/*
Number of TE replicas to deploy.
*/}}
{{- define "database.te.replicas" -}}
{{ .Values.database.te.replicas }}
{{- end -}}

{{/*
Number of non-hotcopy SM replicas to deploy.
*/}}
{{- define "database.sm.noHotCopy.replicas" -}}
{{ .Values.database.sm.noHotCopy.replicas }}
{{- end -}}

{{/*
Number of hotcopy SM replicas to deploy.
*/}}
{{- define "database.sm.hotCopy.replicas" -}}
{{ .Values.database.sm.hotCopy.replicas }}
{{- end -}}

{{/*
SM resources requested and limited
*/}}
{{- define "database.sm.resources" -}}
{{- toYaml .Values.database.sm.resources  -}}
{{- end -}}

{{/*
TE resources requested and limited
*/}}
{{- define "database.te.resources" -}}
{{- toYaml .Values.database.te.resources  -}}
{{- end -}}

{{/*
Any additional volumes that need to go into the SM pod spec.
*/}}
{{- define "database.sm.extraVolumes"}}
{{/*
Extension point that can be overriden by an embedding chart.
*/}}
{{- end }}

{{/*
Any additional volume mounts that need to go into the SM container.
*/}}
{{- define "database.sm.extraMounts"}}
{{/*
Extension point that can be overriden by an embedding chart.
*/}}
{{- end }}

{{/*
Any additional sidecar containers that need to go into the SM pod.
*/}}
{{- define "database.sm.extraSidecars"}}
{{/*
Extension point that can be overriden by an embedding chart.
*/}}
{{- end }}

{{/*
Backup hooks sidecar container resources requested and limited
*/}}
{{- define "database.backupHooks.resources" -}}
{{- toYaml .Values.database.backupHooks.resources  -}}
{{- end -}}

{{/*
Command name to start the SM engine.
*/}}
{{- define "database.sm.entryPoint" -}}
nuosm
{{- end -}}

{{/*
Import user defined ENV vars for the backup hooks sidecar
*/}}
{{- define "database.backupHooks.env" }}
{{- if not (empty .Values.database.backupHooks.env) }}
{{ toYaml .Values.database.backupHooks.env | trim }}
{{- end }}
{{- end -}}

{{/*
Returns true if HPA resource is enabled
*/}}
{{- define "database.hpa.enabled" -}}
{{- if eq (include "defaultfalse" .Values.database.te.autoscaling.hpa.enabled) "true" -}}
true
{{- else -}}
false
{{- end -}}
{{- end -}}

{{/*
Returns true if KEDA ScaledObject resource is enabled
*/}}
{{- define "database.keda.enabled" -}}
{{- if eq (include "defaultfalse" .Values.database.te.autoscaling.keda.enabled) "true" -}}
{{- if eq (include "defaultfalse" .Values.database.te.autoscaling.hpa.enabled) "true" -}}
{{- fail "Can not enable both HPA and KEDA for TE autoscaling" }}
{{- end -}}
true
{{- else -}}
false
{{- end -}}
{{- end -}}

{{/*
Return the API version of the HorizontalPodAutoscaler kind
*/}}
{{- define "database.hpa.apiVersion" -}}
{{- if .Capabilities.APIVersions.Has "autoscaling/v2/HorizontalPodAutoscaler" -}}
autoscaling/v2
{{- else -}}
autoscaling/v2beta1
{{- end -}}
{{- end -}}

{{/*
Return the minReplicas setting for HPA
*/}}
{{- define "database.hpa.minReplicas" -}}
{{ .Values.database.te.autoscaling.minReplicas }}
{{- end -}}

{{/*
Return the maxReplicas setting for HPA
*/}}
{{- define "database.hpa.maxReplicas" -}}
{{ .Values.database.te.autoscaling.maxReplicas }}
{{- end -}}

{{/*
Return the behaviors setting for HPA
*/}}
{{- define "database.hpa.behavior" -}}
{{ toYaml .Values.database.te.autoscaling.hpa.behavior }}
{{- end -}}

{{/*
Return the targetCpuUtilization setting
*/}}
{{- define "database.targetCpuUtilization" -}}
{{ .Values.database.te.autoscaling.hpa.targetCpuUtilization }}
{{- end -}}

{{/*
Return the KEDA triggers
*/}}
{{- define "database.keda.triggers" -}}
{{- $triggers := toYaml .Values.database.te.autoscaling.keda.triggers -}}
{{- tpl $triggers . }}
{{- end -}}

{{/*
Any additional annotations to add to services.
*/}}
{{- define "database.serviceAnnotations"}}
{{/*
Extension point that can be overriden by an embedding chart.
*/}}
{{- end }}

{{/*
Database root user
*/}}
{{- define "database.rootUser" -}}
{{ .Values.database.rootUser }}
{{- end -}}

{{/*
Database root password
*/}}
{{- define "database.rootPassword" -}}
{{ .Values.database.rootPassword }}
{{- end -}}
