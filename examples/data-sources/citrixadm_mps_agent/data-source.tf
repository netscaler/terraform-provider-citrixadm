data "citrixadm_mps_agent" "agent1" {
  name = "10.0.1.91"
}

output "agent_id" {
  value = data.citrixadm_mps_agent.agent1.id
}

output "datacenter_id" {
  value = data.citrixadm_mps_agent.agent1.datacenter_id
}