{{/*
Generic helper function to turn a list into a comma separated string
*/}}
{{- define "helm-toolkit.utils.joinListWithComma" -}}
{{- $local := dict "first" true -}}
{{- range $k, $v := . -}}{{- if not $local.first -}},{{- end -}}{{- $v -}}{{- $_ := set $local "first" false -}}{{- end -}}
{{- end -}}

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
Create chart name as used by the chart label.
*/}}
{{- define "admin.chart" -}}
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
Get Pod securityContext
*/}}
{{- define "securityContext" -}}
{{- if eq (include "defaultfalse" .Values.admin.securityContext.enabled) "true" }}
securityContext:
  runAsUser: {{ default 1000 .Values.admin.securityContext.runAsUser }}
  runAsGroup: 0
  fsGroup: {{ default 1000 .Values.admin.securityContext.fsGroup }}
  {{- include "sc.fsGroupChangePolicy" . }}
{{- else if eq (include "defaultfalse" .Values.admin.securityContext.runAsNonRootGroup) "true" }}
securityContext:
  runAsUser: 1000
  runAsGroup: 1000
  fsGroup: {{ default 1000 .Values.admin.securityContext.fsGroup }}
  {{- include "sc.fsGroupChangePolicy" . }}
{{- else if eq (include "defaultfalse" .Values.admin.securityContext.fsGroupOnly) "true" }}
securityContext:
  fsGroup: {{ default 1000 .Values.admin.securityContext.fsGroup }}
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
  {{- if eq (include "defaultfalse" .Values.admin.securityContext.enabledOnContainer) "true" }}
securityContext:
  privileged: {{ include "defaultfalse" .Values.admin.securityContext.privileged }}
  allowPrivilegeEscalation: {{ include "defaultfalse" .Values.admin.securityContext.allowPrivilegeEscalation }}
  {{- include "sc.capabilities" . | indent 2 }}
  {{- end }}
{{- end -}}

{{/*
Get container securityContext defining capabilities
*/}}
{{- define "sc.capabilities" -}}
  {{- if .Values.admin.securityContext.capabilities }}
    {{- if typeIs "[]interface {}" .Values.admin.securityContext.capabilities }}
capabilities:
      {{- with .Values.admin.securityContext.capabilities }}
  add: {{ . }}
      {{- end }}
    {{- else if or .Values.admin.securityContext.capabilities.add .Values.admin.securityContext.capabilities.drop }}
capabilities:
      {{- if .Values.admin.securityContext.capabilities.add }}
  add:
        {{- toYaml .Values.admin.securityContext.capabilities.add | trim | nindent 4 }}
      {{- end }}
      {{- if .Values.admin.securityContext.capabilities.drop }}
  drop:
        {{- toYaml .Values.admin.securityContext.capabilities.drop | trim | nindent 4 }}
      {{- end }}
    {{- end }}
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

{{/*
Imports NuoAdmin global load balancer configuration via annotations.
The configuration is imported only in the entrypoint cluster.
*/}}
{{- define "admin.loadBalancerConfig" -}}
{{- if .Values.admin.lbConfig }}
{{- if (eq (default "cluster0" .Values.cloud.cluster.name) (default "cluster0" .Values.cloud.cluster.entrypointName)) }}
{{- with .Values.admin.lbConfig.fullSync }}
"nuodb.com/sync-load-balancer-config": {{ . | quote }}
{{- end -}}
{{- with .Values.admin.lbConfig.prefilter }}
"nuodb.com/load-balancer-prefilter": {{ . | quote }}
{{- end -}}
{{- with .Values.admin.lbConfig.default }}
"nuodb.com/load-balancer-default": {{ . | quote }}
{{- end -}}
{{- with .Values.admin.lbConfig.policies }}
{{- range $opt, $val := . }}
"nuodb.com/load-balancer-policy.{{ $opt }}": {{ $val | quote }}
{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}
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
Renders the admin service name for external access based on the service type
*/}}
{{- define "admin.externalServiceName" -}}
  {{- $serviceType := (default "LoadBalancer" .Values.admin.externalAccess.type) -}}
  {{- if eq $serviceType "LoadBalancer" -}}
{{ .Values.admin.domain }}-{{ .Values.admin.serviceSuffix.balancer }}
  {{- else if eq $serviceType "NodePort" -}}
{{ .Values.admin.domain }}-{{ .Values.admin.serviceSuffix.nodeport }}
  {{- else -}}
{{ .Values.admin.domain }}
  {{- end }}
{{- end }}

{{/*
Renders the annotations for the LoadBalancer admin service
*/}}
{{- define "admin.externalAccessAnnotations" -}}
  {{- if eq (default "LoadBalancer" .Values.admin.externalAccess.type) "LoadBalancer" }}
    {{- if .Values.admin.externalAccess.annotations }}
{{ toYaml .Values.admin.externalAccess.annotations | trim }}
    {{- else -}}
      {{- if .Values.cloud.provider }}
        {{- if eq .Values.cloud.provider "amazon" }}
          {{- if .Values.admin.externalAccess.internalIP }}
service.beta.kubernetes.io/aws-load-balancer-internal: "true"
service.beta.kubernetes.io/aws-load-balancer-scheme: "internal"
          {{- else }}
service.beta.kubernetes.io/aws-load-balancer-type: "external"
service.beta.kubernetes.io/aws-load-balancer-nlb-target-type: "ip"
service.beta.kubernetes.io/aws-load-balancer-scheme: "internet-facing"
          {{- end }}
        {{- else if eq .Values.cloud.provider "azure" }}
          {{- if .Values.admin.externalAccess.internalIP }}
service.beta.kubernetes.io/azure-load-balancer-internal: "true"
          {{- end }}
        {{- else if eq .Values.cloud.provider "google" }}
          {{- if .Values.admin.externalAccess.internalIP }}
cloud.google.com/load-balancer-type: "Internal"
networking.gke.io/load-balancer-type: "Internal"
          {{- end -}}
        {{- end -}}
      {{- end -}}
    {{- end }}
  {{- end }}
{{- end -}}
