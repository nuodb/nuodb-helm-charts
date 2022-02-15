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
{{- end }}
{{- end -}}

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
