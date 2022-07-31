variable "RUNNER_TOKEN" {
  description = "The runner token to register the runner."
  type        = string
  sensitive   = true
}

module "runner" {
  source = "../modules/gce_github_actions_runner"

  runner = {
    version = "2.294.0"
    token   = var.RUNNER_TOKEN
  }

  github = {
    username   = "nicklasfrahm"
    repository = "infrastructure"
  }

  vm = {
    hostname     = "echo"
    fqdn         = "nicklasfrahm.xyz"
    machine_type = "e2-micro"
  }
}

output "ip" {
  description = "The IP address of the runner's VM."
  value       = module.runner.ip
}
