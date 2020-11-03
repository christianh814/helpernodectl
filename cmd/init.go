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
	"os/exec"
	"github.com/robertsandoval/ocp4-helpernode/utils"
)

var containerRuntime string
var dhcp string = "quay.io/helpernode/dhcp"
var dns string = "quay.io/helpernode/dns"
var http string = "quay.io/helpernode/http"
var loadbalancer string = "quay.io/helpernode/loadbalancer"
var pxe string = "quay.io/helpernode/pxe"


// A map value containing some key-value pairs.
var images =  map[string]string{
	"dns": "quay.io/helpernode/dns",
	"dhcp": "quay.io/helpernode/dhcp",
	"http": "quay.io/helpernode/http",
	"loadbalancer": "quay.io/helpernode/loadbalancer",
	"pxe": "quay.io/helpernode/pxe",
	}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		verifyContainerRuntime()
		pullImages()

	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func verifyContainerRuntime() {

    out, err := exec.LookPath("podman")
    if err != nil {
		fmt.Println("Podman not found, checking for Docker")

		out, err = exec.LookPath("docker")
		if err != nil {
			fmt.Println("Docker not found...Please install Podman or Docker")
		}
	}
	fmt.Printf("found %s", out)
	containerRuntime = out
}


func pullImages(){
	for k, v := range images {
//		fmt.Printf("key[%s] value[%s]\n", k, v)
		fmt.Println("Pulling : " + k)
		utils.PullImage(v, "latest")
	}
}
