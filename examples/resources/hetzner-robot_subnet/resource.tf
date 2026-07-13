# Configure traffic warnings for a subnet.
resource "hetzner-robot_subnet" "net" {
  subnet_ip        = "2a01:4f8:1:2::"
  traffic_warnings = true
  traffic_monthly  = 50
}
