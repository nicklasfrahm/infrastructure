# Infrastructure

A repository to automate the configuration of my local infrastructure.

[![Pulumi](https://github.com/nicklasfrahm/infrastructure/actions/workflows/pulumi.yml/badge.svg)](https://github.com/nicklasfrahm/infrastructure/actions/workflows/pulumi.yml)

My local infrastructure is my homelab, which I primarily use to learn about state-of-the-art technologies as well as to run applications for potential startup incubations.

As part of this repository you may find interesting learnings and insights that I came across when developing my homelab and its automation. Automating my infrastructure allows me to deploy workloads more efficiently and debug issues more quickly.

## History 📜

I have been running a homelab since 2014, where it started with a single [Cubietech Cubietruck][website-cubietruck]. It served as a NAS, LEMP server, and was a great introduction for me into Linux. Over time the setup grew to include other single board computers, such as the [Hardkernel Odroid XU4][website-odroidxu4].

Since 2018 I started to make significant changes by automating the configuration of my router and my manually provisioned servers via [Ansible][website-ansible]. Using it allowed me to learn the _classic way_ of Baremetal provisioning. It provided a solid tool to get me started with the deployment and operation of [Kubernetes][website-kubernetes] clusters.

Now, I am in the process of upgrading my infrastructure to become cloud-native and fully automated. The goal is to provide a RESTful API to provision baremetal machines from scratch, configure networking, and bootstrap and autoscale Kubernetes clusters.

## Networking🔌

Below, I describe the network setup of my Homelab. Some parts of it are already implemented, while others are still being conceptualized.

### DNS 🌍

**Status:** 🟢 Operational

I register my domains with [Namecheap][website-namecheap], because they often provide good discounts on variety of domain names. Below, you may find a list of limitations that I ran into, which is the reason why I decided to only use them as a registrar, but not as a DNS provider.

- **Limited usability in CI**  
  Their API requires a manually administrated allow-list of trusted IPs. The allow list they offer does not allow CIDR notation IP ranges, e.g. `0.0.0.0/24`. The issue here is that [GitHub Actions'][website-github-actions] hosted runners do not have a static set of IPs. As I want to administrate my entire infrastructure via GitHub actions and pipelines, a self-host runner is not an option for me.

- **Slow API development**  
  The [Namecheap API][website-namecheap-api] is well-documented, but is not RESTful and uses XML. This setup is rather old-fashioned and makes it rather uncomfortable to use. Their customers have requested a more modern API surface, but so far the development of the API focuses on pure maintenance. Having a lightweight, easy-to-use API is critical when relying on it to build your products.

Therefore I decided to use [Google Cloud DNS][website-gcp-dns] for the automated administration of my DNS records. The deployment of my DNS infrastructure is automated using [GitHub Actions][website-github-actions] and [Pulumi][website-pulumi], because it allows me to write my _Infrastructure as Code_ using my favourite programming language [Go][website-go].

For more information you may want to check out the following files and directories:

- `main.go`: The entry point to my Pulumi program as a general overview.
- `pkg/dns/`: The configuration of DNSSEC and my DNS records.

### Routing 🔀

**Status**: 🟢 Operational

While I am migrating my current setup, I am using an [Ubiquity Edgerouter-X][website-erx] (ER-X) with port-forwarding towards a single baremetal [k3s cluster][website-k3s] on a seperate machine. The ER-X is fine, but not great for automation and dynamic configuration. I would like to run a Kubernetes cluster on my router which is not possible on the ER-X.

I am currently running two subnets, **userspace** (`192.168.0.1/24`) and **homelab** (`172.16.0.1/22`), where DHCP is only enabled in the **userspace** subnet. This setup allows me to maintain network connectivity for **userspace** traffic while performing testing or configuration changes in my **homelab**. The configuration is currently hardcoded to specific ports on my router. In the future I am planning to orchestrate software defined networking via VLANs.

**Status (v2):** 💡 Planning

**Keywords:** BGP, Kubernetes Operator, [Kube-Router][website-kube-router], VLAN, SDN  
**Motivation:** Automation, ease-of-operation, client source IPs

### Load balancing 🔁

**Status:** 🔵 Limited availability

Currently, I am only running [Traefik][website-traefik] as an [Ingress Controller][docs-ingress], which does load balancing and TLS termination for my pods in Kubernetes. In the future I would like to have automatic load balancing for my Kubernetes Control Planes and `LoadBalancer` type services.

**Status (v2):** 💡 Planning

**Keywords:** BGP, [HAProxy][website-haproxy], [Gateway API][website-gateway-api]  
**Motivation:** Automation, ease-of-operation, client source IPs

### VLANs 🛡️

**Status:** 🟡 Rollout

The table below describes the set of manually assigned VLANs. The configuration of the DHCPv4 server is based on the following resources:

- [Ubuntu community article about `isc-dhcp-server`][article-ubuntu-dhcp]
- [How to make a simple router on a Ubuntu server][medium-ubuntu-router]

| VLAN ID | Gateway CIDR       | DHCP | Name       | Description                                                        |
| ------- | ------------------ | ---- | ---------- | ------------------------------------------------------------------ |
| 1       | `none`             | No   | default    | A network without any gateway to isolate unassigned hosts.         |
| 4000    | `172.16.0.1/22`    | Yes  | management | A network for the configuration and management of network devices. |
| 4010    | `172.16.4.1/22`    | No   | homelab    | A network for experimental deployments.                            |
| 4020    | `192.168.255.1/24` | Yes  | userspace  | A home network for WiFi and other domestic traffic.                |

### Firewall 🔥

**Status:** 🟡 Rollout

I chose to set up `firewalld` and use `firewall-cmd` to manage the firewall. The table below shows an overview of the defined zones and their associated services. For the implementation the following resources were used:

- [Digital Ocean tutorial about `firewall-cmd`][tutorial-firewall-cmd]
- [RHEL documentation explaining usage of `firewalld`][documentation-rhel-firewalld]
- [Article about IP masquerade with `firewalld`][website-server-world-masquerade]

_TODO: Update zones, masquerading and services._

| Zone     | Description                          | Services                |
| -------- | ------------------------------------ | ----------------------- |
| external | The connection towards the internet. | `ssh`, `kube-apiserver` |

## Limitations❗

Below you may find a list of limitations with my current infrastructure setup.

- **Manual NS and DS records**  
  Due to the limitations of the [Namecheap API][website-namecheap-api], I decided to completely administrate my NS and DS records by hand. If I can't automate it for one registrar, I will not automate it for any of them. This is subject to change based on the amount of domains and the frequency of changes. Usually however NS records rarely change.

- **Manual administration of VLANs**  
  Because my network contains a variety of network devices with different management protocols, there is a currently little value in automating the management of VLANs. If this should prove to be a valuable investment at some point, the links below include information on how to interface with some of my hardware via protocols such as the _Easy Smart Configuration Protocol (ESCP)_.
  - [How I can gain control of your TP-LINK home switch](https://www.pentestpartners.com/security-blog/how-i-can-gain-control-of-your-tp-link-home-switch/)
  - [Information disclosure vulnerability in TP-Link Easy Smart switches](https://www.chrisdcmoore.co.uk/post/tplink-easy-smart-switch-vulnerabilities/)
  - [Ansible Collection - rgl.tp_link_easy_smart_switch](https://github.com/rgl/ansible-collection-tp-link-easy-smart-switch)

## License 📄

This project is licensed under the terms of the [MIT license][file-license].

[website-cubietruck]: https://linux-sunxi.org/Cubietech_Cubietruck
[website-namecheap]: https://namecheap.com
[website-namecheap-api]: https://www.namecheap.com/support/api/intro/
[website-github-actions]: https://github.com/features/actions
[website-gcp-dns]: https://cloud.google.com/dns
[website-pulumi]: https://www.pulumi.com/
[website-go]: https://golang.org
[website-odroidxu4]: https://wiki.odroid.com/odroid-xu4/odroid-xu4
[website-ansible]: https://www.ansible.com/
[website-kubernetes]: https://kubernetes.io/
[website-erx]: https://www.ui.com/edgemax/edgerouter-x/
[website-k3s]: https://k3s.io
[docs-ingress]: https://kubernetes.io/docs/concepts/services-networking/ingress/
[website-traefik]: https://traefik.io/traefik/
[website-kube-router]: https://www.kube-router.io/
[website-haproxy]: http://www.haproxy.org/
[website-gateway-api]: https://gateway-api.sigs.k8s.io/
[file-license]: ./LICENSE.md
[tutorial-firewall-cmd]: https://www.digitalocean.com/community/tutorials/how-to-set-up-a-firewall-using-firewalld-on-centos-8
[article-ubuntu-dhcp]: https://help.ubuntu.com/community/isc-dhcp-server
[medium-ubuntu-router]: https://medium.com/@exesse/how-to-make-a-simple-router-gateway-from-ubuntu-server-18-04-lts-fd40b7bfec9
[documentation-rhel-firewalld]: https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/8/html/configuring_and_managing_networking/using-and-configuring-firewalld_configuring-and-managing-networking
[website-server-world-masquerade]: https://www.server-world.info/en/note?os=CentOS_7&p=firewalld&f=2
