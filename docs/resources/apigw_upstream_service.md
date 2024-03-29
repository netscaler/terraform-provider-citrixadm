---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "citrixadm_apigw_upstream_service Resource - terraform-provider-citrixadm"
subcategory: ""
description: |-
  Configure Upstream service for the provided API Deployment Id
---

# citrixadm_apigw_upstream_service (Resource)

Configure Upstream service for the provided API Deployment Id

## Example Usage

```terraform
data "citrixadm_apigw_deployment" "demo" {
  name = "deploy"
}
resource "citrixadm_apigw_upstream_service" "service1" {
  name              = "tf_upstream_service"
  scheme            = "http"
  service_fqdn      = "www.example.com"
  service_fqdn_port = 80
  deployment_id     = data.citrixadm_apigw_deployment.demo.id
  backend_servers {
    ip4_addr = "2.1.2.4"
    port     = 80
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `backend_servers` (Block List, Min: 1) List of all IPv4 or IPv6 host address and Port for Upstream Service. (see [below for nested schema](#nestedblock--backend_servers))
- `deployment_id` (String) API Deployment Id
- `name` (String) Upstream Service Name

### Optional

- `scheme` (String) Protocol used to exchange data with the service
- `service_fqdn` (String) Hostname of Upstream Service as FQDN where API traffic will be sent
- `service_fqdn_port` (Number) Optional Listening port for Service FQDN

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--backend_servers"></a>
### Nested Schema for `backend_servers`

Required:

- `ip4_addr` (String) IPv4 host address for the Upstream Server/Service where API traffic will be sent
- `port` (Number) Port number for Upstream Server/Service


