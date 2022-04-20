terraform {
  required_providers {
    citrixadm = {
      source = "citrix/citrixadm"
    }
  }
}
provider "citrixadm" {
  # host          = "https://adm.cloud.com"                # Optionally use CITRIXADM_HOST env var
  # client_id     = "12345678-90ab-cdef-ghij-klmnopqrstuv" # Optionally use CITRIXADM_CLIENT_ID env var
  # client_secret = "ABcdefgHIJklmnopqrstuv=="             # Optionally use CITRIXADM_CLIENT_SECRET env var
  # host_location = "us"                                   # us, eu # Optionally use CITRIX_ADM_HOST_LOCATION env var
  # customer_id   = "abcdefghijkl"                         # Optionally use CITRIX_ADM_CUSTOMER_ID env var
}
