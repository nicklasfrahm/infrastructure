module "google_vm_juliett" {
  source = "../modules/google_vm"

  hostname        = "juliett"
  fqdn            = "nicklasfrahm.xyz"
  github_username = "nicklasfrahm"
}

output "juliett_ip" {
  value = module.google_vm_juliett.ip
}
