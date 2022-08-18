data "citrixadm_apigw_deployment" "demo" {
  name = "deploy"
}
resource "citrixadm_apigw_upstream_service" "service1" {
  name              = "tf_upstream_service"
  scheme            = "http"
  service_fqdn      = "www.example.com"
  service_fqdn_port = 80
  deployment_id     = data.citrixadm_apigw_deployment.demo.id
  backend_servers {
    ip4_addr = "2.1.2.4"
    port     = 80
  }
}
