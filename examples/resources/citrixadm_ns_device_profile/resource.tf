resource "citrixadm_ns_device_profile" "profile1" {
  name       = "tf_ns_profile"
  username   = "nsroot"
  password   = "verysecretpassword"
  http_port  = "80"
  https_port = "443"
}
