data "citrixadm_managed_device" "device1" {
  ip_address = "10.0.1.42"
}

output "device1_id" {
  value = data.citrixadm_managed_device.device1.id
}