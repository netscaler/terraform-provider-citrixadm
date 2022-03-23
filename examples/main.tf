resource "citrixadm_ns_device_profile" "profile1" {
    // concatenate name and random_string
    name = "tf_test_profile"
    # name = format("tf_test_profile_%s", random_string.random.result)
    username = "nsroot"
    password = "tfnsroot"
    snmpsecurityname = "tf_test_snmp"
}

# resource "random_string" "random" {
#   length           = 1
#   special = false
# }