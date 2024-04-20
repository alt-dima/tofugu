package utils

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"os"
	"os/exec"
)

func PrepareTemp(tofiPath string, sharedModulesPath string, tmpFolderName string) string {
	cmdTempDir := os.TempDir() + "/tofugu" + getMD5Hash(tmpFolderName)

	command := exec.Command("cp", "-R", tofiPath+"/.", cmdTempDir)
	output, err := command.CombinedOutput()
	if err != nil {
		os.RemoveAll(cmdTempDir)
		log.Printf("failed %s", output)
		log.Fatalf("failed to copit tofi to tempdir %s\n", err)
	}

	command = exec.Command("ln", "-sf", sharedModulesPath, cmdTempDir)
	output, err = command.CombinedOutput()
	if err != nil {
		os.RemoveAll(cmdTempDir)
		log.Printf("failed %s", output)
		log.Fatalf("failed to copit tofi to tempdir %s\n", err)
	}

	return cmdTempDir
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
