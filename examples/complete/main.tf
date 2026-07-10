provider "hetzner-robot" {
  # Credentials come from the environment:
  #   export HETZNERROBOT_USERNAME=... HETZNERROBOT_PASSWORD=...
  # (or set username/password here directly).
}

# 1a. (Optional) Reference an EXISTING key by name — looked up via GET /key.
data "hetzner-robot_ssh_key" "existing" {
  count = var.existing_ssh_key_name != "" ? 1 : 0
  name  = var.existing_ssh_key_name
}

# 1b. (Optional) Upload a NEW key. Only created when var.ssh_public_key is set.
resource "hetzner-robot_ssh_key" "node" {
  count = var.ssh_public_key != "" ? 1 : 0
  name  = var.upload_key_name
  data  = var.ssh_public_key
}

locals {
  # authorize the existing and/or uploaded key fingerprints; empties dropped.
  authorized_key_fingerprints = compact(concat(
    [for k in data.hetzner-robot_ssh_key.existing : k.fingerprint],
    [for k in hetzner-robot_ssh_key.node : k.fingerprint],
  ))
}

# 2. Read-only lookups: what can we order, and what do we already have?
data "hetzner-robot_server_products" "available" {}

data "hetzner-robot_servers" "existing" {}

# Orderable install options (dist/lang/arch) for the product we're ordering —
# use these to fill dist/lang on a real order.
data "hetzner-robot_server_product" "selected" {
  id = var.product_id
}

# 3. Validate a dedicated-server order.
#    test = true  -> NO charge, nothing is provisioned (just validated).
#    test = false -> places a REAL, billable order.
resource "hetzner-robot_server_order" "validate" {
  product_id      = var.product_id
  location        = var.location
  authorized_keys = local.authorized_key_fingerprints

  # To also pre-install an OS, add `dist` (and `lang`) using values from that
  # product's orderable options:
  #   curl -su "$HETZNERROBOT_USERNAME:$HETZNERROBOT_PASSWORD" \
  #     https://robot-ws.your-server.de/order/server/product/AX42-1 | jq '.product.dist, .product.lang'

  test = true
}

# ----------------------------------------------------------------------------
# Once you have a real server (ordered above with test=false, or an existing
# one), this is the boot -> reset flow that installs Ubuntu and gets you SSH.
# Left commented because it reboots/reinstalls a real machine.
# ----------------------------------------------------------------------------
# resource "hetzner-robot_boot" "node" {
#   server_number    = 123456
#   active_profile   = "linux"
#   operating_system = "Ubuntu 24.04 LTS minimal"
#   language         = "en"
#   authorized_keys  = local.authorized_key_fingerprints
# }
#
# resource "hetzner-robot_reset" "node" {
#   server_number = hetzner-robot_boot.node.server_number
#   reset_type    = "hw"
#   triggers = {
#     boot_profile = hetzner-robot_boot.node.active_profile
#   }
# }
