package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
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
			fmt.Println("preflight with fix all")
		} else {
			fmt.Println("preflight called")
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

