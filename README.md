# Inventory manager for OpenTofu or Terraform
Manage your infrastructure and environments with Inventory manager and OpenTofu/Terraform!
Avoid duplication of the TF code! Reuse same code for multive enviroments with configuration in dedicated json files.

**Written for learing Go (it means, the code is very ugly)** and cobra and viper

Application manages inventory and executes opentofu with generated configs

## Build

`go build -o tofugu`

## Use

```
./tofugu cook -o demo-org -d account:test-account -d datacenter:staging1 -t vpc -- init
./tofugu cook -o demo-org -d account:test-account -d datacenter:staging1 -t vpc -- plan
./tofugu cook -o demo-org -d account:test-account -d datacenter:staging1 -t vpc -- apply
```

Everything after `--` will be passed as parameters to the `cmd_to_exec`

## Compatibility

`tofugu` is OpenTofu/Terraform version agnostic!

## $HOME/.tofurc

Recommended to enable plugin_cache_dir to reuse providers

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
  s3_bucket_name: asu-tfstates
  s3_bucket_region: us-east-2
```

- `tofies_path` = relative path to the folder with terraform code (`tofi`)
- `shared_modules_path` = relative path to the folder with shared TF modules maybe used by any `tofi`
- `inventory_path` =  relative path to the folder with jsons
- `cmd_to_exec` = name of the binary to execute (`tofu` or `terraform`)
- `s3_bucket_name` = name of the S3 bucket to store state
- `s3_bucket_region` = region of the S3 bucket to store state

S3 key (path) will be generated like `dimName1/dimNameN/tofiName.tfstate`

## License

`tofugu` is licensed with Apache License Version 2.0.
Please read the LICENSE file for more details.
