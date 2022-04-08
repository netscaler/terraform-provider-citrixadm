data "citrixadm_mps_agent" "agent1" {
  name = "10.0.1.91"
}

resource "citrixadm_ns_device_profile" "profile1" {
  name       = "sample_profile"
  username   = "nsroot"
  password   = "notnsroot" #FIXME: make this a secret
  http_port  = "80"
  https_port = "443"
}

resource "citrixadm_managed_device" "device1" {
  ip_address    = "10.0.1.166"
  profile_name  = citrixadm_ns_device_profile.profile1.name
  datacenter_id = data.citrixadm_mps_agent.agent1.datacenter_id
  agent_id      = data.citrixadm_mps_agent.agent1.id
}

# Create Configpack for Existing Device
data "citrixadm_managed_device" "device2" {
  ip_address = "10.0.1.42"
}

resource "citrixadm_stylebook_configpack" "cfgpack1" {
  stylebook {
    name      = "lb"
    namespace = "com.citrix.adc.stylebooks"
    version   = "1.1"
  }
  parameters = {
    lb-appname       = "tf-sample-lb1"
    lb-service-type  = "HTTP"
    lb-virtual-ip    = "4.3.3.4"
    lb-virtual-port  = "80"
    svc-service-type = "HTTP"
  }
  targets {
    instance_id = citrixadm_managed_device.device1.id # new device created in this module
  }
  targets {
    instance_id = data.citrixadm_managed_device.device2.id # existing device
  }
}
