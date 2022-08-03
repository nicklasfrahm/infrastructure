variable "domain" {
  description = "The domain name of the GitHub pages site."
  type        = string

  validation {
    condition     = substr(var.domain, -1, -1) == "."
    error_message = "The domain name must end with a dot."
  }
}

variable "zone" {
  description = "The name of the DNS zone."
  type        = string
}

variable "organization" {
  description = "The name of the GitHub organization."
  type        = string
}

resource "google_dns_record_set" "a" {
  # Only create A records if the domain is an apex domain.
  count = length(var.domain) - length(replace(var.domain, ".", "")) > 2 ? 0 : 1

  name         = var.domain
  managed_zone = var.zone
  type         = "A"
  ttl          = 600

  rrdatas = [
    "185.199.108.153",
    "185.199.109.153",
    "185.199.110.153",
    "185.199.111.153",
  ]
}

resource "google_dns_record_set" "cname" {
  name         = "www.${var.domain}"
  managed_zone = var.zone
  type         = "CNAME"
  ttl          = 600

  rrdatas = ["${var.organization}.github.io."]
}
