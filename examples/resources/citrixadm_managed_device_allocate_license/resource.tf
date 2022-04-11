# Apply License for new device
data "citrixadm_mps_agent" "agent1" {
  name = "10.0.1.91"
}

resource "citrixadm_ns_device_profile" "profile1" {
  name       = "sample_profile"
  username   = "nsroot"
  password   = "verysecretpassword"
  http_port  = "80"
  https_port = "443"
}

resource "citrixadm_managed_device" "device1" {
  ip_address    = "10.0.1.166"
  profile_name  = citrixadm_ns_device_profile.profile1.name
  datacenter_id = data.citrixadm_mps_agent.agent1.datacenter_id
  agent_id      = data.citrixadm_mps_agent.agent1.id
}

resource "citrixadm_managed_device_allocate_license" "lic1" {
  managed_device_id = citrixadm_managed_device.device1.id
  license_edition   = "Platinum"
  plt_bw_config     = 700 # in Mbps
}




# Apply License for Existing Device
data "citrixadm_managed_device" "device2" {
  ip_address = "10.0.1.42"
}

resource "citrixadm_managed_device_allocate_license" "lic2" {
  managed_device_id = data.citrixadm_managed_device.device2.id
  license_edition   = "Platinum"
  plt_bw_config     = 600 # in Mbps
}
