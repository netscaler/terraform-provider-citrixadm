resource "citrixadm_ns_device_profile" "profile1" {
    name = "tf_test_profile"
    username = "nsroot"
    password = "tfnsroot"
}

# resource "citrixadm_managed_device" "device1" {
#     ip_address = "10.0.1.145"
#     profile_name = "sumanth-adms-terraform-provider-standalone-RegisterADC2ADMServiceStack-QKEVDOTC2VMRK0J"
#     datacenter_id = "bcbf82f7-5451-4e48-8261-caec673c18e1"
#     agent_id = "e485c7d7-b54d-4a7a-a078-3c4150d1117d" # FIXME: Ask George if we need to get the ID or Agent IP and then find the ID internally?
# }