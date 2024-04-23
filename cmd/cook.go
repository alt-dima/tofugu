package cmd

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/alt-dima/tofugu/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// cookCmd represents the cook command
var cookCmd = &cobra.Command{
	Use:   "cook",
	Short: "Execute OpenTofu",
	Args:  cobra.MinimumNArgs(1),
	Long:  `Execute OpenTofu with generated config from inventory and parameters after --`,
	Run: func(cmd *cobra.Command, args []string) {
		sigs := make(chan os.Signal, 2)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

		//for key, value := range viper.GetViper().AllSettings() {
		//	log.Printf("key %v val %v", key, value)
		//}

		tofiName, _ := cmd.Flags().GetString("tofi")
		orgName, _ := cmd.Flags().GetString("org")
		dimensionsArgs, _ := cmd.Flags().GetStringSlice("dimension")

		cmdToExec := utils.GetConfigFromViperString("cmd_to_exec", orgName)
		currentDir, _ := os.Getwd()
		tofiPath := currentDir + "/" + utils.GetConfigFromViperString("tofies_path", orgName) + "/" + orgName + "/" + tofiName
		manifest := utils.ParseTofiManifest(tofiPath + "/tofi_manifest.json")

		//log.Println(manifest.Dimensions)
		parsedDimensions := utils.ParseDimensions(manifest.Dimensions, dimensionsArgs)

		var stateS3Path string
		if !viper.IsSet(orgName + ".s3_bucket_name") {
			stateS3Path = stateS3Path + "org_" + orgName + "/"
		}
		for _, dimension := range manifest.Dimensions {
			stateS3Path = stateS3Path + dimension + "_" + parsedDimensions[dimension] + "/"
		}
		stateS3Path = stateS3Path + tofiName + ".tfstate"
		stateS3Region := utils.GetConfigFromViperString("s3_bucket_region", orgName)
		stateS3Name := utils.GetConfigFromViperString("s3_bucket_name", orgName)

		cmdArgs := args
		if args[0] == "init" {
			cmdArgs = append(cmdArgs, "-backend-config=bucket="+stateS3Name)
			cmdArgs = append(cmdArgs, "-backend-config=key="+stateS3Path)
			cmdArgs = append(cmdArgs, "-backend-config=region="+stateS3Region)
		}

		cmdWorkTempDir := utils.PrepareTemp(tofiPath, currentDir+"/"+utils.GetConfigFromViperString("shared_modules_path", orgName), orgName+stateS3Path+tofiName)

		utils.GenerateVarsByDims(parsedDimensions, cmdWorkTempDir, currentDir+"/"+utils.GetConfigFromViperString("inventory_path", orgName)+"/"+orgName)
		utils.GenerateVarsByEnvVars(cmdWorkTempDir)

		log.Println("TofuGu starting cooking: " + cmdToExec + " " + strings.Join(cmdArgs, " "))
		execChildCommand := exec.Command(cmdToExec, cmdArgs...)
		execChildCommand.Dir = cmdWorkTempDir
		execChildCommand.Env = os.Environ()
		execChildCommand.Stdin = os.Stdin
		execChildCommand.Stdout = os.Stdout
		execChildCommand.Stderr = os.Stderr
		err := execChildCommand.Start()
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
			log.Println("OpenTofu failed " + err.Error())
		} else if execChildCommand.ProcessState.ExitCode() == 143 {
			exitCodeFinal = 0
		} else {
			exitCodeFinal = execChildCommand.ProcessState.ExitCode()
		}

		if args[0] == "apply" || args[0] == "destroy" {
			os.RemoveAll(cmdWorkTempDir)
		}

		log.Printf("OpenTofu finished with code %v", exitCodeFinal)
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
	//viper.BindPFlag("org", cookCmd.Flags().Lookup("org"))
	cookCmd.MarkFlagRequired("tofi")
	cookCmd.MarkFlagRequired("org")
}
