package cmd

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
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

		cmdToExec := "tofu"

		log.Println("ToFuGu starting OpenTofu with args: " + strings.Join(args, " "))
		execChildCommand := exec.Command(cmdToExec, args...)
		execChildCommand.Env = os.Environ()
		execChildCommand.Stdout = os.Stdout
		execChildCommand.Stderr = os.Stderr
		err := execChildCommand.Start()
		if err != nil {
			log.Fatalf("cmd.Start() failed with %s\n", err)
		}

		go func() {
			sig := <-sigs
			log.Println("Got singla +" + sig.String())
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
	cookCmd.Flags().StringP("account", "a", "", "specify account from invetory")
	cookCmd.MarkFlagRequired("account")
}
