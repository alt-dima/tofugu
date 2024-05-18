package utils

import (
	"log"
	"os"
	"os/exec"
)

func (tofuguStruct *Tofugu) PrepareTemp() {
	if tofuguStruct.StateS3Path == "" {
		log.Fatalf("StateS3Path is empty \n")
	}

	tmpFolderNameSuffix := tofuguStruct.OrgName + tofuguStruct.StateS3Path + tofuguStruct.TofiName
	cmdTempDirFullPath := os.TempDir() + "/tofugu-" + GetMD5Hash(tmpFolderNameSuffix)

	command := exec.Command("rsync", "-a", "--delete", "--exclude=.terraform*", "--exclude=tofi_manifest.json", tofuguStruct.TofiPath+"/.", cmdTempDirFullPath)
	output, err := command.CombinedOutput()
	if err != nil {
		os.RemoveAll(cmdTempDirFullPath)
		log.Printf("failed %s", output)
		log.Fatalf("failed to rsync tofi to tempdir %s\n", err)
	}

	if tofuguStruct.SharedModulesPath != "" {
		command = exec.Command("ln", "-sf", tofuguStruct.SharedModulesPath, cmdTempDirFullPath)
		output, err = command.CombinedOutput()
		if err != nil {
			os.RemoveAll(cmdTempDirFullPath)
			log.Printf("failed %s", output)
			log.Fatalf("failed symlink shared_modules to tempdir %s\n", err)
		}
		log.Println("TofuGu symlinked shared_modules to tempdir : " + tofuguStruct.SharedModulesPath)
	}

	tofuguStruct.CmdWorkTempDir = cmdTempDirFullPath
	log.Println("TofuGu prepared tofi in temp dir: " + tofuguStruct.CmdWorkTempDir)
}
