# Infrastructure layers configuration orchestrator for OpenTofu or Terraform
`tofugu` is an infrastructure layers configuration orchestrator that dynamically manages OpenTofu or Terraform. It provides infrastructure configuration definitions from outside the Terraform code, using either files or an Infrastructure Layers Configuration Management Database (CMDB) called Toaster-ToasterDB. This allows you to reuse Terraform code across multiple environments instead of duplicating it.

- Environment/layer configuration stored outside the Terraform/OpenTofu code
- Terraform/OpenTofu code, called **tofi**, should be generic enough to handle provided configuration to deploy the same resources with different configurations
- `tfvars` and `variables` are automatically generated in the temporary folder with selected terraform code (**tofi**) resulting in a full set of the Terraform code and configuration variables
- After the temporary folder is ready, it executes `terraform` or `tofu` with specified parameters
- Maintains separate state files for each environment/layer, automatically providing configuration for remote state management (different path on the storage regarding configured layers/dimensions). So the deployed set (configuration + terraform) is stored in different `tfstate` files in remote storage (S3, GCS)

## Quick start with AI Coding Assistants

Getting started with `tofugu` is even easier using AI coding assistants:

1. Open the repository in your preferred editor with an AI assistant installed:
   - GitHub Copilot in VS Code
   - Cursor
   - Claude/Anthropic
   - WindSurf
   - Cline

2. Ask questions like:
   - "How do I set up tofugu for AWS resources?"
   - "Help me create a new tofi for a GCP instance"
   - "How do I use Toaster-ToasterDB with tofugu?"
   - "Show me how to pass environment variables to my terraform code"

These instructions provide context about:
- The architecture and key concepts of tofugu
- How to work with tofies, dimensions, and inventory sources
- Integration with Toaster-ToasterDB for centralized configuration
- Project-specific conventions and workflows

For more complex tasks, you can ask the AI assistant to guide you step-by-step through creating configurations, setting up the backend, or troubleshooting deployment issues.

## What about alternative tools?

