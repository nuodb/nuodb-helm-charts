{{/*
Expand the name of the chart.
*/}}
{{- define "admin.name" -}}
{{- default .Chart.Name .Values.admin.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 50 chars because some Kubernetes name fields are limited to 63 chars (by the DNS naming spec)
and we need to allow space for any suffixes that may be added, and the "-NN" where NN is the pod number.
If release name contains chart name it will be used as a full name.
*/}}
{{- define "admin.fullname" -}}
{{- $domain := default "domain" .Values.admin.domain -}}
{{- $cluster := default "cluster0" .Values.cloud.cluster.name -}}
{{- if .Values.admin.fullnameOverride -}}
{{- .Values.admin.fullnameOverride | trunc 50 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.admin.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- printf "%s-%s-%s" .Release.Name .Values.admin.domain $cluster | trunc 50 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s-%s-%s" .Release.Name .Values.admin.domain $cluster $name | trunc 50 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "admin.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
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
{{- $registryName :=  .Values.busybox.image.registry -}}
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
{{- printf "%s.%s.svc" $domain .Release.Namespace  | trunc 63 | trimSuffix "-" -}}
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
Add capabilities in a securityContext
*/}}
{{- define "admin.capabilities" -}}
{{- with .Values.admin.securityContext.capabilities }}
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
   You Have Been Warned!
*/}}
{{- define "admin.envFrom" }}
envFrom: [{{- range $index, $map := .Values.admin.envFrom.configMapRef }}{{if gt $index 0}},{{end}} configMapRef: { name: {{$map}} } {{ end }}]
{{- end -}}

{{/*
Define the cluster domains
*/}}
{{- define "cluster.domain" -}}
{{- .Values.cloud.cluster.domain | default "cluster.local" }}
{{- end -}}

{{- define "cluster.entrypointDomain" -}}
{{- .Values.cloud.cluster.entrypointDomain | default (include "cluster.domain" .) }}
{{- end -}}

{{/*
Define the fully qualified NuoDB Admin address for the domain entrypoint.
*/}}
{{- define "admin.entrypointFullname" -}}
{{- $domain := default "domain" .Values.admin.domain -}}
{{- $cluster := default "cluster0" .Values.cloud.cluster.entrypointName -}}
{{- if .Values.admin.fullnameOverride -}}
{{- .Values.admin.fullnameOverride | trunc 50 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.admin.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- printf "%s-%s-%s" .Release.Name .Values.admin.domain $cluster | trunc 50 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s-%s-%s" .Release.Name .Values.admin.domain $cluster $name | trunc 50 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{- define "nuodb.domainEntrypoint" -}}
{{ include "admin.entrypointFullname" . }}-0.{{ .Values.admin.domain }}.$(NAMESPACE).svc.{{ include "cluster.entrypointDomain" . }}
{{- end -}}

{{- define "nuodb.altAddress" -}}
$(POD_NAME).{{ .Values.admin.domain }}.$(NAMESPACE).svc.{{ include "cluster.domain" . }}
{{- end -}}