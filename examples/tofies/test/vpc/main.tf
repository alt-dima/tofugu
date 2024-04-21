//use shared-module
module "vpc" {
  source = "./shared-modules/create_vpc"
  cidr = var.tofugu_datacenter_manifest[var.tofugu_account_name].cidr
}