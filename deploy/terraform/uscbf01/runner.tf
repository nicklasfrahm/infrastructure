data "google_compute_image" "ubuntu" {
  project = "ubuntu-os-cloud"
  family  = "ubuntu-2204-lts"
}

resource "google_compute_instance" "github_actions_runner" {
  name         = "github-actions-runner"
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
      # Provision an ephemeral IP address for the instance.
    }
  }

  service_account {
    email  = "compute@nicklasfrahm.iam.gserviceaccount.com"
    scopes = ["cloud-platform"]
  }

  # TODO: Install the GitHub Actions runner.
  metadata_startup_script = "echo hi > /test.txt"
}
