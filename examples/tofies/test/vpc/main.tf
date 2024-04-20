module "vpc" {
  source = "./shared-modules/create_vpc"
  cidr = "10.1.0.0/16"
}