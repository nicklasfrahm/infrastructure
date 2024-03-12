# Observability UI

This helm chart installs the following components:

- `grafana`: A single pane of glass for metrics, logs, and traces.

## Prerequisites

Make sure that the namespace where you install this helm chart has the following secrets.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: o11y-ui-grafana-oauth
  labels:
    app.kubernetes.io/name: grafana
    app.kubernetes.io/instance: o11y-ui
    app.kubernetes.io/part-of: o11y-ui
type: Opaque
stringData:
  CLIENT_ID: REDACTED_GOOGLE_CLIENT_ID
  CLIENT_SECRET: REDACTED_GOOGLE_CLIENT_SECRET
---
apiVersion: v1
kind: Secret
metadata:
  name: o11y-ui-grafana-admin
  labels:
    app.kubernetes.io/name: grafana
    app.kubernetes.io/instance: o11y-ui
    app.kubernetes.io/part-of: o11y-ui
type: Opaque
stringData:
  admin-user: admin
  admin-password: REDACTED_ADMIN_PASSWORD
```
