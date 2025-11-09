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
	ToasterUrl        string
	Workspace         string
}

type tofiManifestStruct struct {
	Dimensions []string
}

type ToasterResponse struct {
	Error      string
	Dimensions []DimensionInToaster
}

type DimensionInToaster struct {
	ID        string
	WorkSpace string
	DimData   map[string]interface{}
}

type HistoryStruct struct {
	CmdToExec  string
	CmdArgs    []string
	CmdMainArg string
	ExitCode   int
	Dimensions map[string]string
}
