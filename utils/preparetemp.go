package utils

import (
	"log"
	"os"
	"os/exec"
)

func (tofuguStruct *Tofugu) PrepareTemp() {
	tmpFolderNameSuffix := tofuguStruct.OrgName + tofuguStruct.StateS3Path + tofuguStruct.TofiName
	cmdTempDirFullPath := os.TempDir() + "/tofugu-" + GetMD5Hash(tmpFolderNameSuffix)

	command := exec.Command("rsync", "-a", "--delete", "--exclude=.terraform*", "--exclude=tofi_manifest.json", tofuguStruct.TofiPath+"/.", cmdTempDirFullPath)
	output, err := command.CombinedOutput()
	if err != nil {
		os.RemoveAll(cmdTempDirFullPath)
		log.Printf("failed %s", output)
		log.Fatalf("failed to rsync tofi to tempdir %s\n", err)
	}

	command = exec.Command("ln", "-sf", tofuguStruct.SharedModulesPath, cmdTempDirFullPath)
	output, err = command.CombinedOutput()
	if err != nil {
		os.RemoveAll(cmdTempDirFullPath)
		log.Printf("failed %s", output)
		log.Fatalf("failed symlink shared_modules to tempdir %s\n", err)
	}

	tofuguStruct.CmdWorkTempDir = cmdTempDirFullPath
	log.Println("TofuGu prepared tofi in temp dir: " + tofuguStruct.CmdWorkTempDir)
}
