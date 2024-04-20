# Inventory manager for OpenToFu
Simple Go application to compare two Kong configurations via Go bindings for Kong's Admin API.

**Written for learing Go (it means, the code is very ugly)** and cobra and viper

Application manages inventory and executes opentofu with generated configs

## Build

`go build -o tofugu`

## Use

`./tofugu -a test -- init`

## Compatibility

`tofugu` is compatible with any OpenTofu or Terraform version

## $HOME/.tofurc

```
plugin_cache_dir   = "$HOME/.terraform.d/plugin-cache"
plugin_cache_may_break_dependency_lock_file = true
```

`mkdir "$HOME/.terraform.d/plugin-cache"`

## $HOME/.tofugu

```
defaults:
  tofies_path: examples/tofies
  shared_modules_path: examples/tofies/shared-modules
  inventory_path: examples/inventory
  cmd_to_exec: tofu
```

## License

`tofugu` is licensed with Apache License Version 2.0.
Please read the LICENSE file for more details.
