# Set an install/rescue boot profile, then reset the server so it takes effect.
# `triggers` ties the reset to the boot config so changing the profile re-runs it.
resource "hetzner-robot_boot" "example" {
  server_id        = 123456
  active_profile   = "linux"
  architecture     = "64"
  operating_system = "Ubuntu 24.04 LTS minimal"
  language         = "en"
  authorized_keys  = ["aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99"]
}

resource "hetzner-robot_reset" "example" {
  server_id  = hetzner-robot_boot.example.server_id
  reset_type = "hw" # hw | sw | power | man

  triggers = {
    boot_profile = hetzner-robot_boot.example.active_profile
  }
}
