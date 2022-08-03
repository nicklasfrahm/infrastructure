variable "domain" {
  description = "The domain name of the DNS zone."
  type        = string

  validation {
    condition     = substr(var.domain, -1, -1) != "."
    error_message = "The domain name must not end with a dot."
  }
}

variable "description" {
  description = "Purpose of the DNS zone."
  type        = string
}

resource "google_dns_managed_zone" "default" {
  name        = replace(var.domain, ".", "-")
  dns_name    = "${var.domain}."
  description = var.description

  dnssec_config {
    state         = "on"
    non_existence = "nsec3"

    default_key_specs {
      algorithm  = "ecdsap384sha384"
      key_type   = "keySigning"
      key_length = 384
    }
  }
}

output "id" {
  description = "The ID of the DNS zone."
  value       = google_dns_managed_zone.default.id
}
