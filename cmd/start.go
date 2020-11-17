package cmd

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"github.com/spf13/cobra"
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
		/* this calls preflight, but I need to make preflight smarter
		   right now, errors do not exit the cli. So, it'll try to start
		   even if it find conflicts
		preflightCmd.Run(cmd, []string{})
		*/
		startImages()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// For now, start all images
func startImages() {

	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		fmt.Println("Please specify a config file")
		os.Exit(153)
	} else {

		// Open file on disk
		f, _ := os.Open(cfgFile)

		// Read file into a byte
		reader := bufio.NewReader(f)
		content, _ := ioutil.ReadAll(reader)

		//Encode to base64
		encoded := base64.StdEncoding.EncodeToString(content)

		// run the containers using the encoding
		for k, v := range images {
			if IsImageRunning("helpernode-" + k) {
				fmt.Println("SKIPPING: Container helpernode-" + k + " already running.")
			} else {
				StartImage(v, "latest", encoded, k)
			}
		}
	}
}
