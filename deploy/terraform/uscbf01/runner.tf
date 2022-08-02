variable "PERSONAL_ACCESS_TOKEN" {
  description = "A personal access token with the `repo` scope."
  type        = string
  sensitive   = true
}

module "runner" {
  source = "../modules/gce_github_actions_runner"

  runner = {
    version = "2.294.0"
    token   = var.PERSONAL_ACCESS_TOKEN
  }

  github = {
    username   = "nicklasfrahm"
    repository = "infrastructure"
  }

  vm = {
    hostname     = "delta"
    fqdn         = "nicklasfrahm.xyz"
    machine_type = "e2-micro"
  }
}

output "ip" {
  description = "The IP address of the runner's VM."
  value       = module.runner.ip
}
