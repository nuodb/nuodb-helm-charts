{{/*
Docker image for aws sidecar for restore from snapshot
*/}}
{{- define "aws.sidecar" -}}
{{- $registryName :=  default "docker.io" .Values.backupmgr.image.registry -}}
{{- $repositoryName := .Values.backupmgr.image.repository -}}
{{- $tag := default "latest" .Values.backupmgr.image.tag | toString -}}
{{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
{{- end -}}

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

