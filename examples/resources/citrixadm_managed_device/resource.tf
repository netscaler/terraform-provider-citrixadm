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