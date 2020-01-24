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
{{- define "database.fullname" -}}
{{- $domain := default "domain" .Values.admin.domain -}}
{{- $cluster := default "cluster0" .Values.cloud.cluster.name -}}
{{- if .Values.database.fullnameOverride -}}
{{- .Values.database.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.database.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- printf "%s-%s-%s-%s" .Release.Name $domain $cluster .Values.database.name | trunc 43 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s-%s-%s-%s" .Release.Name $domain $cluster .Values.database.name $name | trunc 43 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "nuodb.chart" -}}
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
Import ENV vars from configMaps
**BEWARE!!**
   The values for envFrom are formated into a single line because some parsers
   - either in k8s or rancher - throw errors occasionally if the multi-line format is used.
   You Have Been Warned.
*/}}
{{- define "restore.envFrom" }}
envFrom: [ configMapRef: { name: {{ .Values.database.name }}-restore } {{- range $map := .Values.restore.envFrom.configMapRef }}, configMapRef: { name: {{$map}} } {{- end }} ]
{{- end -}}

{{/*
Return the restore target.
*/}}
{{- define "restore.target" -}}
{{- default .Values.database.name .Values.restore.target -}}
{{- end -}}

{{/*
Return the restore source.
*/}}
{{- define "restore.source" -}}
{{- default ":latest" .Values.restore.source | trimSuffix "-" -}}
{{- end -}}
