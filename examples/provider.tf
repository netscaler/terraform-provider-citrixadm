terraform {
	required_providers {
	  citrixadm = {
		source = "citrix/citrixadm"
	  }

	}
  }
  provider "citrixadm" {
	host          = "https://adm.cloud.com"
	client_id     = "987a9390-6a65-4f78-8587-790feb82d63a"
	client_secret = "BZnvbtdZaWJ3jYJwYtnsCw=="
	host_location = "us" // eu
	customer_id = "vbd3nm32fn5w"
  }
