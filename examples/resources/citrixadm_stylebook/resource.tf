resource "citrixadm_stylebook" "stylebook1" {
  name      = "basic-lb-config"
  namespace = "com.example.stylebooks"
  version   = "0.1"
  source    = file("./sample_stylebook.yaml")
}
