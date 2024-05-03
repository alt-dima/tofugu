# Cloud Native Inventory Manager for OpenTofu or Terraform
Manage your infrastructure and environments with Inventory manager and OpenTofu/Terraform!
Avoid duplication of the TF code! Reuse same code for multive enviroments with configuration in dedicated json files.

**Written for learing Go** (and cobra and viper)

Application manages inventory and executes opentofu with terraform variables.

No need to manually create any `tfvars` or `variables` files/directives -> [empty variables.tf](examples/tofies/demo-org/vpc/variables.tf)

## Usage

Org with AWS resources and state stored in S3
```bash
./tofugu cook --config examples/.tofugu -o demo-org -d account:test-account -d datacenter:staging1 -t vpc -- init
./tofugu cook --config examples/.tofugu -o demo-org -d account:test-account -d datacenter:staging1 -t vpc -- plan
./tofugu cook --config examples/.tofugu -o demo-org -d account:test-account -d datacenter:staging1 -t vpc -- apply
```

Org with Google Cloud resources and state stored in Google Cloud Storage
```bash
./tofugu cook --config examples/.tofugu -o gcp-org -d account:free-tier -t free_instance -- init
./tofugu cook --config examples/.tofugu -o gcp-org -d account:free-tier -t free_instance -- plan
./tofugu cook --config examples/.tofugu -o gcp-org -d account:free-tier -t free_instance -- apply
```

- Everything after `--` will be passed as parameters to the `cmd_to_exec`
- `-c` = to remove temp dir after any `tofugu` execution (after `apply` or `destroy` and exitcode=0 temp dir removed automatically)
- `-o` = name of the `organization` (subfolder in **Inventory**, **tofies** folders and in `.tofugu` config section)
- `-d` = `dimension` to attach to tofu/terraform. You may specify as many `-d` pairs as you need!
- `-t` = name of the `tofi` in the `tofies` folder

## Tofi Manifest

Special json file with name `tofi_manifest.json` in `tofi` folder provides options for TofuGu.

Currently only `dimensions` with list of the required/expecting dimensions (from **Inventory Store**)

[tofi_manifest.json example](examples/tofies/demo-org/vpc/tofi_manifest.json)

## Inventory (dimensions) Store

### Cloud Native Inventory Storage (Toaster-ToasterDB)
You could set env variable `toasterurl` to point to TofuGu-Toaster, like:
```bash
export toasterurl='https://accountid:accountpass@toaster.altuhov.su'
```

