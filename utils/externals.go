package utils

import (
	"bytes"
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

func (tofuguStruct *Tofugu) SetupBackendConfig() []string {
	var backendFinalConfig []string

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
	for param, value := range backendTofuguConfig {
		replacedVar := strings.Replace(value.(string), "$tofugu_state_path", tofuguStruct.StateS3Path, 1)
		backendFinalConfig = append(backendFinalConfig, "-backend-config="+param+"="+replacedVar)
	}

	return backendFinalConfig
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

func (tofuguStruct *Tofugu) SendHistoryData(cmdToExec string, cmdArgs []string, exitCodeFinal int) {
	if tofuguStruct.ToasterUrl != "" {

		var historyData HistoryStruct
		historyData.CmdToExec = cmdToExec
		historyData.CmdMainArg = cmdArgs[0]
		if len(cmdArgs) > 1 {
			historyData.CmdArgs = cmdArgs[1:]
		}
		historyData.ExitCode = exitCodeFinal
		historyData.Dimensions = tofuguStruct.ParsedDimensions

		byteStream := new(bytes.Buffer)
		err := json.NewEncoder(byteStream).Encode(historyData)
		if err != nil {
			log.Printf("tofugu: failed to prepare json data: %s", err)
		}

		resp, err := http.Post(tofuguStruct.ToasterUrl+"/api/history/"+tofuguStruct.OrgName+"/"+tofuguStruct.Workspace+"/"+tofuguStruct.TofiName, "application/json; charset=UTF-8", byteStream)
		if err != nil {
			log.Printf("tofugu toaster: history post request failed: %s", err)
		} else if resp.StatusCode != 200 {
			resp.Body.Close()
			log.Printf("tofugu toaster: history post request failed with response: %v", resp.StatusCode)
		}
		defer resp.Body.Close()
	}
}
