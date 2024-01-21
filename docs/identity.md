# Identity

This document describes how identity and access management is handled.

## Scenarios

- **A service running as a Kubernetes workload**  
  To support this scenario, the Kubernetes cluster must have the `charts/workload-identity` Helm chart installed. Additionally, the service must be able to obtain a service account token from the Kubernetes API server. The service account token can then be exchanged with other OpenID Connect providers to authenticate to other services.
