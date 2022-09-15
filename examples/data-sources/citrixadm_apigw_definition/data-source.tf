data "citrixadm_apigw_definition" "demo" {
  name = "new_at"
}
output "demoId" {
  value = data.citrixadm_apigw_definition.demo.id
}