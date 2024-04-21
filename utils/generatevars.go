package utils

import (
	"encoding/json"
	"log"
	"os"
)

func GenerateVarsByDims(parsedDimensions map[string]string, cmdWorkTempDir string, inventoryPath string) {
	for dimKey, dimValue := range parsedDimensions {
		var inventroyJsonMap interface{}

		inventroyJsonPath := inventoryPath + "/" + dimKey + "/" + dimValue + ".json"
		targetVarsTfPath := cmdWorkTempDir + "/tofugu_" + dimKey + "_vars.tf.json"
		targetAutoTfvarsPath := cmdWorkTempDir + "/tofugu_" + dimKey + ".auto.tfvars.json"

		inventroyJsonBytes, err := os.ReadFile(inventroyJsonPath)
		if err != nil {
			log.Fatal("Error when opening file: ", err)
		}
		json.Unmarshal(inventroyJsonBytes, &inventroyJsonMap)

		targetVarsTfMap := map[string]interface{}{
			"variable": map[string]interface{}{
				"tofugu_" + dimKey + "_manifest": map[string]interface{}{},
				"tofugu_" + dimKey + "_name":     map[string]string{"type": "string"},
			},
		}

		targetAutoTfvarMap := map[string]interface{}{
			"tofugu_" + dimKey + "_manifest": inventroyJsonMap,
			"tofugu_" + dimKey + "_name":     dimValue,
		}

		targetVarsTfJson, _ := json.Marshal(targetVarsTfMap)
		os.WriteFile(targetVarsTfPath, targetVarsTfJson, os.ModePerm)

		targetAutoTfvarfJson, _ := json.Marshal(targetAutoTfvarMap)
		os.WriteFile(targetAutoTfvarsPath, targetAutoTfvarfJson, os.ModePerm)
	}
}
