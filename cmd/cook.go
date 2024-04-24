package cmd

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/alt-dima/tofugu/utils"
	"github.com/spf13/cobra"
)

// cookCmd represents the cook command
var cookCmd = &cobra.Command{
	Use:   "cook",
	Short: "Execute OpenTofu",
	Args:  cobra.MinimumNArgs(1),
	Long:  `Execute OpenTofu with generated config from inventory and parameters after --`,
	Run: func(cmd *cobra.Command, args []string) {
		//Creating signal to be handled and send to the child tofu/terraform
		sigs := make(chan os.Signal, 2)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		var err error

		// Creating Tofug shared structure and filling with values
		tofuguStruct := &utils.Tofugu{}

		tofuguStruct.TofiName, _ = cmd.Flags().GetString("tofi")
		tofuguStruct.OrgName, _ = cmd.Flags().GetString("org")
		tofuguStruct.DimensionsFlags, _ = cmd.Flags().GetStringSlice("dimension")
		tofuguStruct.TofiPath, _ = filepath.Abs(tofuguStruct.GetStringFromViperByOrgOrDefault("tofies_path") + "/" + tofuguStruct.OrgName + "/" + tofuguStruct.TofiName)
		tofuguStruct.SharedModulesPath, _ = filepath.Abs(tofuguStruct.GetStringFromViperByOrgOrDefault("shared_modules_path"))
		tofuguStruct.InventoryPath, _ = filepath.Abs(tofuguStruct.GetStringFromViperByOrgOrDefault("inventory_path") + "/" + tofuguStruct.OrgName)

		tofuguStruct.ParseTofiManifest("tofi_manifest.json")
		tofuguStruct.ParseDimensions()

		tofuguStruct.SetupStateS3Path()

		tofuguStruct.PrepareTemp()

		tofuguStruct.GenerateVarsByDims()
		tofuguStruct.GenerateVarsByEnvVars()

		//Local variables for child execution
		stateS3Region := tofuguStruct.GetStringFromViperByOrgOrDefault("s3_bucket_region")
		stateS3Name := tofuguStruct.GetStringFromViperByOrgOrDefault("s3_bucket_name")
		forceCleanTempDir, _ := cmd.Flags().GetBool("clean")
		cmdArgs := args
		if args[0] == "init" {
			cmdArgs = append(cmdArgs, "-backend-config=bucket="+stateS3Name)
			cmdArgs = append(cmdArgs, "-backend-config=key="+tofuguStruct.StateS3Path)
			cmdArgs = append(cmdArgs, "-backend-config=region="+stateS3Region)
		}
		cmdToExec := tofuguStruct.GetStringFromViperByOrgOrDefault("cmd_to_exec")

		// Starting child and Waiting for it to finish, passing signals to it
		log.Println("TofuGu starting cooking: " + cmdToExec + " " + strings.Join(cmdArgs, " "))
		execChildCommand := exec.Command(cmdToExec, cmdArgs...)
		execChildCommand.Dir = tofuguStruct.CmdWorkTempDir
		execChildCommand.Env = os.Environ()
		execChildCommand.Stdin = os.Stdin
		execChildCommand.Stdout = os.Stdout
		execChildCommand.Stderr = os.Stderr
		err = execChildCommand.Start()
		if err != nil {
			log.Fatalf("cmd.Start() failed with %s\n", err)
		}

		go func() {
			sig := <-sigs
			log.Println("Got singnal +" + sig.String())
			execChildCommand.Process.Signal(sig)
		}()

		err = execChildCommand.Wait()
		exitCodeFinal := 0
		if err != nil && execChildCommand.ProcessState.ExitCode() < 0 {
			exitCodeFinal = 1
			log.Println(cmdToExec + " failed " + err.Error())
		} else if execChildCommand.ProcessState.ExitCode() == 143 {
			exitCodeFinal = 0
		} else {
			exitCodeFinal = execChildCommand.ProcessState.ExitCode()
		}

		if (exitCodeFinal == 0 && (args[0] == "apply" || args[0] == "destroy")) || forceCleanTempDir {
			os.RemoveAll(tofuguStruct.CmdWorkTempDir)
			log.Println("TofuGu removed tofi temp dir: " + tofuguStruct.CmdWorkTempDir)
		}

		log.Printf("TofuGu: %v finished with code %v", cmdToExec, exitCodeFinal)
		os.Exit(exitCodeFinal)
	},
}

func init() {
	rootCmd.AddCommand(cookCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cookCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	cookCmd.Flags().StringSliceP("dimension", "d", []string{}, "specify dimensions from invetory like dim:name")
	//viper.BindPFlag("account", cookCmd.Flags().Lookup("account"))
	cookCmd.Flags().StringP("tofi", "t", "", "specify tofu unit")
	//viper.BindPFlag("tofi", cookCmd.Flags().Lookup("tofi"))
	cookCmd.Flags().StringP("org", "o", "", "specify org")
	cookCmd.Flags().BoolP("clean", "c", false, "remove tmp after execution")
	//viper.BindPFlag("org", cookCmd.Flags().Lookup("org"))
	cookCmd.MarkFlagRequired("tofi")
	cookCmd.MarkFlagRequired("org")
}
