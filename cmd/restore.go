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
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore <backup-set>",
	Short: "Performs a database restore using the specified backup-set",
	Run: func(cmd *cobra.Command, args []string) {

		// TODO: Add Use command to Database and use active database for the instance
		if database == "" {
			fmt.Println("no database specified")
			return
		}

		// No instance specified then use active instance
		instance, err := utils.GetInstanceName(host, instance)
		cobra.CheckErr(err)

		fmt.Printf("Restore initiated on database '%s' running on instance '%s' using backupset: %s\n", database, instance, args[0])

		// Check if the .backupsets folder exists
		_, err = utils.DatabaseAPIGet(host, instance, database, user, password, "Contents('Files')/Contents('.backupsets')")
		if err != nil {
			// TODO: Presuming failure implies the .backupsets folder doesn't exist, should check 404 Not Found response to be sure
			folderEntryPayload := map[string]any{"@odata.type": "#ibm.tm1.api.v1.Folder", "Name": ".backupsets"}
			_, err := utils.DatabaseAPIPost(host, instance, database, user, password, "Contents('Files')/Contents", folderEntryPayload)
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
		_, err = utils.DatabaseAPIPost(host, instance, database, user, password, "Contents('Files')/Contents('.backupsets')/Contents", documentEntryPayload)
		cobra.CheckErr(err)

		// Now that we created this new, temporary, document in the .backupsets folder, lets make sure we dispose of it as well!
		defer func() {
			path := fmt.Sprintf("Contents('Files')/Contents('.backupsets')/Contents('%s')", backupsetTempName)
			err := utils.DatabaseAPIDelete(host, instance, database, user, password, path)
			if err != nil {
				err = fmt.Errorf("temporary backupset '%s', stored in '.backupsets' under files, could not be delete due to: %w", backupsetTempName, err)
				fmt.Println("Warning:", err)
			}
		}()

		// Now let's upload the contents of the backupset to the newly created entry
		path := fmt.Sprintf("Contents('Files')/Contents('.backupsets')/Contents('%s')/Content", backupsetTempName)
		err = utils.DatabaseAPIPutFile(host, instance, database, user, password, path, backupsetPath)
		// TODO: CheckErr os.exists on error which prevents the defered method cleaning up from executing!
		cobra.CheckErr(err)

		// Now that the backupset is available to the database we can perform the restore
		restorePayload := map[string]any{"URL": backupsetTempName}
		_, err = utils.DatabaseAPIPost(host, instance, database, user, password, "tm1s.Restore", restorePayload)
		cobra.CheckErr(err)
	},
}

func init() {

	restoreCmd.Flags().StringVar(&host, "host", "", "The host on which the instance is running, if not specified the active host will be used")
	restoreCmd.Flags().StringVar(&instance, "instance", "", "The instance to be used, if not specified the active instance will be used")
	restoreCmd.Flags().StringVar(&database, "database", "", "The database you want to restore")
	restoreCmd.Flags().StringVar(&user, "user", "", "The user name needed to authenticate with the TM1 instance")
	restoreCmd.Flags().StringVar(&password, "password", "", "The password needed to authenticate with the TM1 instance")
	rootCmd.AddCommand(restoreCmd)
}
