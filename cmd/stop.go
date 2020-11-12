package cmd

import (
	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops containers that are running for the HelperNode",
	Long: `This stops containers used for running the services needed by the HelperNode.
This stops the container immedietly. Example:

	helpernodectl stop

It's important to note that an "rm" is passed to the runtime.
However, the image is preserved on the host.`,
	Run: func(cmd *cobra.Command, args []string) {
		stopImages()
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stopCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stopCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func stopImages() {
	for k, _ := range images {
		if !IsImageRunning("helpernode-" + k) {
			fmt.Println("SKIPPING: Container helpernode-" + k + " already stopped.")
		} else {
			StopImage(k)
		}
	}
}
