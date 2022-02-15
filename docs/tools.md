# Tools 🧰

The infrastructure repository also contains several tools in the form of self-contained Go binaries. They can be found in the `cmd/` directory of this repository. Their documentation and usage information is contained within this file, but may also be obtained by running `<cmd> -h`.

The API of these tools is subject to change and not stable. Once an API is stable, they will **graduate** and be moved into their own repository. Until then, the tools will follow the versioning **of the their parent `infrastructure` repository**. Hence, a major breaking change **may or may not** affect the API of a tool. You are using experimental tools; here be dragons. 🐲

## Google Cloud DNS Management ☁️ via `dnsadm`

**Status: Beta 🧪**

### Usage

```shell
$ dnsadm -h
A command line interface to automate DNS management tasks. Currently
it only uses Google Cloud DNS, but this may change in the future.

Usage:
  dnsadm [flags]
  dnsadm [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  ddns        Dynamically update DNS records
  help        Help about any command

Flags:
  -a, --auth-file string   set cloud credentials path (default "/etc/secrets/credentials.json")
  -h, --help               display help for command

Use "dnsadm [command] --help" for more information about a command.
```
