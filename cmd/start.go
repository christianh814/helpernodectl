package cmd

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts HelperNode containers based on the provided manifest.",
	Long: `This will start the containers needed for the HelperNode to run.
It will run the services depending on what manifest is passed.

Examples:

	helpernodectl start --config=helpernode.yaml
	
	cp helpernode.yaml ~/.helpernode.yaml
	helpernodectl start

This manifest should have all the information the services need to start
up successfully.`,
	Run: func(cmd *cobra.Command, args []string) {
		// get any options passed
		skippreflight, _ := cmd.Flags().GetBool("skip-preflight")
		svc, _ := cmd.Flags().GetString("service")
		///
		if len(svc) > 0 {
                        // check if the string passed matches a service
                        if _, exists := images[svc]; exists {
                                //if it matches; create a "single service" map and pass that to the stop function
                                singleservicemap := map[string]string{svc:images[svc]}
                                startImages(singleservicemap)
                        } else {
                                // If I didn't find it...tell them
                                fmt.Println("Invalid service: " + svc)
                                os.Exit(12)
                        }
		} else {
			if skippreflight {
				fmt.Printf("Skipping Preflightchecks\n======================\n")
				startImages(images)
			} else {
				preflightCmd.Run(cmd, []string{})
				fmt.Printf("Starting Containers\n======================\n")
				startImages(images)
			}
		}
		///
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().BoolP("skip-preflight", "", false, "Skips preflight checks and tries to start the containers")
	startCmd.Flags().String("service", "", "start a service/container (preflight NOT performed). Valid names: dns, dhcp, http, loadbalancer, pxe")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// Take the map given and start images based on that map: Default (above) `images` gets passed in
func startImages(imgs map[string]string) {

	if _, err := os.Stat(viper.ConfigFileUsed()); os.IsNotExist(err) {
		fmt.Println("Please specify a config file")
		os.Exit(153)
	} else {

		// Open file on disk
		f, _ := os.Open(viper.ConfigFileUsed())

		// Read file into a byte
		reader := bufio.NewReader(f)
		content, _ := ioutil.ReadAll(reader)

		//Encode to base64
		encoded := base64.StdEncoding.EncodeToString(content)

		// run the containers using the encoding
		for k, v := range imgs {
			if IsImageRunning("helpernode-" + k) {
				fmt.Println("SKIPPING: Container helpernode-" + k + " already running.")
			} else {
				if !DoISkip(k) {
					StartImage(v, "latest", encoded, k)
				}
			}
		}
	}
}
