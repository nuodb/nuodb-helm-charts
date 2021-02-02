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
{{- if .Values.database.securityContext.enabled -}}
  {{ .Values.database.securityContext.runAsUser }}
{{- else -}}
   "1000"
{{- end -}}
{{- end -}}

{{/*
Resolve the os.group
*/}}
{{- define "os.group" -}}
{{- if .Values.database.securityContext.enabled -}}
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
{{- define "hotcopy.group" -}}
{{ default .Values.cloud.cluster.name .Values.database.sm.hotCopy.backupGroup }}
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
Takes a boolean as argument return it's value if it was defined or return true otherwise
Note: Sprig's default function on an empty/not defined variable returns false, workaround
it by calling typeIs "bool" https://github.com/Masterminds/sprig/issues/111
*/}}
{{- define "defaulttrue" -}}
{{- if typeIs "bool" . -}}
{{- . -}}
{{- else -}}
{{- default true . -}}
{{- end -}}
{{- end -}}

{{/*
Imports per-database load balancer configuration via annotations.
The configuration is imported only in the entrypoint cluster.
*/}}
{{- define "database.loadBalancerConfig" -}}
{{- if .Values.database.lbConfig }}
{{- if (eq (default "cluster0" .Values.cloud.cluster.name) (default "cluster0" .Values.cloud.cluster.entrypointName)) }}
{{- with .Values.database.lbConfig.prefilter }}
"nuodb.com/load-balancer-prefilter": {{ . | quote }}
{{- end -}}
{{- with .Values.database.lbConfig.default }}
"nuodb.com/load-balancer-default": {{ . | quote }}
{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Render database restore init container
*/}}
{{- define "database.restoreInitContainer" }}
- name: restore
  image: {{ template "nuodb.image" . }}
  imagePullPolicy: {{ .Values.nuodb.image.pullPolicy }}
  command: ['restorearchive']
  env:
  - name: DB_NAME
    valueFrom:
      secretKeyRef:
        name: {{ .Values.database.name }}.nuodb.com
        key: database-name
  - name: DATABASE_RESTORE_CREDENTIALS
    valueFrom:
      secretKeyRef:
        name: {{ .Values.database.name }}.nuodb.com
        key: database-restore-credentials
  - name: DATABASE_BACKUP_CREDENTIALS
    valueFrom:
      secretKeyRef:
        name: {{ .Values.database.name }}.nuodb.com
        key: database-backup-credentials
  - { name: NUOCMD_API_SERVER,   value: "{{ template "admin.address" . }}:8888" }
  - { name: NUODB_BACKUP_GROUP,  value: "{{ include "hotcopy.group" . }}" }
  envFrom: 
  - configMapRef: { name: {{ .Values.database.name }}-restore }
  volumeMounts:
  - name: log-volume
    mountPath: /var/log/nuodb
  - name: restore-common
    mountPath: /opt/nuodb/etc/restore_common.sh
    subPath: restore_common.sh
  - name: nuobackup
    mountPath: /usr/local/bin/nuobackup
    subPath: nuobackup
  - name: restore-archive
    mountPath: /usr/local/bin/restorearchive
    subPath: restorearchive
  - mountPath: /var/opt/nuodb/archive
    name: archive-volume
  {{- if .smType }}
  {{- if eq .smType "hotcopy" }}
  - mountPath: /var/opt/nuodb/backup
    name: backup-volume
  {{- end }}
  {{- end }}
  {{- if .Values.admin.tlsCACert }}
  - name: tls-ca-cert
    mountPath: /etc/nuodb/keys/ca.cert
    subPath: {{ .Values.admin.tlsCACert.key }}
  {{- end }}
  {{- if .Values.admin.tlsClientPEM }}
  - name: tls-client-pem
    mountPath: /etc/nuodb/keys/nuocmd.pem
    subPath: {{ .Values.admin.tlsClientPEM.key }}
  {{- end }}
  {{- if .Values.admin.tde }}
  {{- if .Values.admin.tde.secrets }}
  {{- if hasKey .Values.admin.tde.secrets .Values.database.name }}
  {{- range $dbName, $secret := .Values.admin.tde.secrets }}
  {{- if eq $dbName $.Values.database.name }}
  - name: tde-volume-{{ $dbName }}
    mountPath: {{ default "/etc/nuodb/tde" $.Values.admin.tde.storagePasswordsDir }}/{{ $dbName }}
    readOnly: true
  {{- end }}
  {{- end }}
  {{- end }}
  {{- end }}
  {{- end }}
{{- end -}}

{{/*
Render database auto import init container
*/}}
{{- define "database.importInitContainer" }}
- name: import
  image: {{ template "nuodb.image" . }}
  imagePullPolicy: {{ .Values.nuodb.image.pullPolicy }}
  command: ['importarchive']
  env:
  - name: DB_NAME
    valueFrom:
      secretKeyRef:
        name: {{ .Values.database.name }}.nuodb.com
        key: database-name
  - name: DATABASE_IMPORT_CREDENTIALS
    valueFrom:
      secretKeyRef:
        name: {{ .Values.database.name }}.nuodb.com
        key: database-import-credentials
  - { name: NUOCMD_API_SERVER,   value: "{{ template "admin.address" . }}:8888" }
  envFrom: 
  - configMapRef: { name: {{ .Values.database.name }}-restore }
  volumeMounts:
  - name: log-volume
    mountPath: /var/log/nuodb
  - name: restore-common
    mountPath: /opt/nuodb/etc/restore_common.sh
    subPath: restore_common.sh
  - name: import-archive
    mountPath: /usr/local/bin/importarchive
    subPath: importarchive
  - mountPath: /var/opt/nuodb/archive
    name: archive-volume
  {{- if .Values.admin.tlsCACert }}
  - name: tls-ca-cert
    mountPath: /etc/nuodb/keys/ca.cert
    subPath: {{ .Values.admin.tlsCACert.key }}
  {{- end }}
  {{- if .Values.admin.tlsClientPEM }}
  - name: tls-client-pem
    mountPath: /etc/nuodb/keys/nuocmd.pem
    subPath: {{ .Values.admin.tlsClientPEM.key }}
  {{- end }}
  {{- if .Values.admin.tde }}
  {{- if .Values.admin.tde.secrets }}
  {{- if hasKey .Values.admin.tde.secrets .Values.database.name }}
  {{- range $dbName, $secret := .Values.admin.tde.secrets }}
  {{- if eq $dbName $.Values.database.name }}
  - name: tde-volume-{{ $dbName }}
    mountPath: {{ default "/etc/nuodb/tde" $.Values.admin.tde.storagePasswordsDir }}/{{ $dbName }}
    readOnly: true
  {{- end }}
  {{- end }}
  {{- end }}
  {{- end }}
  {{- end }}
{{- end -}}

{{/*
Render database init containers
*/}}
{{- define "database.initContainers" }}
initContainers:
- name: init-disk
  image: {{ template "init.image" . }}
  imagePullPolicy: {{ default "" .Values.busybox.image.pullPolicy | quote }}
  command: ['chmod' , '770', '/var/opt/nuodb/archive', '/var/log/nuodb']
  volumeMounts:
  - name: archive-volume
    mountPath: /var/opt/nuodb/archive
  - name: log-volume
    mountPath: /var/log/nuodb
{{- include "database.restoreInitContainer" . }}
{{- include "database.importInitContainer" . }}
{{- end -}}