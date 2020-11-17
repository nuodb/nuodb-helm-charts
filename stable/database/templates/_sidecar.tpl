{{- define "nuodb.sidecar" -}}
{{- if .Values.nuocollector }}
{{- if .Values.nuocollector.enabled }}
- name: nuocollector
  image: {{ template "nuocollector.image" . }}
  imagePullPolicy: {{ .Values.nuocollector.image.pullPolicy }}
  tty: true
  volumeMounts:
  - mountPath: /etc/telegraf/telegraf.d/dynamic/
    name: nuocollector-config
  - name: log-volume
    mountPath: /var/log/nuodb
- name: nuocollector-config
  image: {{ template "nuocollector.watcher" . }}
  imagePullPolicy: {{ .Values.nuocollector.watcher.pullPolicy }}
  env:
  - name: LABEL
    value: "nuodb.com/nuocollector-plugin in ({{ template "database.fullname" $ }}, insights)"
  - name: FOLDER
    value: /etc/telegraf/telegraf.d/dynamic/
  - name: REQ_URL
    value: http://127.0.0.1:5000/reload
  volumeMounts:
  - name: nuocollector-config
    mountPath: /etc/telegraf/telegraf.d/dynamic/
  - name: log-volume
    mountPath: /var/log/nuodb
shareProcessNamespace: true
{{- end }}
{{- end }}
{{- end -}}

{{- define "tde.sidecar" -}}
{{- if .Values.database.tde }}
{{- if .Values.database.tde.enabled }}
- name: tdemonitor
  image: {{ template "nuodb.image" . }}
  imagePullPolicy: {{ .Values.nuodb.image.pullPolicy }}
  args:
    - "nuotde"
    - "update"
    - "--monitor"
  env:
  - name: DB_NAME
    valueFrom:
      secretKeyRef:
        name: {{ .Values.database.name }}.nuodb.com
        key: database-name
  - { name: NUOCMD_API_SERVER,   value: "{{ template "admin.address" . }}:8888" }
  - { name: NUODB_TDE_FILES_PATH, value: "{{ .Values.database.tde.filesMountPath | default "/etc/nuodb/tde"}}" }
  volumeMounts:
  - name: nuotde
    mountPath: /usr/local/bin/nuotde
    subPath: nuotde
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
  {{- if .Values.database.tde.targetPassword }}
  - name: tde-passwords-volume
    mountPath: {{ default "/etc/nuodb/tde" .Values.database.tde.filesMountPath }}
    readOnly: true
  {{- end }}
{{- end }}
{{- end }}
{{- end -}}

{{- define "nuodb.sidecar.volumes" -}}
{{- if .Values.nuocollector }}
{{- if .Values.nuocollector.enabled }}
- name: nuocollector-config
  emptyDir: {}
{{- end }}
{{- end }}
{{- end -}}

{{/*
Return the proper NuoDB Collector image name
*/}}
{{- define "nuocollector.image" -}}
{{- $registryName := .Values.nuocollector.image.registry -}}
{{- $repositoryName := .Values.nuocollector.image.repository -}}
{{- $tag := .Values.nuocollector.image.tag | toString -}}
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
Return the proper NuoDB Collector configuration watcher image name
*/}}
{{- define "nuocollector.watcher" -}}
{{- $registryName := .Values.nuocollector.watcher.registry -}}
{{- $repositoryName := .Values.nuocollector.watcher.repository -}}
{{- $tag := .Values.nuocollector.watcher.tag | toString -}}
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
