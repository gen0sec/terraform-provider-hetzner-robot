# Generate a virtual MAC (vMAC) for an IP address.
resource "hetzner-robot_ip_mac" "node" {
  ip = "1.2.3.4"
}
