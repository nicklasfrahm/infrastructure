# Google DNS GitHub Pages

This module configures a domain on Google Cloud DNS to point to a GitHub Pages site.

## Usage

```hcl
module "google_dns_github_pages_example_com" {
  source = "../modules/google_dns_github_pages"

  organization = "example"
  domain       = "example.com"
  zone         = google_managed_dns_zone.example_com.name
}
```
