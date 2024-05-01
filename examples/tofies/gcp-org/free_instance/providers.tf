# provider "aws" {
#     //region = var.tofugu_account_manifest.region
#     region = var.tofugu_envvar_awsregion
# }

provider "google" {
  project = var.tofugu_account_data.project
  region  = var.tofugu_account_data.region
  zone    = var.tofugu_account_data.zone
}