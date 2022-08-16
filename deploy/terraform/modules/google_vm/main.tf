variable "github_username" {
  description = "The GitHub username used to fetch the SSH keys."
  type        = string
}

variable "hostname" {
  description = "The hostname for the VM."
  type        = string
}

variable "fqdn" {
  description = "The fully qualified domain name for the VM."
  type        = string
}

resource "google_compute_address" "vm" {
  name         = var.hostname
  network_tier = "STANDARD"
}

data "google_compute_image" "ubuntu" {
  project = "ubuntu-os-cloud"
  family  = "ubuntu-2204-lts"
}

data "http" "ssh_keys" {
  url = "https://github.com/${var.github_username}.keys"
}

resource "google_compute_instance" "vm" {
  name         = var.hostname
  machine_type = "e2-micro"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.ubuntu.self_link
      size  = 25
    }
  }

  network_interface {
    network = "default"
    access_config {
      nat_ip       = google_compute_address.vm.address
      network_tier = google_compute_address.vm.network_tier
    }
  }

  metadata = {
    "ssh-keys" = "${var.github_username}:${data.http.ssh_keys.response_body}"
    # Ensure that the runner follows the same naming conventions as all
    # other machines of my infrastructure.
    "user-data" = <<EOT
#cloud-config

hostname: ${var.hostname}
fqdn: ${var.hostname}.${var.fqdn}
prefer_fqdn_over_hostname: true
EOT
  }

  depends_on = [google_compute_address.vm]
}

output "ip" {
  description = "The IP address of the VM."
  value       = google_compute_address.vm.address
}
