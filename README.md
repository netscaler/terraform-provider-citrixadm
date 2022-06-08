# Terraform `Citrix ADM Service` Provider

Terraform provider for [Citrix ADM Service](https://docs.citrix.com/en-us/citrix-application-delivery-management-service/citrix-application-delivery-management-service.html) provides [Infrastructure as Code (IaC)](https://en.wikipedia.org/wiki/Infrastructure_as_code) to manage your ADCs via ADM. Using the terraform provider you can onboard ADCs in ADM, assign licenses, create and trigger stylebooks, run configpacks etc.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 1.x.x
- [Go](https://golang.org/doc/install) 1.11+ (to build the provider plugin)

## Building The Provider

> In [Makefile](./Makefile) Change the `OS_ARCH` variable to the architecture of your system.
> For Eg: `OS_ARCH=linux_amd64` OR `OS_ARCH=darwin_amd64` OR `OS_ARCH=windows_amd64`

```sh
git clone git@github.com:citrix/terraform-provider-citrixadc
cd terraform-provider-citrixadc
make
```

## Using the provider

Documentation can be found [here](./PROVIDER_USAGE.md).
