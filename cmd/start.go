/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/robertsandoval/ocp4-helpernode/utils"
	"os"
	"bufio"
	"encoding/base64"
	"io/ioutil"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		runContianers()
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

func runContianers() {
	// Check to see if file exists
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		fmt.Println("File " + cfgFile + " does not exist")
	} else {
		// Open file on disk
		f, _ := os.Open(cfgFile)
		// Read file into a byte slice
		reader := bufio.NewReader(f)
		content, _ := ioutil.ReadAll(reader)
		//Encode to base64
		encoded := base64.StdEncoding.EncodeToString(content)
		// run the containers using the encoding
		for k, v := range images {
			utils.StartImage(v, "latest", encoded, k)
		}
	}
}
