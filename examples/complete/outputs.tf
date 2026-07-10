output "uploaded_key_fingerprint" {
  description = "Fingerprint of the SSH key uploaded to the Robot key store."
  value       = hetzner-robot_ssh_key.node.fingerprint
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
