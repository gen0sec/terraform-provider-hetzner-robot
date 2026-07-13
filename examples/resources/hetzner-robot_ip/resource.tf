# Configure traffic warnings for a single IP address.
resource "hetzner-robot_ip" "node" {
  ip               = "1.2.3.4"
  traffic_warnings = true
  traffic_hourly   = 200
  traffic_daily    = 2000
  traffic_monthly  = 20
}
