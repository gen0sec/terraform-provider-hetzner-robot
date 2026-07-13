# Create a snapshot of a Storage Box.
resource "hetzner-robot_storagebox_snapshot" "nightly" {
  storagebox_id = 123456
}