To generate your own credentials please go to [https://toaster.altuhov.su/](https://toaster.altuhov.su/) , fill the form with Account Name, Email and press `Create User` and you will receive generated credentials and ready-to-use export command like:
```
Please execute in shell to set toasterurl:

export toasterurl=https://6634b72292e9e996105de19e:generatedpassword@toaster.altuhov.su
```

With correct `toasterurl` TofuGu will connect and receive all the required dimension data from the Toaster-ToasterDB.
Additional parameter could be passed to tofugu `-w workspacename`. In general `workspacename` is the branch name of the source repo where the dimension is stored. If TofuGu-Toaster will not find dimension with specified `workspacename` it will try to return dimension from `master` workspace/branch!

**Toaster-ToasterDB** Provides additional features for your CI and CD pipelines. For example, you need to receive a [first-app.json](examples/inventory/demo-org/application/first-app.json) in the CI pipeline, to check application configuration.
Or you need a list of all the datacenters in [datacenter dimension](examples/inventory/demo-org/datacenter) in [Jenkins drop-down](https://github.com/alt-dima/tofugu/issues/10#issuecomment-2090932416) list to select to which datacenter to deploy application.

[Swagger API docs (full API documentation and examples)](https://app.swaggerhub.com/apis-docs/altuhovsu/tofugu_toaster_api/)

To upload/update dimensions in Toaster from your Inventory Files repo you could use [inventory-to-toaster.sh script example](examples/inventory-to-toaster.sh) and execute it like `bash examples/inventory-to-toaster.sh  examples/inventory/`

Please join the [Toaster-ToasterDB beta-testers!](https://github.com/alt-dima/tofugu/issues/10)


### Inventory Files repo

If env variable `toasterurl` is not set, TofuGu will use file-based inventory storage, by the path configured in `inventory_path`

Examples:

- [staging1.json in Inventory Files](examples/inventory/demo-org/datacenter/staging1.json)
- [dim_defaults.json in Inventory Files](examples/inventory/demo-org/datacenter/dim_defaults.json)

### Dimensions usage in tf-code

When you set dimensions in the tofugu flags `-d datacenter:staging1 `, TofuGu will provide you inside tf-code next variables:

- var.tofugu_datacenter_name = will contain string `staging1`
- var.tofugu_datacenter_manifest = will contain whole object from `staging1.json`
- var.tofugu_datacenter_defaults = will contain whole object from `dim_defaults.json` IF file `dim_defaults.json` exists!

Examples:

- [datacenter object with defaults used in tf-code](examples/tofies/demo-org/vpc/main.tf#L5)

## Passing environment variables from shell

For example, you need to pass a variable (AWS region) from shell to the terraform code, simply set it and use!

**Environment variable must start with `tofugu_envvar_` prefix!**
```bash
export tofugu_envvar_awsregion=us-east-1
```
In the TF code:
```
provider "aws" {
    region = var.tofugu_envvar_awsregion
}
```

[Env variables used in code example](examples/tofies/demo-org/vpc/providers.tf#L3)

## $HOME/.tofugu

Config file (in YAML format) path maybe provided by the `--config` flag, for example: `
```bash
./tofugu --config path_to_config/tofuguconfig cook -o demo-org -d account:test-account -d datacenter:staging1 -t vpc -- init
```
If `--config` flag is not set, then it will try to load from default location `$HOME/.tofugu`

[.tofugu example](examples/.tofugu):
```yaml
defaults:
  tofies_path: examples/tofies
  shared_modules_path: examples/tofies/shared-modules
  inventory_path: examples/inventory
  cmd_to_exec: tofu
  backend:
    bucket: default-tfstates
    key: $tofugu_state_path
    region: us-east-2
gcp-org:
  backend:
    bucket: gcp-tfstates
    prefix: $tofugu_state_path
```

- `tofies_path` = relative path to the folder with terraform code (`tofi`)
- `shared_modules_path` = relative path to the folder with shared TF modules maybe used by any `tofi`
- `inventory_path` =  relative path to the folder with jsons
- `cmd_to_exec` = name of the binary to execute (`tofu` or `terraform`)
- `backend` = Config values for backend provider. All the child key:values will be provided to `init` and `$tofugu_state_path` will be replaced by generated path.
For example, it will look like `tofu init -backend-config=bucket=gcp-tfstates -backend-config=prefix=account_free-tier/free_instance.tfstate`

At least 
```yaml
defaults:
  backend:
    bucket: default-tfstates
    key: $tofugu_state_path
```
must be set in the config file! With key:values specific for the backend provider being used in org!

Other options contain hard-coded defaults:
```go
	viper.SetDefault("defaults.inventory_path", "examples/inventory")
	viper.SetDefault("defaults.shared_modules_path", "examples/tofies/shared-modules")
	viper.SetDefault("defaults.tofies_path", "examples/tofies")
	viper.SetDefault("defaults.cmd_to_exec", "tofu")
```

## Remote state in S3

AWS, Google Cloud and some other backends are supported! You could configure any backend provider in `tofugu config file`

[For AWS S3 your terraform code (`tofi`) should contains at least:](examples/tofies/demo-org/vpc/versions.tf#L4):
```
terraform {
  backend "s3" {}
}
```

[For Google Cloud Storage your terraform code (`tofi`) should contains at least:](examples/tofies/gcp-org/free_instance/versions.tf#L4):
```
terraform {
  backend "gcs" {}
}
```

If for the `demo-org` config `bucket` is set, then `$tofugu_state_path` will be like: `dimName1_dimValue1/dimNameN_dimValueN/tofiName.tfstate`

If for the `demo-org` config `bucket` is NOT set, then `$tofugu_state_path` will be like `org_demo-org/dimName1_dimValue1/dimNameN_dimValueN/tofiName.tfstate`

This could be useful, if you want to store by default tfstate for all the organisations in the same/default bucket `default-tfstates` but for some specific organisation you need to store tfstates in dedicated bucket `demo-org-tfstates`


## $HOME/.tofurc

Recommended to enable plugin_cache_dir to reuse providers.

[.tofurc example](examples/.tofurc):

```
plugin_cache_dir   = "$HOME/.terraform.d/plugin-cache"
plugin_cache_may_break_dependency_lock_file = true
```
Do not forget to create plugin-cache dir: `mkdir "$HOME/.terraform.d/plugin-cache"`

## Compatibility

`tofugu` is OpenTofu/Terraform version agnostic!
Required external tools/binaries: `rsync`, `ln`

## Why not Terragrunt?

Not sure, but for me looks like same general idea, but for different cases.
For example: https://terragrunt.gruntwork.io/docs/features/keep-your-terraform-code-dry/#keep-your-terraform-code-dry

> In a separate repo, called, for example, live, you define the code for all of your environments, which now consists of just one terragrunt.hcl file per component (e.g. app/terragrunt.hcl, mysql/terragrunt.hcl, etc).

And you need to configure/copy terragrunt.hcl (and maybe other files) to each folder/environment (prod,qa, stage) with subfolders like app,mysql,vpc
But if I have 20 environments (stage1-stage20) and 50 units (app,mysql,vpc,eks,redis,....) then, if I need to add stage21 I will need to copy all of the files again.

Maybe better when TF code and inventory split by the repos and adding new environment does not require any changes in the TF repo, only add stage21.json in the inventory repo and deploy every unit
like
```bash
./tofugu cook -o demo-org -d account:test-account -d datacenter:staging21 -t vpc -- init
./tofugu cook -o demo-org -d account:test-account -d datacenter:staging21 -t vpc -- apply -auto-approve

./tofugu cook -o demo-org -d account:test-account -d datacenter:staging21 -t eks -- init
./tofugu cook -o demo-org -d account:test-account -d datacenter:staging21 -t eks -- apply -auto-approve
```
P.S. I very respect terragrunt it is prod-grade tool! this "tool" is just go-learning :)

## License

`tofugu` is licensed with Apache License Version 2.0.
Please read the LICENSE file for more details.
