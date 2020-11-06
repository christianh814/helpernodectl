package cmd

import (
	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pulls images into your node",
	Long: `This will pull the images onto your local host. These images are used to
start all the services needed for the HelperNode. These images are:

quay.io/helpernode/pxe
quay.io/helpernode/http
quay.io/helpernode/loadbalancer
quay.io/helpernode/dns
quay.io/helpernode/dhcp`,
	Run: func(cmd *cobra.Command, args []string) {
		pullImages()
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pullCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pullCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// Loop through images and pull them
func pullImages() {
	for _, v := range images {
		PullImage(v, "latest")
	}
}
