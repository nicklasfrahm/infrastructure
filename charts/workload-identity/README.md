# Workload Identity

This helm chart enables unauthenticated access to the OpenID discovery endpoint of the Kubernetes API server. This is required if you want to authenticate to another service using a Kubernetes service account. For reference, see the [OIDC issuer discovery section in the official Kubernetes documentation][oidc-issuer-discovery].

## Prerequisites

- Kubernetes 1.21+
- Kubernetes API server must be reachable from the federating OIDC provider
- Anonymous authentication enabled via `--anonymous-auth=true` Kubernetes API server flag

Additionally, it may be helpful to explicitly set the following two API server flags:

- `--service-account-issuer=https://k8s.example.com:6443`
- `--service-account-jwks-uri=https://k8s.example.com:6443/openid/v1/jwks`

[oidc-issuer-discovery]: https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/#service-account-issuer-discovery
