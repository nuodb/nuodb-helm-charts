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
{{- $domain := default "domain" .Values.admin.domain -}}
{{- $cluster := default "cluster0" .Values.cloud.cluster.name -}}
{{- if .Values.database.fullnameOverride -}}
{{- .Values.database.fullnameOverride | trunc 50 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.database.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- printf "%s-%s-%s-%s" .Release.Name $domain $cluster .Values.database.name | trunc 50 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s-%s-%s-%s" .Release.Name $domain $cluster .Values.database.name $name | trunc 50 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
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
{{- $domain := default "nuodb" .Values.admin.domain -}}
{{- $namespace := default .Release.Namespace .Values.admin.namespace -}}
{{- printf "%s.%s.svc" $domain $namespace  | trunc 63 | trimSuffix "-" -}}
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
Resolve the os.user
*/}}
{{- define "os.user" -}}
{{- if eq (include "defaultfalse" .Values.database.securityContext.enabled) "true" -}}
  {{ .Values.database.securityContext.runAsUser }}
{{- else -}}
   "1000"
{{- end -}}
{{- end -}}

{{/*
Resolve the os.group
*/}}
{{- define "os.group" -}}
{{- if eq (include "defaultfalse" .Values.database.securityContext.enabled) "true" -}}
  {{ .Values.database.securityContext.fsGroup }}
{{- else -}}
   "0"
{{- end -}}
{{- end -}}

{{/*
Add capabilities in a securityContext
*/}}
{{- define "database.capabilities" -}}
{{- with .Values.database.securityContext.capabilities }}
securityContext:
  capabilities:
    add: {{ . }}
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
envFrom: [ configMapRef: { name: {{ .Values.database.name }}-restore } {{- range $map := .Values.database.envFrom.configMapRef }}, configMapRef: { name: {{$map}} } {{- end }} ]
{{- end -}}

{{/*
Return options as $key $value
*/}}
{{- define "opt.key-values" -}}
{{- range $opt, $val := . }} {{$opt}} {{$val}} {{- end}}
{{- end -}}

{{/*
Return the hotcopy group
*/}}
{{- define "hotcopy.groupPrefix" -}}
{{ default .Values.cloud.cluster.name .Values.database.sm.hotCopy.backupGroupPrefix }}
{{- end -}}

{{/*
Return the hotcopy cronjob schedule by hotcopy type. It will take into account
any schedule overwrites configured per backup group.
*/}}
{{- define "hotcopy.schedule" -}}
  {{- $scheduleProp := printf "%sSchedule" .hotcopyType -}}
  {{- $val := index .Values.database.sm.hotCopy ( print $scheduleProp ) -}}
  {{- $overwrite := "" -}}
  {{- if eq .hotcopyType "journal" -}}
    {{- $val = .Values.database.sm.hotCopy.journalBackup.journalSchedule -}}
  {{- end -}}
  {{- if .Values.database.sm.hotCopy.backupGroups -}}
    {{- $group := (index .Values.database.sm.hotCopy.backupGroups ( print .backupGroup )) -}}
    {{- if $group -}}
      {{- $overwrites := index $group "overwrites" -}}
      {{- if $overwrites -}}
        {{- $overwrite = index $overwrites ( print $scheduleProp ) -}}
      {{- end -}}
    {{- end -}}
  {{- end -}}
{{ default $val $overwrite }}
{{- end -}}


{{/*
Return backup group selector
*/}}
{{- define "hotcopy.group.labels" -}}
  {{- $defaultGroupSelector := dict "pod-name" (printf "sm-%s-hotcopy-%d" (include "database.fullname" .) .index) -}}
  {{- $groupSelector := dict -}}
  {{- if .Values.database.sm.hotCopy.backupGroups -}}
    {{- $group := index .Values.database.sm.hotCopy.backupGroups (print .backupGroup) -}}
    {{- if $group -}}
      {{- $groupSelector = index $group "labelSelector" -}}
    {{- end -}}
  {{- end -}}
{{- include "opt.key-values" (default $defaultGroupSelector $groupSelector) }}
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
"nuodb.com/automatic-database-protocol-upgrade.te-preference-policy": {{ . | quote }}
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