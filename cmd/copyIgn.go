package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"github.com/spf13/cobra"
)

// copyIgnCmd represents the copyIgn command
var copyIgnCmd = &cobra.Command{
	Use:   "copy-ign",
	Aliases:   []string{"copy-ignition", "cp-ign", "copyign", "cpign", "copyignition"},
	Short: "Copies ignition configs from given direcotry into the http container.",
	Long: `This command takes ignition configurations from the given directory,
and copies those files into the http contianer. For example:

	helpernodectl copy-ign --dir=/path/to/install/dir

This command must be run on the host that is to be the
helpernode. There is no support for copying the ignition
files to an external webserver.`,
	Run: func(cmd *cobra.Command, args []string) {
		// take what what passed into --dir
		dir, _ := cmd.Flags().GetString("dir")

		//check if it exists
		if _, err := os.Stat(dir); os.IsNotExist(err) {
		        fmt.Println("Please specify your install directory where the ignition files are.")
		        os.Exit(153)
		} else {
			copyIgnToWebServer(dir)
		}

	},
}

func init() {
	rootCmd.AddCommand(copyIgnCmd)
	copyIgnCmd.PersistentFlags().String("dir", "", "the directory where the ignition files are located")

	//make the --dir flag required
	copyIgnCmd.MarkPersistentFlagRequired("dir")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// copyIgnCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// copyIgnCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func copyIgnToWebServer(dir string) {

	//if the webserver isn't running, don't bother running
	if !IsImageRunning("helpernode-http") {
		fmt.Fprintf(os.Stderr, "ERROR: helpernode-http isn't running\n")
		os.Exit(253)
	}

	// get the list of the ignition files from the directory
	dirglob := dir + "/*.ign"
	ignfiles, err := filepath.Glob(dirglob)

	// check for error of command
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting ignition files %s: %s\n", ignfiles, err)
		os.Exit(253)
	}

	// Check to see if there's ANY igniton files there
	if len(ignfiles) == 0 {
		fmt.Fprintf(os.Stderr, "ERROR: No ignition files found in: " + dir + "\n")
		os.Exit(253)
	}

	// itterate over them
	clientpath := "/var/www/html/ignition/"
	for _, v := range ignfiles {
		fmt.Println("Copying over " + v + " to http container")
		// copy the igniton file to the http server - I should create a generic get/put function
		cmd, err := exec.Command("podman", "cp", v, "helpernode-http:" + clientpath).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running podman-cp command %s: %s\n", cmd, err)
			os.Exit(253)
		}
	}

	// Now we must fix permissions
	fixcmd, err := exec.Command("podman", "exec", "-it", "helpernode-http", "chmod", "+r", "-R", clientpath).Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running podman-exec command %s: %s\n", fixcmd, err)
		os.Exit(253)
	}
}
