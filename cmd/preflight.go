package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net"
)

// preflightCmd represents the preflight command
var preflightCmd = &cobra.Command{
	Use:   "preflight",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fixall, _ := cmd.Flags().GetBool("fix-all")
		if fixall {
			fmt.Printf("Checking for conflicts\nBEST EFFORT IN FIXING ERRORS\n============================\n")
			portCheck()
		} else {
			fmt.Printf("Checking for conflicts\n======================\n")
			portCheck()
		}
	},
}

func init() {
	rootCmd.AddCommand(preflightCmd)
	preflightCmd.Flags().BoolP("fix-all", "x", false, "Does the needful and fixes errors it finds")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// preflightCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// preflightCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func portCheck() {
	// check each port
	porterrorcount := 0
	for _, p := range ports {
		//check if you can listen on this port on TCP
		t, err := net.Listen("tcp", ":" + p)

		// If this returns an error, then something else is listening on this port
		if err != nil {
			fmt.Println("WARNING: Port tcp:" + p + " is in use")
			porterrorcount += 1
		} else {
			t.Close()
		}

		//now check if you can listen on this port on UDP
		u, err := net.ListenPacket("udp", ":" + p)

		// If this returns an error, then something else is listening on this port
		if err != nil {
			fmt.Println("WARNING: Port udp:" + p + " is in use")
			porterrorcount += 1
		} else {
			u.Close()
		}

	}

	// Display that no errors were found
	if porterrorcount == 0 {
		fmt.Println("No port confilcts were found")
	}
}
