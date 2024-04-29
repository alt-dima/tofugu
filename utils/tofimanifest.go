package utils

import (
	"encoding/json"
	"log"
	"os"
)

func (tofuguStruct *Tofugu) ParseTofiManifest(tofiManifestFileName string) {
	tofiManifestPath := tofuguStruct.TofiPath + "/" + tofiManifestFileName
	// Let's first read the `config.json` file
	content, err := os.ReadFile(tofiManifestPath)
	if err != nil {
		log.Fatal("tofugu error when opening file: ", err)
	}

	// Now let's unmarshall the data into `payload`
	var tofiManifest tofiManifestStruct
	err = json.Unmarshal(content, &tofiManifest)
	if err != nil {
		log.Fatal("tofugu error during Unmarshal(): ", err)
	}

	tofuguStruct.TofiManifest = tofiManifest
	log.Println("TofuGu loaded tofi manifest: " + tofiManifestPath)
}
