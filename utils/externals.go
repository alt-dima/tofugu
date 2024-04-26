package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (tofuguStruct *Tofugu) GetStringFromViperByOrgOrDefault(keyName string) string {
	if viper.IsSet(tofuguStruct.OrgName + "." + keyName) {
		return viper.GetString(tofuguStruct.OrgName + "." + keyName)
	} else {
		return viper.GetString("defaults." + keyName)
	}
}

func (tofuguStruct *Tofugu) SetupStateS3Path() {
	var stateS3Path string
	if !viper.IsSet(tofuguStruct.OrgName + ".s3_bucket_name") {
		stateS3Path = stateS3Path + "org_" + tofuguStruct.OrgName + "/"
	}
	for _, dimension := range tofuguStruct.TofiManifest.Dimensions {
		stateS3Path = stateS3Path + dimension + "_" + tofuguStruct.ParsedDimensions[dimension] + "/"
	}
	tofuguStruct.StateS3Path = stateS3Path + tofuguStruct.TofiName + ".tfstate"
}

func (tofuguStruct *Tofugu) GetDimData(dimensionKey string, dimensionValue string) map[string]interface{} {
	var dimensionJsonMap map[string]interface{}
	var dimensionJsonBytes []byte
	var err error

	if tofuguStruct.ToasterUrl == "" {
		inventroyJsonPath := tofuguStruct.InventoryPath + "/" + dimensionKey + "/" + dimensionValue + ".json"
		dimensionJsonBytes, err = os.ReadFile(inventroyJsonPath)
		if err != nil {
			log.Fatal("Error when opening file: ", err)
		}
	} else {
		resp, err := http.Get(tofuguStruct.ToasterUrl + "/dimension/" + tofuguStruct.OrgName + "/" + dimensionKey + ":" + dimensionValue)

		if err != nil {
			log.Fatalf("request to Toaster Failed: %s", err)
		}
		defer resp.Body.Close()
		dimensionJsonBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("reading Toaster response failed: %s", err)
		}
	}
	json.Unmarshal(dimensionJsonBytes, &dimensionJsonMap)
	return dimensionJsonMap
}
