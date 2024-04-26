package utils

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

func (tofuguStruct *Tofugu) GenerateVarsByDims() {
	for dimKey, dimValue := range tofuguStruct.ParsedDimensions {
		dimensionJsonMap := tofuguStruct.GetDimData(dimKey, dimValue)

		targetAutoTfvarMap := map[string]interface{}{
			"tofugu_" + dimKey + "_data": dimensionJsonMap,
			"tofugu_" + dimKey + "_name": dimValue,
		}

		writeTfvarsMaps(targetAutoTfvarMap, dimKey, tofuguStruct.CmdWorkTempDir)
		log.Println("TofuGu generated tfvars for dimension: " + dimKey)

	}
}

func (tofuguStruct *Tofugu) GenerateVarsByEnvVars() {
	targetAutoTfvarMap := make(map[string]interface{})

	for _, envVar := range os.Environ() {
		if strings.HasPrefix(envVar, "tofugu_envvar_") {
			envVarList := strings.SplitN(envVar, "=", 2)
			targetAutoTfvarMap[envVarList[0]] = envVarList[1]
		}
	}

	if len(targetAutoTfvarMap) > 0 {
		writeTfvarsMaps(targetAutoTfvarMap, "envivars", tofuguStruct.CmdWorkTempDir)
		log.Println("TofuGu generated tfvars for env variables")
	}
}

func writeTfvarsMaps(targetAutoTfvarMap map[string]interface{}, fileName string, cmdWorkTempDir string) {
	targetVarsTfPath := cmdWorkTempDir + "/tofugu_" + fileName + "_vars.tf.json"
	targetAutoTfvarsPath := cmdWorkTempDir + "/tofugu_" + fileName + ".auto.tfvars.json"

	targetVarsTfMap := make(map[string]interface{})

	for key, value := range targetAutoTfvarMap {
		switch value.(type) {
		case string:
			targetVarsTfMap[key] = map[string]string{"type": "string"}
		default:
			targetVarsTfMap[key] = map[string]interface{}{}
		}

	}

	targetVarsTfMapFull := map[string]interface{}{
		"variable": targetVarsTfMap,
	}

	marshalJsonAndWrite(targetVarsTfMapFull, targetVarsTfPath)
	marshalJsonAndWrite(targetAutoTfvarMap, targetAutoTfvarsPath)
}

func marshalJsonAndWrite(jsonMap map[string]interface{}, jsonPath string) {
	targetAutoTfvarMapBytes, _ := json.Marshal(jsonMap)
	err := os.WriteFile(jsonPath, targetAutoTfvarMapBytes, os.ModePerm)
	if err != nil {
		log.Fatal("error writing file: ", err)
	}
}
