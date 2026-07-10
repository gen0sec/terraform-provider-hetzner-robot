# Browse available products, then order one. `test = true` (the default) only
# validates the order — flip it to false to place a REAL, billable order.
data "hetzner-robot_server_products" "available" {}

resource "hetzner-robot_server_order" "worker" {
  product_id      = "EX44"
  location        = "FSN1"
  dist            = "Ubuntu 24.04 LTS minimal"
  lang            = "en"
  authorized_keys = ["aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99"]

  test = true # set to false to place a real, billable order
}

output "ordered_server_number" {
  value = hetzner-robot_server_order.worker.server_number
}
