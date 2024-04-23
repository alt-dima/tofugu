package utils

import (
	"log"
	"os"
	"os/exec"
)

func PrepareTemp(tofiPath string, sharedModulesPath string, tmpFolderName string) string {
	cmdTempDir := os.TempDir() + "/tofugu-" + GetMD5Hash(tmpFolderName)

	command := exec.Command("rsync", "-a", "--delete", "--exclude=.terraform*", "--exclude=tofi_manifest.json", tofiPath+"/.", cmdTempDir)
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

	log.Println("TofuGu prepared tofi in temp dir: " + cmdTempDir)
	return cmdTempDir
}
