data "citrixadm_apigw_deployment" "demo" {
  name = "tf_deploy"
}
output "demoId" {
  value = data.citrixadm_apigw_deployment.demo.id
}
data "citrixadm_apigw_upstream_service" "upstreamdemo" {
  name          = "sdf"
  deployment_id = data.citrixadm_apigw_deployment.demo.id
}
output "demoId1" {
  value = data.citrixadm_apigw_upstream_service.upstreamdemo.id
}