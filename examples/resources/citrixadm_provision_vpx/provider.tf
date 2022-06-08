terraform {
  required_providers {
    citrixadm = {
      source = "citrix/citrixadm"
    }
  }
}
provider "citrixadm" {
  # host          = "https://railay.adm.cloud.com"                # Optionally use CITRIXADM_HOST env var
  # client_id     = "" # Optionally use CITRIXADM_CLIENT_ID env var
  # client_secret = "ABcdefgHIJklmnopqrstuv=="             # Optionally use CITRIXADM_CLIENT_SECRET env var
  # host_location = "us"                                   # us, eu # Optionally use CITRIXADM_HOST_LOCATION env var
  # customer_id   = "abcdefghijkl"                         # Optionally use CITRIXADM_CUSTOMER_ID env var
  # fail_on_stall = true                                   # Optionally use CITRIXADM_FAIL_ON_STALL env var
}
