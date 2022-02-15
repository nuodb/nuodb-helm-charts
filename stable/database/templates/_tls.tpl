{{- define "tls.client.mounts" -}}
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
{{- end -}}

{{- define "tls.mounts" -}}
{{- include "tls.client.mounts" . }}
{{- if .Values.admin.tlsKeyStore }}
- name: tls-keystore
  mountPath: /etc/nuodb/keys/nuoadmin.p12
  subPath: {{ .Values.admin.tlsKeyStore.key }}
{{- end }}
{{- end -}}


{{- define "tls.client.volumes" -}}
{{- if .Values.admin.tlsCACert }}
- name: tls-ca-cert
  secret:
    secretName: {{ .Values.admin.tlsCACert.secret }}
    defaultMode: 0440
{{- end }}
{{- if .Values.admin.tlsClientPEM }}
- name: tls-client-pem
  secret:
    secretName: {{ .Values.admin.tlsClientPEM.secret }}
    defaultMode: 0440
{{- end }}
{{- end -}}

{{- define "tls.volumes" -}}
{{- include "tls.client.volumes" . }}
{{- if .Values.admin.tlsKeyStore }}
- name: tls-keystore
  secret:
    secretName: {{ .Values.admin.tlsKeyStore.secret }}
    defaultMode: 0440
{{- end }}
{{- end -}}

{{- define "tls.client.env" -}}
{{- if .Values.admin.tlsKeyStore }}
- name: NUOCMD_VERIFY_SERVER
  value: /etc/nuodb/keys/ca.cert
- name: NUOCMD_CLIENT_KEY
  value: /etc/nuodb/keys/nuocmd.pem
- name: PATH
  value: /usr/bin:/bin:/usr/local/bin:/opt/nuodb/bin
- name: NUODB_HOME
  value: /opt/nuodb
{{- end }}
{{- end -}}


NUOCMD_PROCESS_FORMAT=[{engine_type}] {address::<UNKNOWN ADDRESS>} [start_id = {start_id}] [archive_id = {archive_id::}] [server_id = {server_id}] [pid = {pid::}] [node_id = {node_id::}]  [last_ack = {last_ack:5.2f:>60}] {durable_state}:{engine_state}
NUOCMD_API_SERVER=nuodb.nuodb.svc:8888
NUOCMD_VERIFY_SERVER=

{{/*
Add to environment list variables related to enabling TLS.
*/}}
{{- define "tls.env" -}}
{{- if .Values.admin.tlsKeyStore }}
{{- if .Values.admin.tlsKeyStore.password }}
- { name: NUODOCKER_KEYSTORE_PASSWORD,    value: {{ .Values.admin.tlsKeyStore.password | quote }} }
{{- end }}
{{- end }}
{{- end -}}