Yes, you should check other Infrastructure as Code (IaC) orchestration tools for Terraform:
- [Atmos](https://github.com/cloudposse/atmos)
- [Digger](https://github.com/diggerhq/digger)
- [Spacelift](https://github.com/spacelift-io)
- [Terragrunt](https://github.com/gruntwork-io/terragrunt)
- [Terramate](https://github.com/terramate-io/terramate)

So why another tool?
1. The more open source tools, the better (for GitHub Copilot, not for Earth).
2. The more choice, the better (there are countries where people have no choice).
3. Imagine there is only AWS CloudFormation.

## Usage

Organization with AWS resources and state stored in S3
```bash
./tofugu cook --config examples/.tofugu -o demo-org -d account:test-account -d datacenter:staging1 -t vpc -- init
./tofugu cook --config examples/.tofugu -o demo-org -d account:test-account -d datacenter:staging1 -t vpc -- plan
./tofugu cook --config examples/.tofugu -o demo-org -d account:test-account -d datacenter:staging1 -t vpc -- apply
```

Organization with Google Cloud resources and state stored in Google Cloud Storage
```bash
./tofugu cook --config examples/.tofugu -o gcp-org -d account:free-tier -t free_instance -- init
./tofugu cook --config examples/.tofugu -o gcp-org -d account:free-tier -t free_instance -- plan
./tofugu cook --config examples/.tofugu -o gcp-org -d account:free-tier -t free_instance -- apply
```

- Everything after `--` will be passed as parameters to the `cmd_to_exec`
- `-c` = to remove temp dir after any `tofugu` execution (after `apply` or `destroy` and exit code=0 temp dir removed automatically)
- `-o` = name of the `organization` (subfolder in **Inventory**, **tofies** folders and in `.tofugu` config section)
- `-d` = `dimension` to attach to tofu/terraform. You may specify as many `-d` pairs as you need!
- `-t` = name of the `tofi` in the `tofies` folder

## Tofi Manifest

Special JSON file with the name `tofi_manifest.json` in the `tofi` folder provides options for TofuGu.

Currently, only `dimensions` with a list of the required/expected dimensions (from **Inventory Store**)

[tofi_manifest.json example](examples/tofies/demo-org/vpc/tofi_manifest.json)

## Infrastructure layers/dimensions configurations Storage

### Infrastructure layers Configuration Management Database (CMDB). (Toaster-ToasterDB)
You could set the env variable `toasterurl` to point to TofuGu-Toaster, like:
```bash
export toasterurl='https://accountid:accountpass@toaster.altuhov.su'
```

To generate your own credentials please go to [https://toaster.altuhov.su/](https://toaster.altuhov.su/), fill the form with Account Name, Email, and press `Create User` and you will receive generated credentials and a ready-to-use export command like:
```
Please execute in shell to set toasterurl:

export toasterurl=https://6634b72292e9e996105de19e:generatedpassword@toaster.altuhov.su
```

With the correct `toasterurl`, TofuGu will connect and receive all the required dimension data from the Toaster-ToasterDB.
An additional parameter could be passed to tofugu `-w workspacename`. In general, `workspacename` is the branch name of the source repo where the dimension is stored. If TofuGu-Toaster does not find the dimension with the specified `workspacename`, it will try to return the dimension from the `master` workspace/branch!

**Toaster-ToasterDB** provides additional features for your CI and CD pipelines. For example, you need to receive a [first-app.json](examples/inventory/demo-org/application/first-app.json) in the CI pipeline, to check the application configuration.
Or you need a list of all the datacenters in the [datacenter dimension](examples/inventory/demo-org/datacenter) in a [Jenkins drop-down](https://github.com/alt-dima/tofugu/issues/10#issuecomment-2090932416) list to select to which datacenter to deploy the application.

[Swagger API docs (full API documentation and examples)](https://app.swaggerhub.com/apis-docs/altuhovsu/tofugu_toaster_api/)

To upload/update dimensions in Toaster from your Inventory Files repo you could use [inventory-to-toaster.sh script example](examples/inventory-to-toaster.sh) and execute it like `bash examples/inventory-to-toaster.sh examples/inventory/`

Please join the [Toaster-ToasterDB beta-testers!](https://github.com/alt-dima/tofugu/issues/10)

### File-based Infrastructure layers configuration Storage. (Inventory Files)

If the env variable `toasterurl` is not set, TofuGu will use file-based configuration Storage (probably dedicated git repo), specified by the path configured in `inventory_path`.

Examples:

- [staging1.json in Inventory Files](examples/inventory/demo-org/datacenter/staging1.json)
- [dim_defaults.json in Inventory Files](examples/inventory/demo-org/datacenter/dim_defaults.json)

### Dimensions usage in tf-code

When you set dimensions in the tofugu flags `-d datacenter:staging1`, TofuGu will provide you inside tf-code next variables:

- var.tofugu_datacenter_name = will contain string `staging1`
- var.tofugu_datacenter_data = will contain the whole object from `staging1.json`
- var.tofugu_datacenter_defaults = will contain the whole object from `dim_defaults.json` IF the file `dim_defaults.json` exists!

Examples:

- [datacenter object with defaults used in tf-code](examples/tofies/demo-org/vpc/main.tf#L5)

## Passing environment variables from shell

For example, you need to pass a variable (AWS region) from shell to the terraform code, simply set it and use it!

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

Config file (in YAML format) path may be provided by the `--config` flag, for example:
```bash
tofugu --config path_to_config/tofuguconfig cook -o demo-org -d account:test-account -d datacenter:staging1 -t vpc -- init
```
If the `--config` flag is not set, then it will try to load from the default location `$HOME/.tofugu`

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
- `inventory_path` = relative path to the folder with JSONs
- `cmd_to_exec` = name of the binary to execute (`tofu` or `terraform`)
- `backend` = Config values for backend provider. All the child key:values will be provided to `init` and `$tofugu_state_path` will be replaced by the generated path.
For example, when you execute `tofugu cook ...... -- init`, TofuGu actually will execute `init -backend-config=bucket=gcp-tfstates -backend-config=prefix=account_free-tier/free_instance.tfstate`

At least 
```yaml
defaults:
  backend:
    bucket: default-tfstates
    key: $tofugu_state_path
```
must be set in the config file! With key:values specific for the backend provider being used in org!

Other options contain hard-coded defaults:
```yaml
defaults:
  inventory_path: "examples/inventory"
  shared_modules_path: ""
  tofies_path: "examples/tofies"
  cmd_to_exec: "tofu"
```

## Shared modules support

It is a good practice to move some generic terraform code to the `modules` and reuse those modules in multiple terraform code (**tofies**)

Path to the folder with such private shared modules is configured by the `shared_modules_path` parameter in the `.tofugu` configuration file.

This folder will be mounted/linked to every temporary folder (set) so you could use any module by short path like
```
//use shared-module
module "vpc" {
  source = "./shared-modules/create_vpc"
}
```
Examples:
- [Shared module for VPC creation](examples/tofies/shared-modules/create_vpc)
- [Shared module for VPC creation used in code](examples/tofies/demo-org/vpc/main.tf#L3)

## Remote state (Terraform Backend where state data files are stored)

AWS, Google Cloud, and some other backends are supported! You could configure any backend provider in the [TofuGu Config file](#hometofugu)

[For AWS S3 your terraform code (`tofi`) should contain at least:](examples/tofies/demo-org/vpc/versions.tf#L4):
```
terraform {
  backend "s3" {}
}
```

[For Google Cloud Storage your terraform code (`tofi`) should contain at least:](examples/tofies/gcp-org/free_instance/versions.tf#L4):
```
terraform {
  backend "gcs" {}
}
```

If for the `demo-org` config `bucket` is set, then `$tofugu_state_path` will be like: `dimName1_dimValue1/dimNameN_dimValueN/tofiName.tfstate`

If for the `demo-org` config `bucket` is NOT set, then `$tofugu_state_path` will be like `org_demo-org/dimName1_dimValue1/dimNameN_dimValueN/tofiName.tfstate`

This could be useful if you want to store by default tfstate for all the organizations in the same/default bucket `default-tfstates` but for some specific organization you need to store tfstates in a dedicated bucket `demo-org-tfstates`

## Data Source Configuration (data "terraform_remote_state")

To simplify "Data Source Configuration" (`data "terraform_remote_state" "tfstate" { }`) it will be nice to have backend config values as tfvars.

`var.tofugu_backend_config` will contain all the parameters from [TofuGu config (backend Section)](#hometofugu)

[For example, for AWS S3](examples/tofies/demo-org/vpc/data.tf):
```
data "terraform_remote_state" "network" {
  backend = "s3"
  config = {
    bucket = var.tofugu_backend_config.bucket
    key    = "network/terraform.tfstate"
    region = var.tofugu_backend_config.region
  }
}
```
[And for GCS](examples/tofies/gcp-org/free_instance/data.tf):
```
data "terraform_remote_state" "free_instance" {
  backend = "gcs"
  config = {
    bucket  = var.tofugu_backend_config.bucket
    prefix  = "account_free-tier/free_instance.tfstate"
  }
}
```

You will set `key/prefix` to another tofie's tfstate, which outputs you want to use.

## $HOME/.tofurc

Recommended to enable plugin_cache_dir to reuse providers.

[.tofurc example](examples/.tofurc):

```ini
plugin_cache_dir   = "$HOME/.terraform.d/plugin-cache"
plugin_cache_may_break_dependency_lock_file = true
```
Do not forget to create the plugin-cache dir: `mkdir "$HOME/.terraform.d/plugin-cache"`

## Compatibility

`tofugu` is OpenTofu/Terraform version agnostic!
Required external tools/binaries: `rsync`, `ln`

## License

`tofugu` is licensed with Apache License Version 2.0.
Please read the LICENSE file for more details.
