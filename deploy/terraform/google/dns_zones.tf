module "nicklasfrahm_dev" {
  source = "../modules/google_dns_zone"

  domain      = "nicklasfrahm.dev"
  description = "My customer-facing domain"
}
