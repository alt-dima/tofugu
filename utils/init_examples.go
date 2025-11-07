package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const exampleOrg = "example-org"

// writeExampleFile is a helper function to write a file with given content and log the result
func writeExampleFile(targetPath, content, fileType string) {
	// Create parent directories if they don't exist
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		log.Fatalf("Failed to create directory for %s: %v", fileType, err)
	}

	// Write the file
	if err := os.WriteFile(targetPath, []byte(content), 0644); err != nil {
		log.Fatalf("Failed to create %s file: %v", fileType, err)
	}

	fmt.Printf("Created %s file: %s\n", fileType, targetPath)
}

// CreateExampleStructure creates an example organization structure with sample configuration files
func CreateExampleStructure(targetDir string, useToasterDB bool) string {
	var exampleCmd string

	// Create simple manifest file for demo tofi
	manifestContent := `{
  "dimensions": ["account", "datacenter"]
}`
	manifestPath := filepath.Join(targetDir, "examples", "tofies", exampleOrg, "demo", "tofi_manifest.json")
	writeExampleFile(manifestPath, manifestContent, "example tofi manifest")

	// Create output.tf file
	outputContent := `output "account_region_from_inv" {
    value = var.tofugu_account_data.region
}

output "account_id_from_inv" {
    value = var.tofugu_account_data.account_id
}

output "datacenter_name_from_inv" {
    value = var.tofugu_datacenter_name
}
	
output "datacenter_vpc_cidr_from_inv" {
    value = var.tofugu_datacenter_data.vpc_cidr
}
	
output "datacenter_az_count_from_inv" {
    value = var.tofugu_datacenter_data.az_count
}`
	outputPath := filepath.Join(targetDir, "examples", "tofies", exampleOrg, "demo", "output.tf")
	writeExampleFile(outputPath, outputContent, "example output.tf")

	if !useToasterDB {

		// Create a simple example account.json file
		exampleAccountContent := `{
  "account_id": "123456789012",
  "region": "us-east-1"
}`
		accountPath := filepath.Join(targetDir, "examples", "inventory", exampleOrg, "account", "dev.json")
		writeExampleFile(accountPath, exampleAccountContent, "example account")

		// Create example datacenter.json file
		exampleDevDatacenterContent := `{
  "vpc_cidr": "10.0.0.0/16",
  "az_count": 1
}`
		datacenterDevPath := filepath.Join(targetDir, "examples", "inventory", exampleOrg, "datacenter", "dev.json")
		writeExampleFile(datacenterDevPath, exampleDevDatacenterContent, "example dev datacenter")

		// Create example datacenter.json file
		exampleProdDatacenterContent := `{
  "vpc_cidr": "10.2.0.0/16",
  "az_count": 3
}`
		datacenterProdPath := filepath.Join(targetDir, "examples", "inventory", exampleOrg, "datacenter", "prod.json")
		writeExampleFile(datacenterProdPath, exampleProdDatacenterContent, "example prod datacenter")
	} else {
		exampleCmd = "   export toasterurl=https://662cab7c5e116819738b01fe:supertoaster@toaster.altuhov.su\n"
	}

	exampleCmd = exampleCmd +
		"   tofugu cook -o " + exampleOrg + " -d account:dev -d datacenter:dev -t demo -- init\n" +
		"   tofugu cook -o " + exampleOrg + " -d account:dev -d datacenter:dev -t demo -- plan\n" +
		"Notice output from opentofu and now execute with prod datacenter:\n" +
		"   tofugu cook -o " + exampleOrg + " -d account:dev -d datacenter:prod -t demo -- init\n" +
		"   tofugu cook -o " + exampleOrg + " -d account:dev -d datacenter:prod -t demo -- plan\n" +
		"Notice different output from opentofu based on different inventory data but same Terraform code!\n"

	return exampleCmd
}

// GetTofuguConfigContent returns the content for the .tofugu configuration file
func CreateExampleTofuguConfigFile(tofuguConfigPath string) {
	content := `defaults:
  tofies_path: examples/tofies
#  shared_modules_path: examples/tofies/shared-modules
  inventory_path: examples/inventory
  cmd_to_exec: tofu
#  backend:
#    bucket: my-tfstates
#    key: $tofugu_state_path
#    region: us-east-1
# Add additional organization-specific configurations as needed:
# example-org:
#   backend:
#     bucket: example-org-tfstates
#     prefix: $tofugu_state_path
`
	writeExampleFile(tofuguConfigPath, content, ".tofugu configuration")
}
