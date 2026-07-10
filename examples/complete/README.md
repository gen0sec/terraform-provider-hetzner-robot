# Complete example

A safe, runnable walkthrough of the provider:

1. **Uploads your SSH public key** to the Robot key store (`hetzner-robot_ssh_key`) — real, but trivially reversible.
2. **Lists** available server products and existing servers (read-only data sources).
3. **Validates** a dedicated-server order with `test = true` — **no charge, nothing provisioned**.

The `boot` → `reset` install flow is included but commented out (it reinstalls a real machine).

## What you add
Only your **SSH public key**. Copy the example tfvars and paste it in:

```bash
cp terraform.tfvars.example terraform.tfvars
# edit terraform.tfvars: set ssh_public_key = "ssh-ed25519 AAAA... you@host"
```

Credentials come from the environment:

```bash
export HETZNERROBOT_USERNAME='#ws+xxxxxxxx'   # Robot webservice user
export HETZNERROBOT_PASSWORD='...'
```

## Running against the local build

The provider isn't published to a registry, so build it and point Terraform at the
binary with a dev override (no `terraform init` needed):

```bash
# 1. build the provider (from the repo root)
go build -o terraform-provider-hetzner-robot .

# 2. add a dev override to ~/.terraformrc (adjust the path to this repo)
cat > ~/.terraformrc <<EOF
provider_installation {
  dev_overrides {
    "gen0sec/hetzner-robot" = "/absolute/path/to/terraform-provider-hetzner-robot"
  }
  direct {}
}
EOF

# 3. plan / apply (dev_overrides skips init)
cd examples/complete
terraform plan
terraform apply
```

## Going further
- To place a **real** order, set `test = false` on `hetzner-robot_server_order` (billable, provisions hardware; `server_number`/`server_ip` populate once ready).
- Then uncomment the `boot` + `reset` resources (fill in the `server_number`) to install Ubuntu with your key. From there the node joins your k3s/k8s cluster with your normal bootstrap (e.g. Tailscale + k3s agent).
