data "citrixadm_managed_device" "device1" {
  ip_address = "10.0.1.42"
}
resource "citrixadm_apigw_proxy" "proxy1" {
  proxy_name           = "tf_proxy101"
  host                 = "2.2.2.2"
  port                 = 22
  protocol             = "https"
  deploy               = "false"
  service_fqdns        = ["api.example.com"]
  tls_security_profile = "High Security"
  tls_certkey_objref {
    adm_certkey {
      id = "abcdefgh-1234-5678-ijkl-mnopqrstuvwx"
    }
  }
  target_apigw {
    target_device {
      id           = data.citrixadm_managed_device.device1.id
      display_name = data.citrixadm_managed_device.device1.ip_address
    }
  }
}
