package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// sayhelloserviceCmd represents the sayhelloservice command
var sayhelloserviceCmd = &cobra.Command{
	Use:   "sayhelloservice",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sayhelloservice called")
	},
}

func init() {
	rootCmd.AddCommand(sayhelloserviceCmd)
	sayhelloserviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
