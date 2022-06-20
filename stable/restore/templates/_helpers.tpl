{{/*
Expand the name of the chart.
*/}}
{{- define "nuodb.name" -}}
{{- default .Chart.Name .Values.restore.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "nuodb.fullname" -}}
{{- $name := default .Chart.Name .Values.restore.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "restore.fullname" -}}
{{- $domain := default "domain" .Values.admin.domain -}}
{{- $cluster := default "cluster0" .Values.cloud.cluster.name -}}
{{- $target := include "restore.target" . -}}
{{- if .Values.restore.fullnameOverride -}}
{{- .Values.restore.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.restore.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- printf "%s-%s-%s-%s" .Release.Name $domain $cluster $target | trunc 43 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s-%s-%s-%s" .Release.Name $domain $cluster $target $name | trunc 43 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "nuodb.chart" -}}
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
{{- else if .Values.nuodb.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.nuodb.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- end -}}
{{- else if .Values.nuodb.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.nuodb.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- end -}}
{{- end -}}

{{/*
Get Pod securityContext
*/}}
{{- define "securityContext" -}}
{{- if eq (include "defaultfalse" .Values.database.securityContext.enabled) "true" }}
securityContext:
  runAsUser: {{ default 1000 .Values.database.securityContext.runAsUser }}
  runAsGroup: 0
  fsGroup: {{ default 1000 .Values.database.securityContext.fsGroup }}
  {{- include "sc.fsGroupChangePolicy" . }}
{{- else if eq (include "defaultfalse" .Values.database.securityContext.runAsNonRootGroup) "true" }}
securityContext:
  runAsUser: 1000
  runAsGroup: 1000
  fsGroup: {{ default 1000 .Values.database.securityContext.fsGroup }}
  {{- include "sc.fsGroupChangePolicy" . }}
{{- else if eq (include "defaultfalse" .Values.database.securityContext.fsGroupOnly) "true" }}
securityContext:
  fsGroup: {{ default 1000 .Values.database.securityContext.fsGroup }}
  {{- include "sc.fsGroupChangePolicy" . }}
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
{{- define "sc.containerSecurityContext" }}
{{- if eq (include "defaultfalse" .Values.database.securityContext.enabled) "true" }}
securityContext:
  privileged: {{ include "defaultfalse" .Values.database.securityContext.privileged }}
  allowPrivilegeEscalation: {{ include "defaultfalse" .Values.database.securityContext.allowPrivilegeEscalation }}
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
Import ENV vars from configMaps
**BEWARE!!**
   The values for envFrom are formated into a single line because some parsers
   - either in k8s or rancher - throw errors occasionally if the multi-line format is used.
   You Have Been Warned.
*/}}
{{- define "restore.envFrom" }}
envFrom: [ configMapRef: { name: {{ .Values.admin.domain }}-{{ include "restore.target" . }}-restore } {{- range $map := .Values.restore.envFrom.configMapRef }}, configMapRef: { name: {{$map}} } {{- end }} ]
{{- end -}}

{{/*
Return the restore target.
*/}}
{{- define "restore.target" -}}
{{- if .Values.database -}}
{{- default .Values.database.name .Values.restore.target -}}
{{- else -}}
{{- .Values.restore.target -}}
{{- end -}}
{{- end -}}

{{/*
Return the restore source.
*/}}
{{- define "restore.source" -}}
{{- if hasPrefix ":" .Values.restore.source }}
{{- $valid := list ":latest" "" }}
{{- if not (has .Values.restore.source $valid) }}
{{- fail (printf "Invalid autorestore source: %s" .Values.restore.source) }}
{{- end -}}
{{- end -}}
{{- default ":latest" .Values.restore.source | trimSuffix "-" -}}
{{- end -}}

{{/*
Return the arguments for nuorestore script which defines the archives 
selector. If archiveIds are specified, they take precedence over labels.
*/}}
{{- define "restore.archives" -}}
{{- if .Values.restore.archiveIds }}
- "--archive-ids"
- {{ join " " .Values.restore.archiveIds | quote }}
{{- else if .Values.restore.labels }}
- "--labels"
{{- if typeIs "string" .Values.restore.labels }}
- {{ .Values.restore.labels }}
{{- else }}
- {{ range $opt, $val := .Values.restore.labels }} {{$opt}} {{$val}} {{- end}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Return the nuorestore script arguments section
*/}}
{{- define "restore.args" -}}
{{- $restoreTarget := include "restore.target" . }}
args:
- "nuorestore"
- "--type"
- "{{ default "database" .Values.restore.type }}"
- "--db-name"
- "{{ $restoreTarget }}"
- "--source"
- "{{ include "restore.source" . }}"
- "--auto"
{{- if hasKey .Values.restore "autoRestart" }}
- {{ .Values.restore.autoRestart | quote }}
{{ else }}
- "true"
{{- end -}}
- "--manual"
{{- if hasKey .Values.restore "manual" }}
- {{ .Values.restore.manual | quote }}
{{ else }}
- "false"
{{- end }}
{{ template "restore.archives" . }}
{{- end -}}

{{/*
Renders the labels for all resources deployed by this Helm chart
*/}}
{{- define "restore.resourceLabels" -}}
app: {{ template "restore.fullname" . }}
group: nuodb
subgroup: restore
domain: {{ .Values.admin.domain }}
database: {{ template "restore.target" . }}
chart: {{ template "nuodb.chart" . }}
release: {{ .Release.Name | quote }}
{{- range $k, $v := .Values.restore.resourceLabels }}
"{{ $k }}": "{{ $v | quote }}"
{{- end }}
{{- end -}}

{{/*
Renders the name of the Secret for the database selected for restore
*/}}
{{- define "database.secretName" -}}
{{ .Values.admin.domain }}-{{ template "restore.target" . }}
{{- end -}}