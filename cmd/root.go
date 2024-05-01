package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var Verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tofugu",
	Short: "Tofies manager and generator for OpenTofu",
	Long:  `Generates additinal vars files based on inventory and executes opentofu`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func SetVersionInfo(version string) {
	rootCmd.Version = version
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tofugu)")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	log.Println("TofuGu version:" + rootCmd.Version)

	viper.SetDefault("defaults.inventory_path", "examples/inventory")
	viper.SetDefault("defaults.shared_modules_path", "examples/tofies/shared-modules")
	viper.SetDefault("defaults.tofies_path", "examples/tofies")
	viper.SetDefault("defaults.cmd_to_exec", "tofu")

	viper.SetConfigType("yaml")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".tofugu")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Println("TofuGu using config file:", viper.ConfigFileUsed())
	} else {
		log.Println(err.Error())
	}
}
