package utils

import (
	"encoding/json"
	"log"
	"os"
)

type tofiManifestStruct struct {
	Dimensions []string
}

func ParseTofiManifest(tofiManifestPath string) tofiManifestStruct {
	// Let's first read the `config.json` file
	content, err := os.ReadFile(tofiManifestPath)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	// Now let's unmarshall the data into `payload`
	var tofiManifest tofiManifestStruct
	err = json.Unmarshal(content, &tofiManifest)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	return tofiManifest
}
