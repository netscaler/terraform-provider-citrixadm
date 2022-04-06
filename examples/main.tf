resource "citrixadm_stylebook_configpack" "cfgpack1" {
  stylebook {
    name      = "lb"
    namespace = "com.citrix.adc.stylebooks"
    version   = "1.1"
  }
  # parameters = {
  parameters {
    lb-appname       = "tf-sample-lb3"
    lb-service-type  = "HTTP"
    lb-virtual-ip    = "4.3.3.4"
    lb-virtual-port  = 80
    svc-service-type = "HTTP"
  }
  targets {
    instance_id = data.citrixadm_managed_device.device1.id
  }
}

data "citrixadm_managed_device" "device1" {
  ip_address = "10.0.1.42"
}

# data "citrixadm_mps_agent" "agent1" {
#   name = "10.0.1.91"
# }

# resource "citrixadm_ns_device_profile" "profile1" {
#   name     = "tf_test_profile2"
#   username = "nsroot"
#   password = "notnsroot"
# }

# resource "citrixadm_managed_device" "device1" {
#   ip_address    = "10.0.1.145"
#   profile_name  = citrixadm_ns_device_profile.profile1.name
#   datacenter_id = data.citrixadm_mps_agent.agent1.datacenter_id
#   agent_id      = data.citrixadm_mps_agent.agent1.id
#   #   description   = "VPX managed device description" # FIXME: API problem. API returns always an empty string
#   entity_tag { # FIXME: API problem. API stores the tags' keys and values in lowercase and if we pass in mixedcase, it's a new update everytime
#     prop_key   = "project"
#     prop_value = "adms"

#   }
#   entity_tag {
#     prop_key   = "environment"
#     prop_value = "staging"
#   }

#   # license_edition= "Platinum"
#   # plt_bw_config= 30

# }

# resource "citrixadm_managed_device" "device2" {
#   ip_address    = "10.0.1.215"
#   profile_name  = citrixadm_ns_device_profile.profile1.name
#   datacenter_id = data.citrixadm_mps_agent.agent1.datacenter_id
#   agent_id      = data.citrixadm_mps_agent.agent1.id

#   # license_edition= "Platinum"
#   # plt_bw_config= 500

# }

# resource "citrixadm_managed_device" "device3" {
#   ip_address    = "10.0.1.233"t
#   profile_name  = citrixadm_ns_device_profile.profile1.name
#   datacenter_id = data.citrixadm_mps_agent.agent1.datacenter_id
#   agent_id      = data.citrixadm_mps_agen.agent1.id
# }
# resource "citrixadm_managed_device" "device4" {
#   ip_address    = "10.0.1.144"
#   profile_name  = citrixadm_ns_device_profile.profile1.name
#   datacenter_id = data.citrixadm_mps_agent.agent1.datacenter_id
#   agent_id      = data.citrixadm_mps_agent.agent1.id
# }

# resource "citrixadm_managed_device" "device5" {
#   ip_address    = "10.0.1.48"
#   profile_name  = "nsroot_notnsroot_profile"
#   datacenter_id = data.citrixadm_mps_agent.agent1.datacenter_id
#   agent_id      = data.citrixadm_mps_agent.agent1.id

# }

# resource "citrixadm_managed_device_allocate_license" "lic1" {
#   managed_device_id = citrixadm_managed_device.device5.id
#   license_edition= "Platinum"
#   plt_bw_config= 1700
# }
