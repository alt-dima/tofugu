package utils

import (
	"log"
	"os"
	"path/filepath"

	"github.com/otiai10/copy"
)

func (tofuguStruct *Tofugu) PrepareTemp() {
	if tofuguStruct.StateS3Path == "" {
		log.Fatalf("StateS3Path is empty \n")
	}

	tmpFolderNameSuffix := tofuguStruct.OrgName + tofuguStruct.StateS3Path + tofuguStruct.TofiName
	cmdTempDirFullPath := os.TempDir() + "/tofugu-" + GetMD5Hash(tmpFolderNameSuffix)

	// Create temp directory if it doesn't exist
	if err := os.MkdirAll(cmdTempDirFullPath, 0755); err != nil {
		log.Fatalf("failed to create temp directory: %v\n", err)
	}

	// Copy options to exclude certain files/directories
	opt := copy.Options{
		Skip: func(info os.FileInfo, src string, dest string) (bool, error) {
			base := filepath.Base(src)
			return base == ".terraform" || base == "tofi_manifest.json", nil
		},
	}

	// Copy the tofi directory to temp directory
	if err := copy.Copy(tofuguStruct.TofiPath, cmdTempDirFullPath, opt); err != nil {
		os.RemoveAll(cmdTempDirFullPath)
		log.Fatalf("failed to copy tofi to tempdir: %v\n", err)
	}

	if tofuguStruct.SharedModulesPath != "" {
		// Remove existing symlink if it exists
		sharedModulesLink := filepath.Join(cmdTempDirFullPath, "shared-modules")
		os.Remove(sharedModulesLink) // Ignore error as file might not exist

		// Create new symlink
		if err := os.Symlink(tofuguStruct.SharedModulesPath, sharedModulesLink); err != nil {
			os.RemoveAll(cmdTempDirFullPath)
			log.Fatalf("failed to create symlink for shared_modules: %v\n", err)
		}
		log.Println("TofuGu symlinked shared_modules to tempdir : " + tofuguStruct.SharedModulesPath)
	}

	tofuguStruct.CmdWorkTempDir = cmdTempDirFullPath
	log.Println("TofuGu prepared tofi in temp dir: " + tofuguStruct.CmdWorkTempDir)
}
