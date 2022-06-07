# Terraform `Citrix ADM Service` Provider

- Website: https://www.terraform.io

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 1.0.x
- [Go](https://golang.org/doc/install) 1.11 (to build the provider plugin)

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
