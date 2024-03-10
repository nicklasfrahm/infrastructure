# Observability client

This helm chart installs the following components:

- `victoria-metrics-operator`: An operator for managing [Victoria Metrics][victoria-metrics] deployments.
- `victoria-metrics-agent`: An agent that scrapes metrics and sends them to a [Victoria Metrics][victoria-metrics] instance.
- `kube-state-metrics`: A service that listens to the Kubernetes API server and generates metrics about the state of the objects.
- `prometheus-node-exporter`: A service that collects metrics about the node it runs on.

[victoria-metrics]: https://victoriametrics.com/
