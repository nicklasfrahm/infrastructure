{{/* Create chart name and version as used by the chart label. */}}
{{- define "workload-identity.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/* Common labels. */}}
{{- define "workload-identity.labels" -}}
helm.sh/chart: {{ include "workload-identity.chart" . }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}
