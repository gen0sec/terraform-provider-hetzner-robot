# Create a subaccount on a Storage Box.
resource "hetzner-robot_storagebox_subaccount" "ci" {
  storagebox_id = 123456
  homedirectory = "/ci-backups"
  ssh           = true
  readonly      = false
  comment       = "CI backup writer"
}
