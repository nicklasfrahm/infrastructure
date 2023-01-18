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

module "google_dns_github_pages_kubestack_nicklasfrahm_dev" {
  source = "../modules/google_dns_github_pages"

  organization = "nicklasfrahm"
  domain       = "kubestack.nicklasfrahm.dev."
  zone         = module.google_dns_zone_nicklasfrahm_dev.name
}

module "google_dns_site_dktil01" {
  source = "../modules/google_dns_site"

  zone     = module.google_dns_zone_nicklasfrahm_dev.name
  location = "dktil01"
  router   = "alfa.nicklasfrahm.xyz."
}

module "google_dns_site_dksjb00" {
  source = "../modules/google_dns_site"

  zone     = module.google_dns_zone_nicklasfrahm_dev.name
  location = "dksjb00"
  router   = "delta.nicklasfrahm.xyz."
}

module "google_dns_site_dksjb01" {
  source = "../modules/google_dns_site"

  zone     = module.google_dns_zone_nicklasfrahm_dev.name
  location = "dksjb01"
  router   = "delta.nicklasfrahm.xyz."
}

module "google_dns_site_dksjb02" {
  source = "../modules/google_dns_site"

  zone     = module.google_dns_zone_nicklasfrahm_dev.name
  location = "dksjb02"
  router   = "delta.nicklasfrahm.xyz."
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

  rrdatas = ["delta.nicklasfrahm.xyz."]
}

resource "google_dns_record_set" "dev_nicklasfrahm_mc_survival" {
  managed_zone = module.google_dns_zone_nicklasfrahm_dev.name
  name         = "mc-survival.nicklasfrahm.dev."
  type         = "CNAME"
  ttl          = 600

  rrdatas = ["delta.nicklasfrahm.xyz."]
}
