package utils

import (
	"log"
	"os"
	"strings"
)

func (tofuguStruct *Tofugu) ParseDimensions() {
	parsedDimArgs := parseDimArgs(tofuguStruct.DimensionsFlags)

	for _, dimension := range tofuguStruct.TofiManifest.Dimensions {
		if _, ok := parsedDimArgs[dimension]; !ok {
			log.Println("dimension " + dimension + " not passed with -d arg")
			os.Exit(1)
		}
	}

	tofuguStruct.ParsedDimensions = parsedDimArgs
}

func parseDimArgs(dimensionsArgs []string) map[string]string {
	parsedDimArgs := make(map[string]string)
	for _, dimension := range dimensionsArgs {
		dimensionSlice := strings.SplitN(dimension, ":", 2)
		parsedDimArgs[dimensionSlice[0]] = dimensionSlice[1]
	}
	return parsedDimArgs
}
