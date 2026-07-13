# Generate a virtual MAC (vMAC) for a subnet.
resource "hetzner-robot_subnet_mac" "net" {
  subnet_ip = "2a01:4f8:1:2::"
}
