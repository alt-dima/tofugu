//use shared-module
module "vpc" {
  source = "./shared-modules/create_vpc"
  cidr = var.tofugu_datacenter_data[var.tofugu_account_name].cidr
  enable_dns_support = try(var.tofugu_datacenter_data[var.tofugu_account_name].enable_dns_support, var.tofugu_datacenter_defaults[var.tofugu_account_name].enable_dns_support)
}

module "vpc_example_simple-vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "5.7.1"

  cidr = var.tofugu_datacenter_data[var.tofugu_account_name].cidr
}