data "citrixadm_apigw_definition" "demo" {
  name = "new_definition"
}
data "citrixadm_managed_device" "device1" {
  ip_address = "10.10.10.10"
}
data "citrixadm_apigw_proxy" "demo_proxy" {
  name = "testing"
}
resource "citrixadm_apigw_deployment" "demo" {
  api_id = data.citrixadm_apigw_definition.demo.id
  apiproxy_ref {
    id = data.citrixadm_apigw_proxy.demo_proxy.id
  }
  name   = "tf_deployment"
  tags   = ["hello"]
  deploy = true
  target_apigw {
    id           = data.citrixadm_managed_device.device1.id
    display_name = data.citrixadm_managed_device.device1.ip_address
  }
  routes {
    name                 = "routing_name_0"
    route_param          = "/user/action"
    route_paramtype      = "resource_path"
    upstreamservice_name = "sdf"
  }
  routes {
    name                 = "route_name2_0"
    route_param          = "/user"
    route_paramtype      = "resource_path"
    upstreamservice_name = "sdf"
  }
  upstreamservices {
    backend_servers {
      ip4_addr = "2.2.2.2"
      port     = 80
    }
    name              = "sdf"
    scheme            = "http"
    service_fqdn      = "www.exam123.com"
    service_fqdn_port = 80
  }
  upstreamservices {
    backend_servers {
      ip4_addr = "2.2.2.2"
      port     = 80
    }
    name              = "sdf1"
    scheme            = "http"
    service_fqdn      = "www.example.com"
    service_fqdn_port = 80
  }
}
