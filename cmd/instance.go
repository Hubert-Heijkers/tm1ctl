package cmd

import (
	"fmt"

	"github.com/Hubert-Heijkers/tm1ctl/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Variables used for flags in any cmd
var (
	host     string
	instance string
	database string
	user     string
	password string
)

// instanceCmd represents the instance command
var instanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "Manage the instances of a TM1 v12 service",
}

var instanceListCmd = &cobra.Command{
	Use:   "list [name]",
	Short: "Get the list of TM1 service instances",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var path string
		if len(args) == 1 && args[0] != "" {
			path = fmt.Sprintf("Instances('%s')", args[0])
		} else {
			path = "Instances"
		}
		data, err := utils.ManageAPIGet(host, path)
		cobra.CheckErr(err)
		// TODO: Highlight/mark the one that is active adding it if no configuration for that instance exists!
		err = utils.OutputCollection(data)
		cobra.CheckErr(err)
	},
}

var instanceCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Creates a new TM1 service instance with the name specified",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		payload := map[string]any{"Name": args[0]}
		data, err := utils.ManageAPIPost(host, "Instances", payload)
		cobra.CheckErr(err)
		err = utils.OutputEntity(data)
		cobra.CheckErr(err)
	},
}

var instanceDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Deletes the TM1 service instance and all its associated databases and artifacts",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		instanceName := args[0]
		path := fmt.Sprintf("Instances('%s')", instanceName)
		err := utils.ManageAPIDelete(host, path)
		cobra.CheckErr(err)
		fmt.Printf("Instance '%s' has been deleted!\n", instanceName)
	},
}

var instanceUseCmd = &cobra.Command{
	Use:   "use [name]",
	Short: "Switch to using the specified instance, or unset if no name given",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get the host name
		host, err := utils.GetHostName(host)
		cobra.CheckErr(err)

		// Lookup the host in list of configured hosts
		hosts := viper.GetStringMap("hosts")
		raw := hosts[host]
		if raw == nil {
			fmt.Printf("no configuration specified for host '%s'", host)
			return
		}
		hostMap, ok := raw.(map[string]any)
		if !ok {
			fmt.Printf("invalid configuration for host '%s', format invalid", host)
			return
		}

		// Update the host's configuration accordingly
		if len(args) == 1 && args[0] != "" {
			name := args[0]
			hostMap["instance"] = name
			hosts[host] = hostMap
			viper.Set("hosts", hosts)
			cobra.CheckErr(utils.SaveConfiguration())
			fmt.Printf("Set active instance on host '%s' to '%s'.\n", host, name)
		} else {
			hostMap["instance"] = ""
			hosts[host] = hostMap
			viper.Set("hosts", hosts)
			cobra.CheckErr(utils.SaveConfiguration())
			fmt.Printf("Reset active instance on host '%s'.\n", host)
		}
	},
}

func init() {

	instanceListCmd.Flags().StringVar(&host, "host", "", "The host to list the instance from, if not specified the active host will be used")
	instanceCmd.AddCommand(instanceListCmd)

	instanceCreateCmd.Flags().StringVar(&host, "host", "", "The host to list the instance from, if not specified the active host will be used")
	instanceCmd.AddCommand(instanceCreateCmd)

	instanceDeleteCmd.Flags().StringVar(&host, "host", "", "The host to list the instance from, if not specified the active host will be used")
	instanceCmd.AddCommand(instanceDeleteCmd)

	instanceUseCmd.Flags().StringVar(&host, "host", "", "The host to list the instance from, if not specified the active host will be used")
	instanceCmd.AddCommand(instanceUseCmd)

	rootCmd.AddCommand(instanceCmd)
}
