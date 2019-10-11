{{/*
Default and null-safe definition of allowVolumeExpansion.
*/}}
{{- define "storageClass.allowVolumeExpansion" -}}
{{- if .Values.storageClass }}
{{- $string := toString .Values.storageClass.allowVolumeExpansion | default "true" -}}
{{- if eq $string "false" }}
allowVolumeExpansion: false
{{- else }}
allowVolumeExpansion: true
{{- end }}
{{- else }}
allowVolumeExpansion: true
{{- end }}
{{- end -}}

{{/*
Amazon EBS encrypted parameter.
*/}}
{{- define "storageClass.encryptedFlag" -}}
{{- $encrypted := (index $.Values.storageClass ( print .className )).encrypted -}}
{{- if not $encrypted -}}
encrypted: "false"
{{- else -}}
encrypted: "true"
{{- end -}}
{{- end -}}

{{/*
Amazon EBS iopsPerGB parameter.
*/}}
{{- define "storageClass.iopsPerGB" -}}
{{- $iopsPerGB := (index $.Values.storageClass ( print .className )).iopsPerGB -}}
{{- if $iopsPerGB -}}
iopsPerGB: {{ $iopsPerGB | quote }}
{{- else -}}
iopsPerGB: "50"
{{- end -}}
{{- end -}}
