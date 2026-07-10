# Provide EITHER ssh_public_key (to upload a new key) OR ssh_key_fingerprint
# (to reference a key already in your Robot key store). Neither is required for a
# test order, but you need one to authorize SSH access on a real server.

variable "ssh_public_key" {
  type        = string
  description = "SSH PUBLIC key to upload as a new Robot key (e.g. 'ssh-ed25519 AAAA... you@host'). Leave empty to use an existing key via ssh_key_fingerprint."
  default     = ""
}

variable "ssh_key_fingerprint" {
  type        = string
  description = "Fingerprint of an EXISTING key in the Robot key store to authorize (alternative to uploading ssh_public_key). Find it in the Robot UI or via GET /key."
  default     = ""
}

variable "ssh_key_name" {
  type        = string
  description = "Name for the uploaded key (only used when ssh_public_key is set)."
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
