{{/*
Return the proper insights image name
*/}}
{{- define "insights.image" -}}
{{- $registryName := $.Values.insights.image.registry -}}
{{- $repositoryName := $.Values.insights.image.repository -}}
{{- $tag := $.Values.insights.image.tag | toString -}}
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
Return the proper insights image name
*/}}
{{- define "insights.watcher" -}}
{{- $registryName := .Values.insights.watcher.registry -}}
{{- $repositoryName := .Values.insights.watcher.repository -}}
{{- $tag := .Values.insights.watcher.tag | toString -}}
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

{{- define "nuodb.sidecar" -}}
{{- if and $.Values.insights.enabled $.Values.admin.insights }}
- name: insights
  image: {{ template "insights.image" . }}
  imagePullPolicy: {{ $.Values.insights.image.pullPolicy }}
  tty: true
  volumeMounts:
  - mountPath: /etc/telegraf/telegraf.d/
    name: insights-config
  - name: log-volume
    mountPath: /var/log/nuodb
- name: insights-config
  image: {{ template "insights.watcher" . }}
  imagePullPolicy: {{ $.Values.insights.watcher.pullPolicy }}
  env:
  - name: LABEL
    value: {{ $.Values.admin.insights | quote }}
  - name: FOLDER
    value: /etc/telegraf/telegraf.d/
  - name: REQ_URL
    value: http://127.0.0.1:5000/reload
  volumeMounts:
  - name: insights-config
    mountPath: /etc/telegraf/telegraf.d/
  - name: log-volume
    mountPath: /var/log/nuodb
shareProcessNamespace: true
{{- end }}
{{- end -}}

{{- define "nuodb.sidecar.volumes" -}}
{{- if and .Values.admin.insights .Values.insights.enabled }}
- name: insights-config
  emptyDir: {}
{{- end }}
{{- end -}}