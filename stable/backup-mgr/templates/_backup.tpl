{{/*
Docker image for aws sidecar
*/}}
{{- define "aws.sidecar" -}}
{{- $registryName :=  default "docker.io" .Values.backupmgr.image.registry -}}
{{- $repositoryName := .Values.backupmgr.image.repository -}}
{{- $tag := default "latest" .Values.backupmgr.image.tag | toString -}}
{{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
{{- end -}}

{{- define "aws.sidecar.setup" -}}
image: {{ template "aws.sidecar" . }}
imagePullPolicy: {{ .Values.backupmgr.image.pullPolicy }}
env:
- name: POD_NAME
  valueFrom:
    fieldRef:
      fieldPath: metadata.name
- name: NAMESPACE
  valueFrom:
    fieldRef:
      fieldPath: metadata.namespace
- { name: NUODB_DOMAIN,      value: "{{ .Values.admin.domain }}" }
- { name: NUOCMD_API_SERVER, value: "{{ template "admin.address" . }}:8888" }
- { name: NUOCMD_PLUGINS,    value: "/opt/nuodb/etc/nuodocker.py" }
- { name: PATH,              value: "/usr/bin:/bin:/usr/local/bin:/opt/nuodb/bin" }
- { name: NUODB_HOME,        value: "/opt/nuodb" }
{{- include "tls.client.env" . }}
volumeMounts:
- name: varlog
  mountPath: /var/log/backupmgr
{{- include "tls.client.mounts" . }}
{{- end -}}

{{- define "restore.serviceaccount" -}}
{{- default "nuodb-restore-mgr"  .Values.backupmgr.restore.serviceAccount.name }}
{{- end -}}

{{- define "snapshot.serviceaccount" }}
{{- default "nuodb-backup-mgr"  .Values.backupmgr.snapshot.serviceAccount.name }}
{{- end -}}




