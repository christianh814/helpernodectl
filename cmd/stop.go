package cmd

import (
	"fmt"
	"os"
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
		svc, _ := cmd.Flags().GetString("service")
		if len(svc) > 0 {
			// check if the string passed matches a service
			if _, exists := images[svc]; exists {
				//if it matches; create a "single service" map and pass that to the stop function
				singleservicemap := map[string]string{svc:images[svc]}
				stopImages(singleservicemap)
			} else {
				// If I didn't find it...tell them
				fmt.Println("Invalid service: " + svc)
				os.Exit(12)
			}
		} else {
			stopImages(images)
		}
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
	stopCmd.Flags().String("service", "", "stop a specifc service/container. Valid service names: dns, dhcp, http, loadbalancer, pxe")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stopCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stopCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func stopImages(imgs map[string]string) {
	for k, _ := range imgs {
		if !IsImageRunning("helpernode-" + k) {
			fmt.Println("SKIPPING: Container helpernode-" + k + " already stopped.")
		} else {
			StopImage(k)
		}
	}
}
