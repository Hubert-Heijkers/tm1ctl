package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/Hubert-Heijkers/tm1ctl/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	userName      string
	userPassword  string
	userVariables string
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage user's credentials and session variables",
}

var userListCmd = &cobra.Command{
	Use:   "list [name]",
	Short: "List all specified users",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		users := viper.GetStringMap("users")

		if len(users) == 0 {
			fmt.Println("No users specified.")
			return
		}

		if len(args) == 1 && args[0] != "" {
			name := args[0]
			user := users[name]
			if user == nil {
				fmt.Printf("no details specified for user '%s'\n", name)
				return
			}
			err := utils.OutputMap(user.(map[string]any), "Name")
			cobra.CheckErr(err)
			return
		}

		// TODO: Highlight/mark the one that's active
		err := utils.OutputMap(users, "Name")
		cobra.CheckErr(err)
	},
}

var userSetCmd = &cobra.Command{
	Use:   "set <name>",
	Short: "Set credential or session variables for the specified user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		// Validate or initialize the users entry
		users := viper.GetStringMap("users")
		raw := users[name]
		var userMap map[string]any

		if raw == nil {
			userMap = make(map[string]any)
		} else {
			cast, ok := raw.(map[string]any)
			if !ok {
				fmt.Printf("Invalid user format for '%s'\n", name)
				return
			}
			userMap = cast
		}

		changed := false

		if userName != "" {
			userMap["name"] = userName
			changed = true
		} else {
			if cmd.Flags().Changed("name") {
				delete(userMap, "name")
				changed = true
			}
		}

		if userPassword != "" {
			userMap["password"] = userPassword
			changed = true
		} else {
			if cmd.Flags().Changed("password") {
				delete(userMap, "password")
				changed = true
			}
		}

		if userVariables != "" {
			// Variables set this way require the string to be a representation of a JSON object where every variable is represented by a property
			var varMap map[string]any
			err := json.Unmarshal([]byte(userVariables), &varMap)
			if err != nil {
				fmt.Printf("Value specified for variables is not a valid map")
				return
			}
			userMap["variables"] = varMap
			changed = true
		} else {
			if cmd.Flags().Changed("variables") {
				delete(userMap, "variables")
				changed = true
			}
		}

		if !changed {
			fmt.Println("No values provided to set. Use --name, --password or --variables.")
			return
		}

		users[name] = userMap
		viper.Set("users", users)
		cobra.CheckErr(utils.SaveConfiguration())
		fmt.Printf("Updated user '%s'\n", name)
	},
}

var userVariableCmd = &cobra.Command{
	Use:   "variable",
	Short: "Manage user's session variables",
}

var userVariableListCmd = &cobra.Command{
	Use:   "list [key]",
	Short: "List all the user's session variables specified",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// No user specified then use active user
		user, err := utils.GetUserName(userName)
		cobra.CheckErr(err)

		// Validate or initialize the users entry
		users := viper.GetStringMap("users")
		raw := users[user]

		if raw != nil {
			cast, ok := raw.(map[string]any)
			if !ok {
				fmt.Printf("Invalid user format for '%s'\n", user)
				return
			}
			raw = cast["variables"]
		}

		var varMap map[string]any

		if raw == nil {
			if len(args) == 1 && args[0] != "" {
				fmt.Printf("No value set for variable '%s'", args[0])
			}
		} else {
			cast, ok := raw.(map[string]any)
			if !ok {
				fmt.Printf("Invalid user variables format for '%s'\n", user)
				return
			}
			varMap = cast
		}

		if len(args) == 1 && args[0] != "" {
			key := args[0]
			val := varMap[key]
			if val != nil && val != "" {
				fmt.Printf("%s = %s\n", key, utils.Stringify(val))
			} else {
				fmt.Printf("No value set for variable '%s'", key)
			}

		} else {
			for key, val := range varMap {
				if val != nil && val != "" {
					fmt.Printf("%s = %s\n", key, utils.Stringify(val))
				}
			}
		}
	},
}

var userVariableSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a user's session variable to the specified value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		// No user specified then use active user
		user, err := utils.GetUserName(userName)
		cobra.CheckErr(err)

		// Validate or initialize the users entry
		users := viper.GetStringMap("users")
		raw := users[user]
		var userMap map[string]any

		if raw == nil {
			userMap = make(map[string]any)
		} else {
			cast, ok := raw.(map[string]any)
			if !ok {
				fmt.Printf("Invalid user format for '%s'\n", user)
				return
			}
			userMap = cast
		}

		raw = userMap["variables"]
		var varMap map[string]any

		if raw == nil {
			varMap = make(map[string]any)
		} else {
			cast, ok := raw.(map[string]any)
			if !ok {
				fmt.Printf("Invalid user variables format for '%s'\n", user)
				return
			}
			varMap = cast
		}

		key := args[0]

		// Note that a JSON value is expected here but if it doesn't parse as JSON we'll treat it as a string (the most common usage presumably)
		var value any
		err = json.Unmarshal([]byte(args[1]), &value)
		if err != nil {
			value = args[1]
		}

		varMap[key] = value
		userMap["variables"] = varMap
		users[user] = userMap
		viper.Set("users", users)
		cobra.CheckErr(utils.SaveConfiguration())
		fmt.Printf("Set variable %s to %s for user %s\n", key, utils.Stringify(value), user)
	},
}

var userUseCmd = &cobra.Command{
	Use:   "use [name]",
	Short: "Switch to using the specified user, or unset if no name given",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 && args[0] != "" {
			name := args[0]

			users := viper.GetStringMap("users")
			if _, exists := users[name]; !exists {
				fmt.Printf("User '%s' is not defined. Please configure a user before making it the active user.\n", name)
				return
			}
			viper.Set("user", name)
			cobra.CheckErr(utils.SaveConfiguration())
			fmt.Printf("Set active user to '%s'.\n", name)
		} else {
			viper.Set("user", "")
			cobra.CheckErr(utils.SaveConfiguration())
			fmt.Println("Reset active user.")
		}
	},
}

var userDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete user from the list of, configured, users",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		users := viper.GetStringMap("users")
		if _, ok := users[name]; !ok {
			fmt.Printf("User '%s' does not exist.\n", name)
			return
		}

		// Unset user if we're deleting the active one
		current := viper.GetString("user")
		if name == current {

			viper.Set("user", "")
		}

		// Delete the user from the list of users
		delete(users, name)
		viper.Set("users", users)
		cobra.CheckErr(utils.SaveConfiguration())
		fmt.Printf("Deleted user '%s'.\n", name)
	},
}

func init() {

	userCmd.AddCommand(userListCmd)

	userSetCmd.Flags().StringVar(&userName, "name", "", "Set the user name for this user")
	userSetCmd.Flags().StringVar(&userPassword, "password", "", "Set the password for this user")
	userSetCmd.Flags().StringVar(&userVariables, "variables", "", "Set the session variables for this user")
	userCmd.AddCommand(userSetCmd)

	userVariableListCmd.Flags().StringVar(&userName, "user", "", "The user to list the variables from, if not specified the active user will be used")
	userVariableCmd.AddCommand(userVariableListCmd)

	userVariableSetCmd.Flags().StringVar(&userName, "user", "", "The user to set the variable for, if not specified the active user will be used")
	userVariableCmd.AddCommand(userVariableSetCmd)

	userCmd.AddCommand(userVariableCmd)

	userCmd.AddCommand(userUseCmd)

	userCmd.AddCommand(userDeleteCmd)

	rootCmd.AddCommand(userCmd)
}
