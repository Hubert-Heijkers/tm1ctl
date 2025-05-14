package cmd

import (
	"fmt"

	"github.com/Hubert-Heijkers/tm1ctl/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var allowedConfigKeys = map[string]bool{
	"database":      true,
	"user":          true,
	"password":      true,
	"output-format": true,
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
	Use:   "list [key]",
	Short: "List all configuration values",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 && args[0] != "" {
			key := args[0]

			if !allowedConfigKeys[key] {
				fmt.Printf("Error: '%s' is not a recognized configuration key\n", key)
				return
			}

			val := viper.Get(key)
			if val != nil && val != "" {
				fmt.Printf("%s = %v\n", key, val)
			} else {
				fmt.Printf("No value set for key '%s'", key)
			}

		} else {
			for key := range allowedConfigKeys {
				val := viper.Get(key)
				if val != nil && val != "" {
					fmt.Printf("%s = %v\n", key, val)
				}
			}
		}
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
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
		err := utils.SaveConfiguration()
		cobra.CheckErr(err)
		fmt.Printf("%s set to %s\n", key, value)
	},
}

func init() {
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configSetCmd)
	rootCmd.AddCommand(configCmd)
}
