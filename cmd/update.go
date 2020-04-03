package cmd

import (
	"fmt"

	"github.com/cloud66-oss/morta/utils"
	"github.com/cloud66-oss/updater"
	"github.com/spf13/cobra"
)

var (
	channel string

	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update morta to the latest version",
		RunE:  updateExec,
	}
)

func init() {
	updateCmd.Flags().StringVarP(&channel, "channel", "c", utils.Channel, "version channel")

	rootCmd.AddCommand(updateCmd)
}

func updateExec(cmd *cobra.Command, args []string) error {
	worker, err := updater.NewUpdater(utils.Version, &updater.Options{
		RemoteURL: "https://s3.amazonaws.com/downloads.cloud66.com/morta/",
		Channel:   channel,
	})
	if err != nil {
		return err
	}

	err = worker.Run(channel != utils.Channel)
	if err != nil {
		return err
	}

	fmt.Println("Update complete")
	return nil
}
