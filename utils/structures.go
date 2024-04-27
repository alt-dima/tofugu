package utils

type Tofugu struct {
	TofiName          string
	OrgName           string
	DimensionsFlags   []string
	TofiPath          string
	SharedModulesPath string
	InventoryPath     string
	TofiManifestPath  string
	ParsedDimensions  map[string]string
	CmdWorkTempDir    string
	TofiManifest      tofiManifestStruct
	StateS3Path       string
}

type tofiManifestStruct struct {
	Dimensions []string
}
