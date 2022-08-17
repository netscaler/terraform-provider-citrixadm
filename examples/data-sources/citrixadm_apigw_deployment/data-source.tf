data "citrixadm_apigw_deployment" "demo" {
  name = "deploy"
}
output "demoId" {
  value = data.citrixadm_apigw_deployment.demo.id
}