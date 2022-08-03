variable "router" {
  description = "The domain name of the site's router."
  type        = string

  validation {
    condition     = substr(var.router, -1, -1) == "."
    error_message = "The domain name must end with a dot."
  }
}

variable "zone" {
  description = "The name of the DNS zone."
  type        = string
}

variable "location" {
  description = "The location code of the site."
  type        = string
}

data "google_dns_managed_zone" "default" {
  name = var.zone
}

resource "google_dns_record_set" "base" {
  name         = "${var.location}.${data.google_dns_managed_zone.default.dns_name}"
  managed_zone = data.google_dns_managed_zone.default.name
  type         = "CNAME"
  ttl          = 600

  rrdatas = [var.router]
}

resource "google_dns_record_set" "wildcard" {
  name         = "*.${var.location}.${data.google_dns_managed_zone.default.dns_name}"
  managed_zone = data.google_dns_managed_zone.default.name
  type         = "CNAME"
  ttl          = 600

  rrdatas = [var.router]
}
