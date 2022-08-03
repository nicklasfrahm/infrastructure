# Google DNS Zone

This module creates a Google Cloud DNS zone. It enforces a very strict DNSSEC security level using NSEC3 and ECDSA.

## Usage

```hcl
module "example_com" {
  source = "../google_dns_zone"

  # This name must not end with a dot.
  name        = "example.com"
  description = "Example DNS zone"
}
```

## Outputs

| Name | Description             |
| ---- | ----------------------- |
| `id` | The ID of the DNS zone. |
