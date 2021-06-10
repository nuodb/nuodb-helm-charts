{{/*
Expand the name of the chart.
*/}}
{{- define "hockey.name" -}}
{{- default .Chart.Name .Values.hockey.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "hockey.fullname" -}}
{{- $domain := default "domain" .Values.admin.domain -}}
{{- if .Values.hockey.fullnameOverride -}}
{{- .Values.hockey.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.hockey.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- printf "%s-%s" .Release.Name .Values.admin.domain | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s-%s" .Release.Name .Values.admin.domain $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "hockey.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Return the proper HOCKEY image name
*/}}
{{- define "hockey.image" -}}
{{- $registryName := .Values.hockey.image.registry -}}
{{- $repositoryName := .Values.hockey.image.repository -}}
{{- $tag := .Values.hockey.image.tag | toString -}}
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
Create a default fully qualified database address for TE direct connections
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "database.address" -}}
{{- $database := .Values.database.name -}}
{{- $namespace := default .Release.Namespace .Values.admin.namespace -}}
{{- printf "%s-clusterip.%s.svc:48006" $database $namespace  | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Return the proper Docker Image Registry Secret Names
*/}}
{{- define "hockey.imagePullSecrets" -}}
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
{{- else if .Values.hockey.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.hockey.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- end -}}
{{- else if .Values.hockey.image.pullSecrets }}
imagePullSecrets:
{{- range .Values.hockey.image.pullSecrets }}
  - name: {{ . }}
{{- end }}
{{- end -}}
{{- end -}}
