//use shared-module
module "vpc" {
  source = "./shared-modules/create_vpc"
  cidr = var.tofugu_datacenter_data[var.tofugu_account_name].cidr
}

module "vpc_example_simple-vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "5.7.1"

  cidr = var.tofugu_datacenter_data[var.tofugu_account_name].cidr
}