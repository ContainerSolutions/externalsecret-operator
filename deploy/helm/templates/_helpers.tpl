{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "externalsecret-operator.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "externalsecret-operator.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "externalsecret-operator.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create the name of the service account to use
*/}}
{{- define "externalsecret-operator.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
    {{ default (include "externalsecret-operator.fullname" .) .Values.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/*
Create the name of the secret that will hold the config
*/}}
{{- define "externalsecret-operator.secretName" -}}
{{- if .Values.secret.create -}}
    {{ default (include "externalsecret-operator.fullname" .) .Values.secret.name }}
{{- else -}}
    {{ default "default" .Values.secret.name }}
{{- end -}}
{{- end -}}

{{/*
Create watchNamespace: if not specified assume is release namespace
*/}}
{{- define "externalsecret-operator.watchNamespace" -}}
{{- if .Values.watchNamespace -}}
    {{ default .Values.watchNamespace }}
{{- else -}}
    {{ default .Release.Namespace }}
{{- end -}}
{{- end -}}
