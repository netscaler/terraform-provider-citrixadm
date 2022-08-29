data "citrixadm_apigw_deployment" "demo" {
  name = "tf_deploy"
}
data "citrixadm_apigw_upstream_service" "upstreamdemo" {
  name          = "sdf"
  deployment_id = data.citrixadm_apigw_deployment.demo.id
}

resource "citrixadm_apigw_policy" "policy" {

  policygroup_name = "tf_policy"
  requestpath {
    order_index = 1
    policy_name = "policyName"
    policytype  = "NoAuth"
    config_spec {
      api_resource_paths {
        delete = false
        endpoints = [
          "/user"
        ]
        get   = true
        patch = false
        post  = true
        put   = false
      }
      custom_rules {
        delete = false
        endpoints = [
          "/user/action"
        ]
        get   = false
        patch = false
        post  = false
        put   = false
      }
    }
  }

  upstreamservice_id = data.citrixadm_apigw_upstream_service.upstreamdemo.id
  deployment_id      = data.citrixadm_apigw_deployment.demo.id
  deploy             = true
}

