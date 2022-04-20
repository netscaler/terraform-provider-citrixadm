terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = "4.10.0"
    }
  }
}

provider "aws" {
  # Configuration options
}


data "aws_ami" "example" {
#   executable_users = ["self"]
  most_recent      = true
#   name_regex       = "^Citrix ADC"
  owners           = ["aws-marketplace"]

  filter {
    name   = "name"
    values = ["Citrix ADC 13.1*63425ded-82f0-4b54-8cdd-6ec8b94bd4f8*"]
  }
}

output "latestami" {
  value = data.aws_ami.example.id
}

