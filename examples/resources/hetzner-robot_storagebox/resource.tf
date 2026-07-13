# Manage the name and service toggles of an existing Storage Box.
resource "hetzner-robot_storagebox" "backup" {
  storagebox_id         = 123456
  name                  = "k8s-backups"
  ssh                   = true
  samba                 = false
  webdav                = false
  external_reachability = true
  zfs                   = false
}
