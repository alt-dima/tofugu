package utils

import (
	"log"
	"os"
	"strings"
)

func ParseDimensions(dimensionsManifest []string, dimensionsArgs []string) map[string]string {
	parsedDimArgs := parseDimArgs(dimensionsArgs)

	for _, dimension := range dimensionsManifest {
		if _, ok := parsedDimArgs[dimension]; !ok {
			log.Println("dimension " + dimension + " not passed with -d arg")
			os.Exit(1)
		}
	}
	return parsedDimArgs
}

func parseDimArgs(dimensionsArgs []string) map[string]string {
	parsedDimArgs := make(map[string]string)
	for _, dimension := range dimensionsArgs {
		dimensionSlice := strings.Split(dimension, ":")
		parsedDimArgs[dimensionSlice[0]] = dimensionSlice[1]
	}
	return parsedDimArgs
}
