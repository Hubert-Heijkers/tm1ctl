/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/Hubert-Heijkers/tm1ctl/internal/utils"
	"github.com/spf13/cobra"
)

// databaseCmd represents the database command
var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Manage the databases of your TM1 v12 service instance",
}

var databaseListCmd = &cobra.Command{
	Use:   "list [name]",
	Short: "Get the list of TM1 databases",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var path string
		if len(args) == 1 && args[0] != "" {
			path = fmt.Sprintf("Databases('%s')", args[0])
		} else {
			path = "Databases"
		}
		data, err := utils.InstanceAPIGet(host, instance, user, password, path)
		cobra.CheckErr(err)
		// TODO: Highlight/mark the one that is active!
		err = utils.OutputCollection(data)
		cobra.CheckErr(err)
	},
}

var databaseCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Creates a new TM1 database with the specified name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		payload := map[string]any{"Name": args[0]}
		data, err := utils.InstanceAPIPost(host, instance, user, password, "Databases", payload)
		cobra.CheckErr(err)
		err = utils.OutputEntity(data)
		cobra.CheckErr(err)
	},
}

var databaseDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Deletes the TM1 database specified with all its artifacts",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		databaseName := args[0]
		path := fmt.Sprintf("Databases('%s')", databaseName)
		err := utils.InstanceAPIDelete(host, instance, user, password, path)
		cobra.CheckErr(err)
		fmt.Printf("Database '%s' has been deleted!\n", databaseName)
	},
}

func init() {

	databaseListCmd.Flags().StringVar(&host, "host", "", "The host on which the instance is running, if not specified the active host will be used")
	databaseListCmd.Flags().StringVar(&instance, "instance", "", "The instance to be used, if not specified the active instance will be used")
	databaseListCmd.Flags().StringVar(&user, "user", "", "The user name needed to authenticate with the TM1 instance")
	databaseListCmd.Flags().StringVar(&password, "password", "", "The password needed to authenticate with the TM1 instance")
	databaseCmd.AddCommand(databaseListCmd)

	databaseCreateCmd.Flags().StringVar(&host, "host", "", "The host on which the instance is running, if not specified the active host will be used")
	databaseCreateCmd.Flags().StringVar(&instance, "instance", "", "The instance to be used, if not specified the active instance will be used")
	databaseCreateCmd.Flags().StringVar(&user, "user", "", "The user name needed to authenticate with the TM1 instance")
	databaseCreateCmd.Flags().StringVar(&password, "password", "", "The password needed to authenticate with the TM1 instance")
	databaseCmd.AddCommand(databaseCreateCmd)

	databaseDeleteCmd.Flags().StringVar(&host, "host", "", "The host on which the instance is running, if not specified the active host will be used")
	databaseDeleteCmd.Flags().StringVar(&instance, "instance", "", "The instance to be used, if not specified the active instance will be used")
	databaseDeleteCmd.Flags().StringVar(&user, "user", "", "The user name needed to authenticate with the TM1 instance")
	databaseDeleteCmd.Flags().StringVar(&password, "password", "", "The password needed to authenticate with the TM1 instance")
	databaseCmd.AddCommand(databaseDeleteCmd)

	rootCmd.AddCommand(databaseCmd)
}
