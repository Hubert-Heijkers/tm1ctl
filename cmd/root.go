package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tm1ctl",
	Short: "TM1 v12 control utility",
	Long:  `The TM1 v12 control utility allows you to manage your TM1 v12 service from the command line.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tm1ctl.json)")
	rootCmd.PersistentFlags().String("output", "", "set the output format for this request, either 'table' or 'json' (defaults to output_format config)")
	viper.BindPFlag("output-format", rootCmd.PersistentFlags().Lookup("output"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Defaults
	hostmap := make(map[string]any)
	hostmap["service_root_url"] = "http://localhost:4444"
	hosts := make(map[string]any)
	hosts["local"] = hostmap
	viper.SetDefault("hosts", hosts)
	viper.SetDefault("host", "local")

	viper.SetDefault("output-format", "table")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			// Error reading specified config file
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".tm1ctl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("json")
		viper.SetConfigName(".tm1ctl")
		if err := viper.ReadInConfig(); err != nil {
			// Create config if it doesn't exist
			viper.SafeWriteConfigAs(filepath.Join(home, ".tm1ctl.json"))
		}
	}
}
