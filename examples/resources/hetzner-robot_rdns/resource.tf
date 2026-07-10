# Set the PTR record for a server IP.
resource "hetzner-robot_rdns" "node" {
  ip  = "1.2.3.4"
  ptr = "node1.k8s.example.com"
}
