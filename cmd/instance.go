package cmd

import (
	"fmt"

	"github.com/Hubert-Heijkers/tm1ctl/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// instanceCmd represents the instance command
var instanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "Manage the instances of a TM1 v12 service",
}

var instanceListCmd = &cobra.Command{
	Use:   "list",
	Short: "Get the list of a TM1 service instances",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := utils.ManageAPIGet("Instances")
		cobra.CheckErr(err)
		err = utils.OutputCollection(data)
		cobra.CheckErr(err)
	},
}

var instanceGetCmd = &cobra.Command{
	Use:   "get [name]",
	Short: "Shows the details of a TM1 service instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := fmt.Sprintf("Instances('%s')", args[0])
		data, err := utils.ManageAPIGet(path)
		cobra.CheckErr(err)
		err = utils.OutputEntity(data)
		cobra.CheckErr(err)
	},
}

var instanceCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Creates a new TM1 service instance with the specified name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		payload := map[string]any{"Name": args[0]}
		data, err := utils.ManageAPIPost("Instances", payload)
		cobra.CheckErr(err)
		err = utils.OutputEntity(data)
		cobra.CheckErr(err)
	},
}

var instanceDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Deletes the TM1 service instance specified with all its associated databases and artifacts",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		instanceName := args[0]
		path := fmt.Sprintf("Instances('%s')", instanceName)
		err := utils.ManageAPIDelete(path)
		cobra.CheckErr(err)
		fmt.Printf("Instance '%s' has been deleted!\n", instanceName)
	},
}

func init() {
	instanceCmd.AddCommand(instanceListCmd)
	instanceCmd.AddCommand(instanceGetCmd)
	instanceCmd.AddCommand(instanceCreateCmd)
	instanceCmd.AddCommand(instanceDeleteCmd)
	rootCmd.AddCommand(instanceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	instanceCmd.PersistentFlags().String("client-id", "", "The root client-id needed to authenticate with the TM1 service")
	viper.BindPFlag("root-client-id", instanceCmd.PersistentFlags().Lookup("client-id"))
	instanceCmd.PersistentFlags().String("client-secret", "", "The root client-secret needed to authenticate with the TM1 service")
	viper.BindPFlag("root-client-secret", instanceCmd.PersistentFlags().Lookup("client-secret"))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// instanceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
