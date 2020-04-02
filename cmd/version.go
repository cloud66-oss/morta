package cmd

import (
	"fmt"

	"github.com/cloud66-oss/trackman/utils"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of shutdown-sequencer",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("shutdown-sequencer")
		fmt.Println("(c) 2020 Cloud66 Inc.")
		fmt.Println("shutdown-sequencer is a CLI to run a sequence of signals against a process")
		fmt.Printf("%s/%s\n", utils.Channel, utils.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
