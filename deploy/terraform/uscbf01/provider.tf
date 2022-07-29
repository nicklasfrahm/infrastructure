terraform {
  backend "local" {
    path = "./uscbf01.tfstate"
  }
}

provider "google" {
  project = "nicklasfrahm"
  # Iowa was chosen, because it qualifies for the GCP Free Tier,
  # and is marked as "low-carbon", while being closest to Europe.
  region = "us-central1"
  zone   = "us-central1-a"
}
