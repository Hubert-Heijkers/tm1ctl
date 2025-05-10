/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/Hubert-Heijkers/tm1ctl/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// databaseCmd represents the database command
var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Manage the databases of your TM1 v12 service instance",
}

var databaseListCmd = &cobra.Command{
	Use:   "list",
	Short: "Get the list of TM1 database of the TM1 service instance in context",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := utils.InstanceAPIGet("Databases")
		cobra.CheckErr(err)
		err = utils.OutputCollection(data)
		cobra.CheckErr(err)
	},
}

var databaseGetCmd = &cobra.Command{
	Use:   "get [name]",
	Short: "Shows the details of a TM1 database",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := fmt.Sprintf("Databases('%s')", args[0])
		data, err := utils.InstanceAPIGet(path)
		cobra.CheckErr(err)
		err = utils.OutputEntity(data)
		cobra.CheckErr(err)
	},
}

var databaseCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Creates a new TM1 database with the specified name on the TM1 service instance in context",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		payload := map[string]any{"Name": args[0]}
		data, err := utils.InstanceAPIPost("Databases", payload)
		cobra.CheckErr(err)
		err = utils.OutputEntity(data)
		cobra.CheckErr(err)
	},
}

var databaseDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Deletes the TM1 database specified with all its artifacts",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		databaseName := args[0]
		path := fmt.Sprintf("Databases('%s')", databaseName)
		err := utils.InstanceAPIDelete(path)
		cobra.CheckErr(err)
		fmt.Printf("Database '%s' has been deleted!\n", databaseName)
	},
}

func init() {
	databaseCmd.AddCommand(databaseListCmd)
	databaseCmd.AddCommand(databaseGetCmd)
	databaseCmd.AddCommand(databaseCreateCmd)
	databaseCmd.AddCommand(databaseDeleteCmd)
	rootCmd.AddCommand(databaseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	databaseCmd.PersistentFlags().String("client-id", "", "The root client-id needed to authenticate with the TM1 service")
	viper.BindPFlag("root-client-id", databaseCmd.PersistentFlags().Lookup("client-id"))
	databaseCmd.PersistentFlags().String("client-secret", "", "The root client-secret needed to authenticate with the TM1 service")
	viper.BindPFlag("root-client-secret", databaseCmd.PersistentFlags().Lookup("client-secret"))
	databaseCmd.PersistentFlags().StringP("instance", "i", "", "The TM1 service instance to be used")
	viper.BindPFlag("service-instance", databaseCmd.PersistentFlags().Lookup("instance"))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//databaseCmd.Flags().String("foo", "", "foo flag")
}
