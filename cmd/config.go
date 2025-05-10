package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var allowedConfigKeys = map[string]bool{
	"service-root-url":   true,
	"root-client-id":     true,
	"root-client-secret": true,
	"service-instance":   true,
	"database":           true,
	"user":               true,
	"password":           true,
	"output-format":      true,
}

var allowedOutputFormats = map[string]bool{
	"json":  true,
	"table": true,
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage tm1ctl configuration",
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration values",
	Run: func(cmd *cobra.Command, args []string) {
		for key := range allowedConfigKeys {
			val := viper.Get(key)
			if val != nil && val != "" {
				fmt.Printf("%s = %v\n", key, val)
			}
		}
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]

		if !allowedConfigKeys[key] {
			fmt.Printf("Error: '%s' is not a recognized configuration key.\n", key)
			fmt.Println("Allowed keys are:")
			for key = range allowedConfigKeys {
				fmt.Println(" -", key)
			}
			return
		}

		val := viper.Get(key)
		fmt.Printf("%s = %v\n", key, val)
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set and save a configuration value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		if !allowedConfigKeys[key] {
			fmt.Printf("Error: '%s' is not a recognized configuration key.\n", key)
			fmt.Println("Allowed keys are:")
			for key := range allowedConfigKeys {
				fmt.Println(" -", key)
			}
			return
		}

		if key == "output-format" && !allowedOutputFormats[value] {
			fmt.Printf("Error: '%s' is not a recognized output format.\n", value)
			fmt.Println("Supported output formats are:")
			for value := range allowedOutputFormats {
				fmt.Println(" -", value)
			}
			return
		}

		viper.Set(key, value)
		if err := viper.WriteConfig(); err != nil {
			fmt.Printf("Failed to save config: %v\n", err)
		} else {
			fmt.Printf("%s set to %s\n", key, value)
		}
	},
}

func init() {
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configSetCmd)
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
