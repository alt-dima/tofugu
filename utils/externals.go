package utils

import (
	"crypto/md5"
	"encoding/hex"

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

func (tofuguStruct *Tofugu) SetupStateS3Path() {
	var stateS3Path string
	if !viper.IsSet(tofuguStruct.OrgName + ".s3_bucket_name") {
		stateS3Path = stateS3Path + "org_" + tofuguStruct.OrgName + "/"
	}
	for _, dimension := range tofuguStruct.TofiManifest.Dimensions {
		stateS3Path = stateS3Path + dimension + "_" + tofuguStruct.ParsedDimensions[dimension] + "/"
	}
	tofuguStruct.StateS3Path = stateS3Path + tofuguStruct.TofiName + ".tfstate"
}
