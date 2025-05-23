package cmd

import (
	"fmt"

	"github.com/Hubert-Heijkers/tm1ctl/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// The config keys we allow to set using the config command
var allowedConfigKeys = map[string]bool{
	"output-format": true,
}

// The config keys we show when list all configurations
var listConfigKeys = map[string]bool{
	"host":          true,
	"instance":      true,
	"database":      true,
	"user":          true,
	"output-format": true,
}

var allowedOutputFormats = map[string]bool{
	"json":  true,
	"table": true,
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage global tm1ctl configuration",
}

func getConfigValue(key string) any {
	switch key {
	case "instance":
		host := getConfigValue("host")
		if host == nil || host == "" {
			return nil
		}
		instance, err := utils.GetInstanceName(host.(string), "")
		if err != nil || instance == "" {
			return nil
		}
		return instance

	case "database":
		instance := getConfigValue("instance")
		if instance == nil || instance == "" {
			return nil
		}

		// TODO: Retrieve active database from the instance once we maintain instances and have the ability to set such active database
		return nil
	}
	return viper.Get(key)
}

var configListCmd = &cobra.Command{
	Use:   "list [key]",
	Short: "List all configuration values",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 && args[0] != "" {
			key := args[0]

			if !listConfigKeys[key] {
				fmt.Printf("Error: '%s' is not a recognized configuration key\n", key)
				return
			}

			val := getConfigValue(key)
			if val != nil && val != "" {
				fmt.Printf("%s = %s\n", key, utils.Stringify(val))
			} else {
				fmt.Printf("No value set for key '%s'", key)
			}

		} else {
			for key := range listConfigKeys {
				val := getConfigValue(key)
				if val != nil && val != "" {
					fmt.Printf("%s = %s\n", key, utils.Stringify(val))
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
