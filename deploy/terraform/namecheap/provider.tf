terraform {
  backend "gcs" {
    bucket = "nicklasfrahm"
    prefix = "terraform/state"
  }

  required_providers {
    namecheap = {
      source  = "namecheap/namecheap"
      version = ">= 2.0.0"
    }
  }
}

provider "namecheap" {
  user_name = "nicklasfrahm"
  api_user  = "nicklasfrahm"
}
