# Usage of the examples

> You can find the setting up of the provider in the [PROVIDER_USAGE.md](./PROVIDER_USAGE.md) file.

Below is the table showing the usecase you would like to achieve and where you can find the respective example.

| Usecase | Terraform Example Folder | Documentation |
|---------|---------|-------------------|
| Add a VPX present in AWS, Azure, GCP  and manage from ADM | [resources/citrixadm_managed_device](examples/resources/citrixadm_managed_device) | [HERE](./docs/resources/managed_device.md)|
| Allocate Pooled license to the managed VPX from ADM | [resources/citrixadm_managed_device_allocate_license](examples/resources/citrixadm_managed_device_allocate_license) | [HERE](./docs/resources/managed_device_allocate_license.md)|
| Create a new NS Device Profile | [resources/citrixadm_ns_device_profile](examples/resources/citrixadm_ns_device_profile) | [HERE](./docs/resources/ns_device_profile.md)|
| Upload a Stylebook to ADM | [resources/citrixadm_stylebook](examples/resources/citrixadm_stylebook) | [HERE](./docs/resources/stylebook.md)|
| Apply an existing Stylebook (by creating a config-pack) in ADM and optionally apply the stylebook config to one or more ADC targets | [resources/citrixadm_stylebook_configpack](examples/resources/citrixadm_stylebook_configpack) | [HERE](./docs/resources/stylebook_configpack.md)|
| Provision VPX on SDX from ADM | [examples/resources/citrixadm_provision_vpx](examples/resources/citrixadm_provision_vpx) | [HERE](./docs/resources/provision_vpx.md)|
