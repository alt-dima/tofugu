package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

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

func (tofuguStruct *Tofugu) GetObjectFromViperByOrgOrDefault(keyName string) map[string]any {
	if viper.IsSet(tofuguStruct.OrgName + "." + keyName) {
		return viper.GetStringMap(tofuguStruct.OrgName + "." + keyName)
	} else {
		return viper.GetStringMap("defaults." + keyName)
	}
}

func (tofuguStruct *Tofugu) SetupBackendConfig() map[string]interface{} {
	var stateS3Path string
	if !viper.IsSet(tofuguStruct.OrgName + ".backend") {
		stateS3Path = stateS3Path + "org_" + tofuguStruct.OrgName + "/"
	}

	for _, dimension := range tofuguStruct.TofiManifest.Dimensions {
		stateS3Path = stateS3Path + dimension + "_" + tofuguStruct.ParsedDimensions[dimension] + "/"
	}
	tofuguStruct.StateS3Path = stateS3Path + tofuguStruct.TofiName + ".tfstate"

	backendTofuguConfig := tofuguStruct.GetObjectFromViperByOrgOrDefault("backend")
	if len(backendTofuguConfig) == 0 {
		log.Println("Tofugu: no backend config provied!")
	}

	var backendTofuguConfigMap = make(map[string]interface{}, len(backendTofuguConfig))
	for param, value := range backendTofuguConfig {
		backendTofuguConfigMap[param] = strings.Replace(value.(string), "$tofugu_state_path", tofuguStruct.StateS3Path, 1)
	}

	return backendTofuguConfigMap
}

func (tofuguStruct *Tofugu) GetDimData(dimensionKey string, dimensionValue string, skipOnNotFound bool) map[string]interface{} {
	var dimensionJsonMap map[string]interface{}

	if tofuguStruct.ToasterUrl == "" {
		inventroyJsonPath := tofuguStruct.InventoryPath + "/" + dimensionKey + "/" + dimensionValue + ".json"
		dimensionJsonBytes, err := os.ReadFile(inventroyJsonPath)
		if err != nil {
			if os.IsNotExist(err) && skipOnNotFound {
				log.Println("TofuGu inventory files: Optional dimension " + tofuguStruct.OrgName + "/" + dimensionKey + "/" + dimensionValue + " not found, skipping")
				return dimensionJsonMap
			}
			log.Fatal("tofugu inventory files: error when opening dim file: ", err.Error())
		}
		err = json.Unmarshal(dimensionJsonBytes, &dimensionJsonMap)
		if err != nil {
			log.Fatal("tofugu error during Unmarshal(): ", err)
		}
	} else {
		resp, err := http.Get(tofuguStruct.ToasterUrl + "/api/dimension/" + tofuguStruct.OrgName + "/" + dimensionKey + "/" + dimensionValue + "?workspace=" + tofuguStruct.Workspace + "&fallbacktomaster=true")
		if err != nil {
			log.Fatalf("tofugu toaster: request Failed: %s", err)
		} else if resp.StatusCode == 404 {
			resp.Body.Close()
			if skipOnNotFound {
				log.Println("TofuGu Toaster: optional dimension " + tofuguStruct.OrgName + "/" + dimensionKey + "/" + dimensionValue + " not found, skipping")
				return dimensionJsonMap
			} else {
				log.Fatalln("tofugu toaster: dimension " + tofuguStruct.OrgName + "/" + dimensionKey + "/" + dimensionValue + " not found")
			}
		} else if resp.StatusCode != 200 {
			resp.Body.Close()
			log.Fatalf("tofugu toaster: request "+tofuguStruct.OrgName+"/"+dimensionKey+"/"+dimensionValue+"?workspace="+tofuguStruct.Workspace+" failed with response: %v", resp.StatusCode)
		}
		defer resp.Body.Close()

		dimensionJsonBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("tofugu toaster: reading body response failed: %s", err)
		}

		var toasterResponse ToasterResponse
		err = json.Unmarshal(dimensionJsonBytes, &toasterResponse)
		if err != nil {
			log.Fatal("tofugu toaster: error during unmarshal json response: ", err)
		}

		if len(toasterResponse.Dimensions) != 1 {
			log.Fatalf("tofugu toaster: should be only one dimension in response")
		}
		if toasterResponse.Error != "" {
			log.Println("TofuGu Toaster: " + toasterResponse.Error)
		}
		dimensionJsonMap = toasterResponse.Dimensions[0].DimData

	}

	return dimensionJsonMap
}
