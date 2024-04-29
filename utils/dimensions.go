package utils

import (
	"log"
	"strings"
)

func (tofuguStruct *Tofugu) ParseDimensions() {
	parsedDimArgs := parseDimArgs(tofuguStruct.DimensionsFlags)

	for _, dimension := range tofuguStruct.TofiManifest.Dimensions {
		if _, ok := parsedDimArgs[dimension]; !ok {
			log.Fatalln("dimension " + dimension + " not passed with -d arg")
		}
	}

	tofuguStruct.ParsedDimensions = parsedDimArgs
}

func parseDimArgs(dimensionsArgs []string) map[string]string {
	parsedDimArgs := make(map[string]string)
	for _, dimension := range dimensionsArgs {
		dimensionSlice := strings.SplitN(dimension, ":", 2)
		if strings.HasPrefix(dimensionSlice[1], "dim_") {
			log.Fatalln("dimension " + dimension + " with dim_ prefix can't be passed with -d arg")
		}
		parsedDimArgs[dimensionSlice[0]] = dimensionSlice[1]
	}
	return parsedDimArgs
}
