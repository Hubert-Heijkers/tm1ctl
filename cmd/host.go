package cmd

import (
	"fmt"

	"github.com/Hubert-Heijkers/tm1ctl/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	serviceRootURL   string
	rootClientId     string
	rootClientSecret string
)

// hostCmd represents the host command
var hostCmd = &cobra.Command{
	Use:   "host",
	Short: "Manage host configuration",
}

var hostListCmd = &cobra.Command{
	Use:   "list [name]",
	Short: "List all configured hosts",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hosts := viper.GetStringMap("hosts")

		if len(hosts) == 0 {
			fmt.Println("No hosts configured.")
			return
		}

		if len(args) == 1 && args[0] != "" {
			name := args[0]
			host := hosts[name]
			if host == nil {
				fmt.Printf("no configuration specified for host '%s'\n", name)
				return
			}
			err := utils.OutputCollection(host)
			cobra.CheckErr(err)
			return
		}

		// TODO: Highlight/mark the one that's active
		err := utils.OutputMap(hosts, "Name")
		cobra.CheckErr(err)
	},
}

var hostSetCmd = &cobra.Command{
	Use:   "set <name>",
	Short: "Set one or more configuration values for the specified host",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		// Validate or initialize the hosts entry
		hosts := viper.GetStringMap("hosts")
		raw := hosts[name]
		var hostMap map[string]any

		if raw == nil {
			hostMap = make(map[string]any)
		} else {
			cast, ok := raw.(map[string]any)
			if !ok {
				fmt.Printf("Invalid host format for '%s'\n", name)
				return
			}
			hostMap = cast
		}

		changed := false

		if serviceRootURL != "" {
			hostMap["service_root_url"] = serviceRootURL
			changed = true
		}

		if rootClientId != "" {
			hostMap["root_client_id"] = rootClientId
			changed = true
		}

		if rootClientSecret != "" {
			hostMap["root_client_secret"] = rootClientSecret
			changed = true
		}

		if !changed {
			fmt.Println("No values provided to set. Use --service_root_url, --root_client_id or --root_client_secret.")
			return
		}

		hosts[name] = hostMap
		viper.Set("hosts", hosts)
		cobra.CheckErr(utils.SaveConfiguration())
		fmt.Printf("Updated host '%s'\n", name)
	},
}

var hostUseCmd = &cobra.Command{
	Use:   "use [name]",
	Short: "Switch to using the specified host, or unset if no name given",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 && args[0] != "" {
			name := args[0]

			hosts := viper.GetStringMap("hosts")
			if _, exists := hosts[name]; !exists {
				fmt.Printf("Host '%s' is not defined. Please configure a host before making it the active host.\n", name)
				return
			}
			viper.Set("host", name)
			cobra.CheckErr(utils.SaveConfiguration())
			fmt.Printf("Set active host to '%s'.\n", name)
		} else {
			viper.Set("host", "")
			cobra.CheckErr(utils.SaveConfiguration())
			fmt.Println("Reset active host.")
		}
	},
}

var hostDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete host from the list of, configured, hosts",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		hosts := viper.GetStringMap("hosts")
		if _, ok := hosts[name]; !ok {
			fmt.Printf("Host '%s' does not exist.\n", name)
			return
		}

		// Unset host if we're deleting the active one
		current := viper.GetString("host")
		if name == current {

			viper.Set("host", "")
		}

		// Delete the host from the list of hosts
		delete(hosts, name)
		viper.Set("Hosts", hosts)
		cobra.CheckErr(utils.SaveConfiguration())
		fmt.Printf("Deleted host '%s'.\n", name)
	},
}

func init() {

	hostCmd.AddCommand(hostListCmd)

	hostSetCmd.Flags().StringVar(&serviceRootURL, "service_root_url", "", "Set the service root URL for this host")
	hostSetCmd.Flags().StringVar(&rootClientId, "root_client_id", "", "Set the root user's client ID for this host")
	hostSetCmd.Flags().StringVar(&rootClientSecret, "root_client_secret", "", "Set the root user's client secret for this host")
	hostCmd.AddCommand(hostSetCmd)

	hostCmd.AddCommand(hostUseCmd)

	hostCmd.AddCommand(hostDeleteCmd)

	rootCmd.AddCommand(hostCmd)
}
