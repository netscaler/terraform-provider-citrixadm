data "citrixadm_apigw_deployment" "demo" {
  name = "deploy"
}
resource "citrixadm_apigw_route" "route1" {
  name                 = "tf_route1"
  route_param          = "/user"
  route_paramtype      = "resource_path"
  upstreamservice_name = "upstream_service10"
  deployment_id        = data.citrixadm_apigw_deployment.demo.id
}