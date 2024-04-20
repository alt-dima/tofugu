package cmd

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

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

		for key, value := range viper.GetViper().AllSettings() {
			log.Printf("key %v val %v", key, value)
		}

		cmdToExec := viper.GetString("defaults.cmd_to_exec")
		// cmdTempDir, err := os.MkdirTemp("", "tofugu")
		// if err != nil {
		// 	log.Fatalf("MkdirTemp failed with %s\n", err)
		// }

		// command := exec.Command(mountCmd, mountArgs...)
		// output, err := command.CombinedOutput()

		tofiName, _ := cmd.Flags().GetString("tofi")
		orgName, _ := cmd.Flags().GetString("org")
		tofiPath := viper.GetString("defaults.tofies_path") + "/" + orgName + "/" + tofiName

		log.Println("ToFuGu starting OpenTofu with args: " + strings.Join(args, " "))
		execChildCommand := exec.Command(cmdToExec, args...)
		execChildCommand.Dir = tofiPath
		execChildCommand.Env = os.Environ()
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
