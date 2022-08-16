# Google VM

This module creates a Google Cloud VM. By default it uses the `e2-micro` machine type, which qualifies for the Free Tier.

## Usage

```hcl
module "juliett" {
  source = "../modules/google_vm"

  hostname        = "one"
  fqdn            = "example.com"
  github_username = "octocat"
}
```
