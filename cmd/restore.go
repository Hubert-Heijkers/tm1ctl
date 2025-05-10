/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Hubert-Heijkers/tm1ctl/internal/utils"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore [backupset]",
	Short: "Performs a database restore using the specified backupset",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Restore initiated on database '%s' running on instance '%s' using backupset: %s\n", viper.GetString("database"), viper.GetString("service-instance"), args[0])

		// Check if the .backupsets folder exists
		_, err := utils.DatabaseAPIGet("Contents('Files')/Contents('.backupsets')")
		if err != nil {
			// TODO: Presuming failure implies the .backupsets folder doesn't exist, should check 404 Not Found response to be sure
			folderEntryPayload := map[string]any{"@odata.type": "#ibm.tm1.api.v1.Folder", "Name": ".backupsets"}
			_, err := utils.DatabaseAPIPost("Contents('Files')/Contents", folderEntryPayload)
			cobra.CheckErr(err)
		}

		// Retrieve the backupset path and check if the file exists
		backupsetPath := args[0]
		_, err = os.Stat(backupsetPath)
		if err != nil {
			cobra.CheckErr(err)
		}

		// Generate a unique, temporary, name to use for the backupset
		backupsetTempName := uuid.New().String() + "-" + filepath.Base(backupsetPath)

		// Create an entry for this backupset in the .backupsets folder
		documentEntryPayload := map[string]any{"@odata.type": "#ibm.tm1.api.v1.Document", "Name": backupsetTempName}
		_, err = utils.DatabaseAPIPost("Contents('Files')/Contents('.backupsets')/Contents", documentEntryPayload)
		cobra.CheckErr(err)

		// Now that we created this new, temporary, document in the .backupsets folder, lets make sure we dispose of it as well!
		defer func() {
			path := fmt.Sprintf("Contents('Files')/Contents('.backupsets')/Contents('%s')", backupsetTempName)
			err := utils.DatabaseAPIDelete(path)
			if err != nil {
				err = fmt.Errorf("temporary backupset '%s', stored in '.backupsets' under files, could not be delete due to: %w", backupsetTempName, err)
				fmt.Println("Warning:", err)
			}
		}()

		// Now let's upload the contents of the backupset to the newly created entry
		path := fmt.Sprintf("Contents('Files')/Contents('.backupsets')/Contents('%s')/Content", backupsetTempName)
		err = utils.DatabaseAPIPutFile(path, backupsetPath)
		// TODO: CheckErr os.exists on error which prevents the defered method cleaning up from executing!
		cobra.CheckErr(err)

		// Now that the backupset is available to the database we can perform the restore
		restorePayload := map[string]any{"URL": backupsetTempName}
		_, err = utils.DatabaseAPIPost("tm1s.Restore", restorePayload)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	restoreCmd.PersistentFlags().StringP("instance", "i", "", "The TM1 service instance in which the database resides")
	viper.BindPFlag("service-instance", restoreCmd.PersistentFlags().Lookup("instance"))
	restoreCmd.PersistentFlags().StringP("database", "d", "", "The TM1 database we want to create/restore")
	viper.BindPFlag("database", restoreCmd.PersistentFlags().Lookup("database"))
	restoreCmd.PersistentFlags().StringP("user", "u", "", "The user name needed to authenticate with the TM1 instance")
	viper.BindPFlag("user", restoreCmd.PersistentFlags().Lookup("user"))
	restoreCmd.PersistentFlags().StringP("client-secret", "p", "", "The password needed to authenticate with the TM1 instance")
	viper.BindPFlag("password", restoreCmd.PersistentFlags().Lookup("password"))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// restoreCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
