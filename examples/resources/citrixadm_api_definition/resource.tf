resource "citrixadm_api_definition" "tf_def1" {
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
