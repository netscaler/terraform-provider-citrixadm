data "citrixadm_cert_store" "demo" {
  name = "tf_certname"
}
output "demoId" {
  value = data.citrixadm_cert_store.demo.id
}