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
