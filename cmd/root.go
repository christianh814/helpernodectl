package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// set the global config file as a string
var cfgFile string

// Env var to the image prefix
var HELPERNODE_IMAGE_PREFIX string

//cli clients as a map
var clients = map[string]string {
	"oc": "openshift-client-linux.tar.gz",
	"openshift-install": "openshift-install-linux.tar.gz",
	"helm": "helm.tar.gz",
}

// Define images and their registry location
var images = map[string]string {
	"dns": "quay.io/helpernode/dns",
	"dhcp": "quay.io/helpernode/dhcp",
	"http": "quay.io/helpernode/http",
	"loadbalancer": "quay.io/helpernode/loadbalancer",
	"pxe": "quay.io/helpernode/pxe",
}

// Define ports needed for preflight check of listening ports
var ports = [10]string{"67", "546", "53", "80", "443", "69", "6443", "22623", "8080", "9000"}

// Define firewalld rules needed to be in place
var fwrule = [13]string {
	"6443/tcp",
	"22623/tcp",
	"8080/tcp",
	"9000/tcp",
	"9090/tcp",
	"67/udp",
	"546/udp",
	"53/tcp",
	"53/udp",
	"80/tcp",
	"443/tcp",
	"22/tcp",
	"69/udp",
}

// Define systemd services we will check
var systemdsvc = map[string]string {
	"resolved": "systemd-resolved.service",
	"dnsmasq" : "dnsmasq.service",
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command {
	Use:   "helpernodectl",
	Short: "Utility for the HelperNode",
	Long: `This cli utility is used to stop/start the HelperNode
on the host it's ran from. You need to provide a helpernode.yaml file
with information about your helper config. A simple example to start
your HelperNode is:

helpernodectl start --config=helpernode.yaml`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	verifyContainerRuntime()
	verifyFirewallCommand()
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.helpernode.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("all", "a", false, "do it for all containers")

	if len(os.Getenv("HELPERNODE_IMAGE_PREFIX")) > 0 {
		// Define images and their registry location based on the env var
		imageprefix := os.Getenv("HELPERNODE_IMAGE_PREFIX")
		images = map[string]string {
			"dns": imageprefix + "/helpernode/dns",
			"dhcp": imageprefix + "/helpernode/dhcp",
			"http": imageprefix + "/helpernode/http",
			"loadbalancer": imageprefix + "/helpernode/loadbalancer",
			"pxe": imageprefix + "/helpernode/pxe",
		}
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".helpernode" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".helpernode")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		//fmt.Println("Using config file:", viper.ConfigFileUsed())
		//fmt.Fprintf(os.Stderr, "Error reading config: %s\n", err)
		//os.Exit(1)
		/*
			Do nothing for now. Checking for the config file will rely on the 
			individual "subcommands" and their children. See cmd/start.go for
			an example.
		*/
	}
}

// Check runtime. Right now it's just podman
func verifyContainerRuntime() {

	_, err := exec.LookPath("podman")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error looking for Podman: %s\n", err)
		os.Exit(1)
	}
}

// Check runtime. Right now it's just podman
func verifyFirewallCommand() {

	_, err := exec.LookPath("firewall-cmd")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error looking for firewall-cmd: %s\n", err)
		os.Exit(1)
	}
}
