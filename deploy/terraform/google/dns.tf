module "google_dns_zone_nicklasfrahm_dev" {
  source = "../modules/google_dns_zone"

  domain      = "nicklasfrahm.dev."
  description = "My customer-facing domain"
}

module "google_dns_github_pages_nicklasfrahm_dev" {
  source = "../modules/google_dns_github_pages"

  organization = "nicklasfrahm"
  domain       = "nicklasfrahm.dev."
  zone         = module.google_dns_zone_nicklasfrahm_dev.name
}

module "google_dns_site_dktil01" {
  source = "../modules/google_dns_site"

  zone     = module.google_dns_zone_nicklasfrahm_dev.name
  location = "dktil01"
  router   = "alfa.nicklasfrahm.xyz."
}

module "google_dns_site_dksjb01" {
  source = "../modules/google_dns_site"

  zone     = module.google_dns_zone_nicklasfrahm_dev.name
  location = "dksjb01"
  router   = "mike.nicklasfrahm.xyz."
}

module "google_dns_site_dksjb02" {
  source = "../modules/google_dns_site"

  zone     = module.google_dns_zone_nicklasfrahm_dev.name
  location = "dksjb02"
  router   = "mike.nicklasfrahm.xyz."
}

module "google_dns_site_deflf01" {
  source = "../modules/google_dns_site"

  zone     = module.google_dns_zone_nicklasfrahm_dev.name
  location = "deflf01"
  router   = "bravo.nicklasfrahm.xyz."
}

module "google_dns_site_deflf02" {
  source = "../modules/google_dns_site"

  zone     = module.google_dns_zone_nicklasfrahm_dev.name
  location = "deflf02"
  router   = "charlie.nicklasfrahm.xyz."
}

resource "google_dns_record_set" "dev_nicklasfrahm_api" {
  managed_zone = module.google_dns_zone_nicklasfrahm_dev.name
  name         = "netlifycms.nicklasfrahm.dev."
  type         = "CNAME"
  ttl          = 600

  rrdatas = ["mike.nicklasfrahm.xyz."]
}
