# Inventory manager for OpenToFu
Simple Go application to compare two Kong configurations via Go bindings for Kong's Admin API.

**Written for learing Go (it means, the code is very ugly)**

Application manages inventory and executes opentofu with generated configs

## Build

`go build -o tofugu`

## Use

`./tofugu -a test -- init`

## Compatibility

`tofugu` is compatible with any OpenTofu or Terraform version

## License

`tofugu` is licensed with Apache License Version 2.0.
Please read the LICENSE file for more details.
