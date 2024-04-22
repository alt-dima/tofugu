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

Config file path maybe provided by the `--config` flag, for example: `
```
./tofugu --config path_to_config/tofuguconfig cook -o demo-org -d account:test-account -d datacenter:staging1 -t vpc -- init
```
If `--config` flag is not set, then it will try to load from default location `$HOME/.tofugu`

Example config file:
```
defaults:
  tofies_path: examples/tofies
  shared_modules_path: examples/tofies/shared-modules
  inventory_path: examples/inventory
  cmd_to_exec: tofu
  s3_bucket_name: default-tfstates
  s3_bucket_region: us-east-2
demo-org:
  s3_bucket_name: demo-org-tfstates
```

- `tofies_path` = relative path to the folder with terraform code (`tofi`)
- `shared_modules_path` = relative path to the folder with shared TF modules maybe used by any `tofi`
- `inventory_path` =  relative path to the folder with jsons
- `cmd_to_exec` = name of the binary to execute (`tofu` or `terraform`)
- `s3_bucket_name` = name of the S3 bucket to store state
- `s3_bucket_region` = region of the S3 bucket to store state

If for the `demo-org` config `s3_bucket_name` is set, then S3 key (path) will be generated like: `s3://demo-org-tfstates/dimName1_dimValue1/dimNameN_dimValueN/tofiName.tfstate`

If for the `demo-org` config `s3_bucket_name` is NOT set, then S3 key (path) will be generated like `s3://default-tfstates/org_demo-org/dimName1_dimValue1/dimNameN_dimValueN/tofiName.tfstate`

This could be useful, if you want to store by default tfstate for all the organisations in the same/default bucket `default-tfstates` but for some specific organisation you need to store tfstates in dedicated bucket `demo-org-tfstates`

At least 
```
defaults:
  s3_bucket_name: default-tfstates
  s3_bucket_region: us-east-2
```
must be set in the config file!

Other options contain hard-coded defaults:
```
	viper.SetDefault("defaults.inventory_path", "examples/inventory")
	viper.SetDefault("defaults.shared_modules_path", "examples/tofies/shared-modules")
	viper.SetDefault("defaults.tofies_path", "examples/tofies")
	viper.SetDefault("defaults.cmd_to_exec", "tofu")
```

## License

`tofugu` is licensed with Apache License Version 2.0.
Please read the LICENSE file for more details.
