# terraform-provider-hetzner-robot

A Terraform provider for the [Hetzner Robot](https://robot.hetzner.com/doc/webservice/en.html)
webservice — manage **dedicated servers** (boot config, resets, firewall, SSH keys,
vSwitches, reverse DNS) and **order** new servers, as code.

Full reference: [docs/](docs/index.md).

## Features

**Resources**

| Resource | Purpose |
|---|---|
| `hetzner-robot_server_order` | Order a dedicated server (`test = true` by default → validate without charge) |
| `hetzner-robot_boot` | Set the boot/install profile — `linux`, `rescue`, `vnc`, `windows` (deactivates on destroy) |
| `hetzner-robot_reset` | Trigger a hardware/software reset (apply a boot profile) |
| `hetzner-robot_ssh_key` | Manage a key in the Robot key store |
| `hetzner-robot_firewall` | Manage a server's firewall rules |
| `hetzner-robot_vswitch` | Manage a vSwitch (private L2 network) |
| `hetzner-robot_rdns` | Manage a reverse-DNS (PTR) record |
| `hetzner-robot_ip` / `hetzner-robot_subnet` | Manage traffic-warning limits for an IP / subnet |
| `hetzner-robot_ip_mac` / `hetzner-robot_subnet_mac` | Manage a virtual MAC (vMAC) for an IP / subnet |
| `hetzner-robot_wol` | Send a Wake-on-LAN packet to a server |
| `hetzner-robot_storagebox` | Manage an existing Storage Box's name and service toggles |
| `hetzner-robot_storagebox_snapshot` | Create a Storage Box snapshot |
| `hetzner-robot_storagebox_subaccount` | Manage a Storage Box subaccount |

**Data sources**

| Data source | Purpose |
|---|---|
| `hetzner-robot_server` / `hetzner-robot_servers` | Look up a server / list all servers |
| `hetzner-robot_server_products` / `hetzner-robot_server_product` | List orderable products / one product's orderable `dist`/`lang`/`arch` |
| `hetzner-robot_ssh_key` | Look up a key **by name or fingerprint** |
| `hetzner-robot_boot` | Read current boot config |
| `hetzner-robot_vswitch` | Read a vSwitch |
| `hetzner-robot_rdns` | Read a PTR record |
| `hetzner-robot_ip` / `hetzner-robot_subnet` | Read an IP / subnet's details |
| `hetzner-robot_traffic` | Query traffic statistics for an IP over a period |
| `hetzner-robot_storagebox` / `hetzner-robot_storageboxes` | Read a Storage Box / list all Storage Boxes |

## Usage

```hcl
provider "hetzner-robot" {
  # credentials from HETZNERROBOT_USERNAME / HETZNERROBOT_PASSWORD (or set here)
}

data "hetzner-robot_ssh_key" "node" {
  name = "k8s-node-key"
}

resource "hetzner-robot_server_order" "worker" {
  product_id      = "AX42-1"
  location        = "FSN1"
  dist            = "Ubuntu 24.04 LTS base"
  authorized_keys = [data.hetzner-robot_ssh_key.node.fingerprint]
  test            = false # ⚠ billable — provisions real hardware
}
```

A complete, runnable walkthrough (order → install → rDNS) is in
[`examples/complete/`](examples/complete).

### Provider configuration

| Argument | Env var | Notes |
|---|---|---|
| `username` | `HETZNERROBOT_USERNAME` | Robot webservice user (`#ws+…`) |
| `password` | `HETZNERROBOT_PASSWORD` | |
| `url` | `HETZNERROBOT_URL` | defaults to `https://robot-ws.your-server.de` |

## Local development

The provider isn't published to a registry, so build it and use a **dev override**
(no `terraform init`):

```bash
go build -o terraform-provider-hetzner-robot .

cat > ~/.terraformrc <<EOF
provider_installation {
  dev_overrides {
    "gen0sec/hetzner-robot" = "/absolute/path/to/this/repo"
  }
  direct {}
}
EOF

cd examples/complete && terraform plan   # NOT init
```

## Build & release

```bash
# local snapshot
goreleaser release --snapshot --clean
```

Releases build automatically via GitHub Actions + GoReleaser on each new `v*` tag.

## Testing

```bash
go test -race -cover ./...   # unit tests (mock the Robot API; no credentials needed)
```

Acceptance tests hit the live API and are gated behind `TF_ACC`:

```bash
TF_ACC=1 HETZNERROBOT_USERNAME=… HETZNERROBOT_PASSWORD=… go test ./hetznerrobot/ -run TestAcc -v
```

## Fork lineage

`mwudka/terraform-provider-hetznerrobot` → `SLoeuillet/terraform-provider-hetznerrobot`
→ `Peters-IT/terraform-provider-hetzner-robot` → `strng-solutions/terraform-provider-hetzner-robot`
→ **this fork** (`gen0sec`), which consolidates those and adds `reset`, `server_order`,
`rdns`, `server_product`, SSH-key-by-name lookup, unit tests, and updated tooling.

This software comes without any guarantee of functionality. PRs welcome.
