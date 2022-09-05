data "citrixadm_apigw_proxy" "demo" {
  name = "tf_proxy_new"
}
output "demoId" {
  value = data.citrixadm_apigw_proxy.demo.id
}