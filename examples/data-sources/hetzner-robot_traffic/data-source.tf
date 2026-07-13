# Query monthly traffic for an IP address.
data "hetzner-robot_traffic" "node" {
  ip   = "1.2.3.4"
  type = "month"
  from = "2024-01"
  to   = "2024-12"
}
