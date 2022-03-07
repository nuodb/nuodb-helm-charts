{{/*
Docket image for aws sidecar for restore from snapshot
*/}}
{{- define "aws.sidecar" -}}
{{- $registryName :=  default "docker.io" .Values.backupmgr.image.registry -}}
{{- $repositoryName := .Values.backupmgr.image.repository -}}
{{- $tag := default "latest" .Values.backupmgr.image.tag | toString -}}
{{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
{{- end -}}
