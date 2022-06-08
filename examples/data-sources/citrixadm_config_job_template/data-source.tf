data "citrixadm_config_job_template" "template1" {
  name = "ShowConfiguration"
}

output "template_id" {
  value = data.citrixadm_config_job_template.template1.id
}