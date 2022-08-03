# Google DNS Site

This module creates the corresponding location code DNS records for the given site based on the hostname of the site's router.

Let's assume the sites location code is `uscbf01`. If the hostname of the site's router is `delta.example.com`, the module will create two `CNAME` records for `uscbf01.${zone_dns_name}` and `*.uscbf01.${zone_dns_name}`, which both point to `delta.example.com`.

## Usage

```hcl
module "google_dns_site_uscbf01" {
  source = "../modules/google_dns_site"

  zone     = google_managed_dns_zone.example_com.name
  location = "uscbf01"
  router   = "delta.example.com."
}
```
