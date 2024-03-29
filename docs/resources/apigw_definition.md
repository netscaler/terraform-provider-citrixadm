---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "citrixadm_apigw_definition Resource - terraform-provider-citrixadm"
subcategory: ""
description: |-
  Create and Manage API Definition
---

# citrixadm_apigw_definition (Resource)

Create and Manage API Definition

## Example Usage

```terraform
resource "citrixadm_apigw_definition" "tf_def1" {
  name     = "tf-def"
  version  = "V2"
  title    = "my_tf_api"
  host     = "example.com"
  basepath = "/"
  schemes  = []
  apiresources {
    paths   = "/user"
    methods = ["GET", "PUT"]
  }
  apiresources {
    paths   = "/user/action"
    methods = ["POST"]
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `apiresources` (Block List, Min: 1) List of HTTP Methods and API Resource paths. (see [below for nested schema](#nestedblock--apiresources))
- `host` (String) Host FQDN where API service is hosted
- `name` (String) API Definition name
- `schemes` (List of String) Schemes of API Definition , HTTP/HTTPS
- `title` (String) Title for API Definition
- `version` (String) API Definition version

### Optional

- `basepath` (String) API Definition base path - this is appended as prefix to all API resources

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--apiresources"></a>
### Nested Schema for `apiresources`

Optional:

- `methods` (List of String) API Method
- `paths` (String) API Path


