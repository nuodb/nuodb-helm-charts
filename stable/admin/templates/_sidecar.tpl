{{/*
Return the proper nuocollector image name
*/}}
{{- define "nuocollector.image" -}}
{{- $registryName := $.Values.nuocollector.image.registry -}}
{{- $repositoryName := $.Values.nuocollector.image.repository -}}
{{- $tag := $.Values.nuocollector.image.tag | toString -}}
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
Return the proper nuocollector config watcher image name
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

{{- define "admin.nuodb.sidecar.collector" -}}
{{- if .Values.nuocollector }}
{{- if eq (include "defaultfalse" .Values.nuocollector.enabled) "true" }}
- name: nuocollector
  image: {{ template "nuocollector.image" . }}
  imagePullPolicy: {{ $.Values.nuocollector.image.pullPolicy }}
  tty: true
  resources:
  {{- toYaml .Values.nuocollector.resources | trim | nindent 4 }}
  {{- include "admin.sc.containerSecurityContext" . | indent 2 }}
  volumeMounts:
  - mountPath: /etc/telegraf/telegraf.d/dynamic/
    name: eph-volume
    subPath: telegraf
  - mountPath: /var/log/nuodb
    {{- if eq (include "defaultfalse" .Values.admin.logPersistence.enabled) "true" }}
    name: log-volume
    {{- else }}
    name: eph-volume
    subPath: log
    {{- end }}
- name: nuocollector-config
  image: {{ template "nuocollector.watcher" . }}
  imagePullPolicy: {{ $.Values.nuocollector.watcher.pullPolicy }}
  resources:
  {{- toYaml .Values.nuocollector.resources | trim | nindent 4 }}
  {{- include "admin.sc.containerSecurityContext" . | indent 2 }}
  env:
  - name: LABEL
    value: "nuodb.com/nuocollector-plugin in ({{ template "admin.fullname" $ }}, insights)"
  - name: FOLDER
    value: /etc/telegraf/telegraf.d/dynamic/
  - name: REQ_URL
    value: http://127.0.0.1:5000/reload
  volumeMounts:
  - mountPath: /etc/telegraf/telegraf.d/dynamic/
    name: eph-volume
    subPath: telegraf
  - mountPath: /tmp
    name: eph-volume
    subPath: tmp-watcher
  - mountPath: /var/log/nuodb
    {{- if eq (include "defaultfalse" .Values.admin.logPersistence.enabled) "true" }}
    name: log-volume
    {{- else }}
    name: eph-volume
    subPath: log
    {{- end }}
shareProcessNamespace: true
{{- end }}
{{- end }}
{{- end -}}
