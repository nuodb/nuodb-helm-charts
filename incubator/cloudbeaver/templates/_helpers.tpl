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
Expand the name of the chart.
*/}}
{{- define "cloudbeaver.name" -}}
{{- default .Chart.Name .Values.cloudbeaver.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "cloudbeaver.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}


{{/*
Create the name of the service account to use
*/}}
{{- define "cloudbeaver.serviceAccountName" -}}
{{- if .Values.cloudbeaver.addServiceAccount }}
{{- default (include "cloudbeaver.fullname" .) .Values.cloudbeaver.serviceAccount }}
{{- else }}
{{- default "default" .Values.cloudbeaver.serviceAccount.name }}
{{- end }}
{{- end }}
