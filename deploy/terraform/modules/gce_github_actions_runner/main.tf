resource "google_compute_address" "runner" {
  name         = var.vm.hostname
  network_tier = "STANDARD"
}

data "google_compute_image" "ubuntu" {
  project = "ubuntu-os-cloud"
  family  = "ubuntu-2204-lts"
}

data "http" "ssh_keys" {
  url = "https://github.com/${var.github.username}.keys"
}

resource "google_compute_instance" "runner" {
  name         = var.vm.hostname
  machine_type = "e2-micro"

  # desc

  boot_disk {
    initialize_params {
      image = data.google_compute_image.ubuntu.self_link
      size  = 25
    }
  }

  network_interface {
    network = "default"
    access_config {
      nat_ip       = google_compute_address.runner.address
      network_tier = google_compute_address.runner.network_tier
    }
  }

  metadata = {
    "ssh-keys" = "${var.github.username}:${data.http.ssh_keys.response_body}"
    # Ensure that the runner follows the same naming conventions as all
    # other machines of my infrastructure.
    "user-data" = <<EOT
#cloud-config

hostname: ${var.vm.hostname}
fqdn: ${var.vm.hostname}.${var.vm.fqdn}
prefer_fqdn_over_hostname: true

runcmd:
  - |
    #!/bin/bash

    # Create a folder.
    mkdir /actions-runner && cd /actions-runner

    # Download the latest runner package.
    curl -o actions-runner-linux-x64-${var.runner.version}.tar.gz -L https://github.com/actions/runner/releases/download/v${var.runner.version}/actions-runner-linux-x64-${var.runner.version}.tar.gz

    # Extract the installer.
    tar xzf ./actions-runner-linux-x64-${var.runner.version}.tar.gz

    # Configure the runner and start the configuration experience.
    RUNNER_ALLOW_RUNASROOT=1 ./config.sh --url https://github.com/${var.github.username}/${var.github.repository} --token ${var.runner.token} --unattended --ephemeral

    # Install the runner as a systemd service.
    RUNNER_ALLOW_RUNASROOT=1 ./svc.sh install

    # Start the runner.
    RUNNER_ALLOW_RUNASROOT=1 ./svc.sh start

    # Check status of the runner.
    RUNNER_ALLOW_RUNASROOT=1 ./svc.sh start
EOT
  }

  # Using this ensures that a new VM will be provisioned if the version of the runner changes.
  metadata_startup_script = "echo ${var.runner.version} > /dev/null"

  depends_on = [google_compute_address.runner]
}
