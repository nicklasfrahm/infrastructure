module "nicklasfrahm_dev" {
  source = "../modules/google_dns_zone"

  domain      = "nicklasfrahm.dev."
  description = "My customer-facing domain"
}

module "nicklasfrahm_dev_github_pages" {
  source = "../modules/google_dns_github_pages"

  organization = "nicklasfrahm"
  domain       = "nicklasfrahm.dev."
  zone         = module.nicklasfrahm_dev.name
}
