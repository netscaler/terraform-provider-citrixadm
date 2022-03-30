data "citrixadm_mps_agent" "agent1" {
  name = "10.0.1.91"
}

resource "citrixadm_ns_device_profile" "profile1" {
  name     = "tf_test_profile2"
  username = "nsroot"
  password = "notnsroot"
}

resource "citrixadm_managed_device" "device1" {
  ip_address    = "10.0.1.145"
  profile_name  = citrixadm_ns_device_profile.profile1.name
  datacenter_id = data.citrixadm_mps_agent.agent1.datacenter_id
  agent_id      = data.citrixadm_mps_agent.agent1.id
  #   description   = "VPX managed device description" # FIXME: API problem. API returns always an empty string
  entity_tag { # FIXME: API problem. API stores the tags' keys and values in lowercase and if we pass in mixedcase, it's a new update everytime
    prop_key   = "project"
    prop_value = "adms"

  }
  entity_tag {
    prop_key   = "environment"
    prop_value = "staging"
  }

  # license_edition= "Platinum"
  # plt_bw_config= 30

}

resource "citrixadm_managed_device" "device2" {
  ip_address    = "10.0.1.215"
  profile_name  = citrixadm_ns_device_profile.profile1.name
  datacenter_id = data.citrixadm_mps_agent.agent1.datacenter_id
  agent_id      = data.citrixadm_mps_agent.agent1.id

  # license_edition= "Platinum"
  # plt_bw_config= 500

}

resource "citrixadm_managed_device" "device3" {
  ip_address    = "10.0.1.233"
  profile_name  = citrixadm_ns_device_profile.profile1.name
  datacenter_id = data.citrixadm_mps_agent.agent1.datacenter_id
  agent_id      = data.citrixadm_mps_agent.agent1.id
}
resource "citrixadm_managed_device" "device4" {
  ip_address    = "10.0.1.144"
  profile_name  = citrixadm_ns_device_profile.profile1.name
  datacenter_id = data.citrixadm_mps_agent.agent1.datacenter_id
  agent_id      = data.citrixadm_mps_agent.agent1.id
}