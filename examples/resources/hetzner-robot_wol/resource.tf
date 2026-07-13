# Send a Wake-on-LAN packet to a server.
resource "hetzner-robot_wol" "node" {
  server_number = 123456
}
