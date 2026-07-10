output "authorized_key_fingerprints" {
  description = "Fingerprints authorized on the order (uploaded and/or existing)."
  value       = local.authorized_key_fingerprints
}

output "available_product_ids" {
  description = "Product IDs you can order."
  value       = [for p in data.hetzner-robot_server_products.available.products : p.id]
}

output "existing_servers" {
  description = "Servers already in the account."
  value = [for s in data.hetzner-robot_servers.existing.servers : {
    number = s.server_number
    name   = s.server_name
  }]
}

output "order_validation_status" {
  description = "Status returned by the test order (test = true)."
  value       = hetzner-robot_server_order.validate.status
}

output "product_install_options" {
  description = "Orderable dist/lang/arch for the selected product — use these for a real order's dist/lang."
  value = {
    dist = data.hetzner-robot_server_product.selected.dist
    lang = data.hetzner-robot_server_product.selected.lang
    arch = data.hetzner-robot_server_product.selected.arch
  }
}

output "rdns" {
  description = "The PTR record set for rdns_ip (null when not configured)."
  value = length(hetzner-robot_rdns.node) > 0 ? {
    ip  = hetzner-robot_rdns.node[0].ip
    ptr = hetzner-robot_rdns.node[0].ptr
  } : null
}
