provider "hetzner-robot" {
  # Credentials come from the environment:
  #   export HETZNERROBOT_USERNAME=... HETZNERROBOT_PASSWORD=...
  # (or set username/password here directly).
}

# 1. (Optional) Upload a new SSH public key to the Robot key store. Only created
#    when var.ssh_public_key is set; otherwise reference an existing key by
#    fingerprint via var.ssh_key_fingerprint.
resource "hetzner-robot_ssh_key" "node" {
  count = var.ssh_public_key != "" ? 1 : 0
  name  = var.ssh_key_name
  data  = var.ssh_public_key
}

locals {
  # authorize the uploaded key's fingerprint and/or an existing one; empties dropped.
  authorized_key_fingerprints = compact([
    var.ssh_public_key != "" ? hetzner-robot_ssh_key.node[0].fingerprint : "",
    var.ssh_key_fingerprint,
  ])
}

# 2. Read-only lookups: what can we order, and what do we already have?
data "hetzner-robot_server_products" "available" {}

data "hetzner-robot_servers" "existing" {}

# 3. Validate a dedicated-server order.
#    test = true  -> NO charge, nothing is provisioned (just validated).
#    test = false -> places a REAL, billable order.
resource "hetzner-robot_server_order" "validate" {
  product_id      = var.product_id
  location        = var.location
  dist            = "Ubuntu 24.04 LTS minimal"
  lang            = "en"
  authorized_keys = local.authorized_key_fingerprints

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
