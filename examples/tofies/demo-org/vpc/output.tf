output "region_from_env" {
    value = var.tofugu_envvar_awsregion
}

output "region_from_inv" {
    value = var.tofugu_account_manifest.region
}