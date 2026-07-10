# Provide EITHER existing_ssh_key_name (reference a key already in your Robot key
# store, looked up by name) OR ssh_public_key (upload a new key). Both optional;
# a test order works with neither.

variable "existing_ssh_key_name" {
  type        = string
  description = "Name of an EXISTING key in your Robot key store to authorize (resolved to a fingerprint via GET /key)."
  default     = ""
}

variable "ssh_public_key" {
  type        = string
  description = "SSH PUBLIC key to upload as a NEW Robot key (e.g. 'ssh-ed25519 AAAA... you@host')."
  default     = ""
}

variable "upload_key_name" {
  type        = string
  description = "Name for the uploaded key (only used when ssh_public_key is set)."
  default     = "k8s-node-key"
}

variable "product_id" {
  type        = string
  description = "Product ID to validate ordering for (see the available_product_ids output, e.g. AX42-1, EX63-1)."
  default     = "AX42-1"
}

variable "location" {
  type        = string
  description = "Datacenter location the product supports (e.g. FSN1 or HEL1)."
  default     = "FSN1"
}

# (Optional) set a reverse-DNS (PTR) record. Opt-in: both must be set, otherwise
# no rDNS change is made — so this never touches an existing server by default.
variable "rdns_ip" {
  type        = string
  description = "IP to set a PTR record for. Leave empty to skip."
  default     = ""
}

variable "rdns_ptr" {
  type        = string
  description = "PTR hostname to set for rdns_ip."
  default     = ""
}
