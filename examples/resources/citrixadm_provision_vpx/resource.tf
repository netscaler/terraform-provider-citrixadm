resource "citrixadm_provision_vpx" "instance1" {
  name                    = "tf-instance4"
  provisioning_profile_id = citrixadm_provisioning_profile.profile1.id
}

data "citrixadm_mps_agent" "agent1" {
  name = "10.222.74.161"
}

data "citrixadm_config_job_template" "template1" {
  name = "ShowConfiguration"
}

resource "citrixadm_provisioning_profile" "profile1" {
  name          = "tf-instance4"
  instance_type = "NetScaler"
  site_id       = data.citrixadm_mps_agent.agent1.datacenter_id

  instance_capacity_details {
    config_job_templates = [
      data.citrixadm_config_job_template.template1.id
    ]
  }

  mas_registration_details {
    mas_agent_id = data.citrixadm_mps_agent.agent1.id
  }

  platform_type      = "sdx"
  deployment_details = <<EOF
  {
			"sdx": {
				"nitro": {
					"burst_priority": "0",
					"config_type": 0,
					"gateway": "10.222.74.129",
					"if_0_1": true,
					"if_0_2": true,
					"image_name": "NSVPX-XEN-13.1-17.42_nc_64.xva",
					"ip_address": "10.222.74.176",
					"ipv4_address": "10.222.74.176",
					"l2_enabled": "false",
					"license": "Standard",
					"max_burst_throughput": "0",
					"name": "tf-instance4",
					"netmask": "255.255.255.192",
					"network_interfaces": [
						{
							"device_channel_name": "",
							"mac_address": "",
							"mac_mode": "default",
							"port_name": "LA/2",
							"receiveuntagged": true,
							"vlan_whitelist_array": [
								"110"
							]
						}
					],
					"nexthop": "",
					"nsvlan_id": "",
					"nsvlan_interfaces": [],
					"nsvlan_tagged": "false",
					"number_of_acu": 0,
					"number_of_cores": "0",
					"number_of_scu": "0",
					"pps": "1000000",
					"profile_name": "nsroot_Notnsroot250",
					"sync_operation": "false",
					"throughput": "1000",
					"throughput_allocation_mode": "0",
					"vlan_id_0_1": "",
					"vlan_id_0_2": "",
					"vlan_type": 1,
					"vm_memory_total": "2048"
				},
				"target": "10.222.74.135"
			}
  }
  EOF
}
