variable "ssh_public_key" {
  type        = string
  description = "Your SSH public key (e.g. 'ssh-ed25519 AAAA... you@host'). Uploaded to the Hetzner Robot key store."
}

variable "ssh_key_name" {
  type        = string
  description = "Name for the uploaded key in the Robot key store."
  default     = "k8s-node-key"
}

variable "product_id" {
  type        = string
  description = "Product ID to validate ordering for (see the available_product_ids output)."
  default     = "EX44"
}

variable "location" {
  type        = string
  description = "Datacenter location."
  default     = "FSN1"
}
