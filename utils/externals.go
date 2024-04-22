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

func GetConfigFromViperString(keyName string, orgName string) string {
	if viper.IsSet(orgName + "." + keyName) {
		return viper.GetString(orgName + "." + keyName)
	} else {
		return viper.GetString("defaults." + keyName)
	}
}
