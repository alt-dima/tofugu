package utils

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

func (tofuguStruct *Tofugu) GenerateVarsByDims() {
	for dimKey, dimValue := range tofuguStruct.ParsedDimensions {
		dimensionJsonMap := tofuguStruct.GetDimData(dimKey, dimValue, false)

		targetAutoTfvarMap := map[string]interface{}{
			"tofugu_" + dimKey + "_data": dimensionJsonMap,
			"tofugu_" + dimKey + "_name": dimValue,
		}

		writeTfvarsMaps(targetAutoTfvarMap, dimKey, tofuguStruct.CmdWorkTempDir)
		log.Println("TofuGu attached dimension in var.tofugu_" + dimKey + "_data and var.tofugu_" + dimKey + "_name")

	}
}

func (tofuguStruct *Tofugu) GenerateVarsByDimOptional(optionType string) {
	for dimKey := range tofuguStruct.ParsedDimensions {
		dimensionJsonMap := tofuguStruct.GetDimData(dimKey, "dim_"+optionType, true)
		if len(dimensionJsonMap) > 0 {
			targetAutoTfvarMap := map[string]interface{}{
				"tofugu_" + dimKey + "_" + optionType: dimensionJsonMap,
			}

			writeTfvarsMaps(targetAutoTfvarMap, dimKey+"_"+optionType, tofuguStruct.CmdWorkTempDir)
			log.Println("TofuGu attached " + optionType + " in var.tofugu_" + dimKey + "_" + optionType)
		}
	}
}

func (tofuguStruct *Tofugu) GenerateVarsByEnvVars() {
	targetAutoTfvarMap := make(map[string]interface{})

	for _, envVar := range os.Environ() {
		if strings.HasPrefix(envVar, "tofugu_envvar_") {
			envVarList := strings.SplitN(envVar, "=", 2)
			targetAutoTfvarMap[envVarList[0]] = envVarList[1]
			log.Println("TofuGu attached env variable in var." + envVarList[0])
		}
	}

	if len(targetAutoTfvarMap) > 0 {
		writeTfvarsMaps(targetAutoTfvarMap, "envivars", tofuguStruct.CmdWorkTempDir)
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
	targetAutoTfvarMapBytes, err := json.Marshal(jsonMap)
	if err != nil {
		log.Fatal("tofugu failed to marshal json: ", err)
	}
	err = os.WriteFile(jsonPath, targetAutoTfvarMapBytes, os.ModePerm)
	if err != nil {
		log.Fatal("tofugu error writing file: ", err)
	}
}
