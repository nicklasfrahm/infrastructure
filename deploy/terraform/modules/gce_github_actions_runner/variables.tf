variable "github" {
  description = "(Required) Configuration for the runner and its connection to GitHub."

  type = object({
    username   = string
    repository = string
  })
}

variable "vm" {
  description = "(Optional) General settings for the runner's VM."

  type = object({
    hostname     = string
    fqdn         = string
    machine_type = string
  })

  default = {
    hostname     = "runner"
    fqdn         = "local"
    machine_type = "e2-micro"
  }
}

variable "runner" {
  description = "(Required) Configuration for the runner."
  sensitive   = true

  type = object({
    version = string
    # This is not a runner registration token, but a personal access token with
    # the `repo` scope, because the runner registration token is too short lived.
    token = string
  })
}
