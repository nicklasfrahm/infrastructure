# GCE GitHub Actions Runner

This module creates a new GitHub action runner using the Google Compute Engine (GCE).

## Usage

```hcl
module "echo" {
  source = "../modules/gce_github_actions_runner"

  runner {
    token   = env.TF_GCE_GITHUB_ACTIONS_RUNNER_TOKEN
    version = "2.294.0"
  }

  github {
    username   = "nicklasfrahm"
    repository = "infrastructure"
  }

  vm {
    hostname = "runner"
    fqdn     = "local"
    # By default, the module will use the "e2-micro"
    # machine type that is part of the GCP Free Tier:
    # https://cloud.google.com/free/docs/gcp-free-tier/#free-tier-usage-limits
    machine_type = "e2-micro"
  }
}
```

## Outputs

| Name | Type     | Description                      |
| ---- | -------- | -------------------------------- |
| `ip` | `string` | The IP address of the runner VM. |
