package cmd

import (
	"fmt"
	"net"
	"github.com/spf13/cobra"
)

// preflightCmd represents the preflight command
var preflightCmd = &cobra.Command{
	Use:   "preflight",
	Short: "Checks for any conflicts on the host.",
	Long: `This checks for conflicts on the host and can optionally fix
errors it finds. For example:
	
	helpernodectl preflight

	helpernodectl preflight --fix-all


This checks for port conflicts, systemd conflicts, and also checks any 
firewall rules. It will optionally fix systemd and firewall rules by
passing the --fix-all option.`,
	Run: func(cmd *cobra.Command, args []string) {
		fixall, _ := cmd.Flags().GetBool("fix-all")
		if fixall {
			fmt.Printf("Checking for conflicts\nBEST EFFORT IN FIXING ERRORS\n============================\n")
			systemdCheck(true)
			portCheck()
		} else {
			fmt.Printf("Checking for conflicts\n======================\n")
			systemdCheck(false)
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

func systemdCheck(fix bool) {
	svcerrorcount := 0
	for _, s := range systemdsvc {
		if IsServiceRunning(s) {
			fmt.Println("WARNING: Service " + s + " is running")
			svcerrorcount += 1
			if fix {
				fmt.Println("STOPPING/DISABLING SERVICE: " + s)
				StopService(s)
				DisableService(s)
			}
		}
	}
	// Display that no errors were found
	if svcerrorcount == 0 {
		fmt.Println("No service confilcts were found")
	}
}
