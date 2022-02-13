# Tools 🧰

The infrastructure repository also contains several tools in the form of self-contained Go binaries. Below you may find usage information for each of them.

## Dynamic Google Cloud DNS ☁️

The `dnsctl` tool can be used to dynamically update the IP address of a Google Cloud DNS record with the public IP address of the host. The tool makes the following assumptions:

- If the top-level domain is `example.com`, the DNS zone is named `example-com`.
- The target DNS record must only contain **A** or **AAAA** records with a **single IP** per record type.

### Usage

```shell
$ dnsctl
A dynamic DNS client to update or create Google Cloud DNS
records. By default it will load the credentials from the
/etc/secrets/gcp.json JSON file. This behaviour can be
overwritten by setting the CREDENTIALS_FILE environment
variable. The service account in question must have the
DNS Admin role.

Usage:
  dnsctl <domain>
```
