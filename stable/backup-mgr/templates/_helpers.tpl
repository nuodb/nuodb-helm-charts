{{/*
Expand the name of the chart.
*/}}
{{- define "backupmgr.name" -}}
{{- default .Chart.Name .Values.backupmgr.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 50 chars because some Kubernetes name fields are limited to 63 chars (by the DNS naming spec)
and we have to allow for added suffixes including "-hotcopy" and "-NN" where NN is the pod number.
*/}}
{{- define "backupmgr.fullname" -}}
{{- $domain := default "domain" .Values.admin.domain -}}
{{- $cluster := default "cluster0" .Values.cloud.cluster.name -}}
{{- if .Values.backupmgr.fullnameOverride -}}
{{- .Values.backupmgr.fullnameOverride | trunc 50 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.backupmgr.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- printf "%s-%s-%s-%s" .Release.Name $domain $cluster .Values.backupmgr.name | trunc 50 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s-%s-%s-%s" .Release.Name $domain $cluster .Values.backupmgr.name $name | trunc 50 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name as used by the chart label.
*/}}
{{- define "backupmgr.chart" -}}
{{- printf "%s" .Chart.Name | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
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
Get Pod securityContext (core/v1/PodSecurityContext)
*/}}
{{- define "securityContext" -}}
{{- if eq (include "defaultfalse" .Values.backupmgr.securityContext.enabled) "true" }}
securityContext:
  runAsUser: {{ default 1000 .Values.backupmgr.securityContext.runAsUser }}
  runAsGroup: 0
  fsGroup: {{ default 1000 .Values.backupmgr.securityContext.fsGroup }}
  {{- include "sc.fsGroupChangePolicy" . }}
{{- else if eq (include "defaultfalse" .Values.backupmgr.securityContext.runAsNonRootGroup) "true" }}
securityContext:
  runAsUser: 1000
  runAsGroup: 1000
  fsGroup: {{ default 1000 .Values.backupmgr.securityContext.fsGroup }}
  {{- include "sc.fsGroupChangePolicy" . }}
{{- else if eq (include "defaultfalse" .Values.backupmgr.securityContext.fsGroupOnly) "true" }}
securityContext:
  fsGroup: {{ default 1000 .Values.backupmgr.securityContext.fsGroup }}
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
  {{- if eq (include "defaultfalse" .Values.backupmgr.securityContext.enabledOnContainer) "true" }}
securityContext:
  privileged: {{ include "defaultfalse" .Values.backupmgr.securityContext.privileged }}
  allowPrivilegeEscalation: {{ include "defaultfalse" .Values.backupmgr.securityContext.allowPrivilegeEscalation }}
  {{- include "sc.capabilities" . | indent 2 }}
  {{- end }}
{{- end -}}

{{/*
Get container securityContext defining capabilities
*/}}
{{- define "sc.capabilities" -}}
  {{- if .Values.backupmgr.securityContext.capabilities }}
    {{- if typeIs "[]interface {}" .Values.backupmgr.securityContext.capabilities }}
capabilities:
      {{- with .Values.backupmgr.securityContext.capabilities }}
  add: {{ . }}
      {{- end }}
    {{- else if or .Values.backupmgr.securityContext.capabilities.add .Values.backupmgr.securityContext.capabilities.drop }}
capabilities:
      {{- if .Values.backupmgr.securityContext.capabilities.add }}
  add:
        {{- toYaml .Values.backupmgr.securityContext.capabilities.add | trim | nindent 4 }}
      {{- end }}
      {{- if .Values.backupmgr.securityContext.capabilities.drop }}
  drop:
        {{- toYaml .Values.backupmgr.securityContext.capabilities.drop | trim | nindent 4 }}
      {{- end }}
    {{- end }}
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
Renders nuodocker options and flags. An option is rendered only if its value is
not empty. Flags can be defined by setting their value to boolean true or "true"
*/}}
{{- define "backupmgr.otherOptions" -}}
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

