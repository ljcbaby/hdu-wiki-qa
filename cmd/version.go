package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "-dev"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of hdu-wiki-qa",
	Long:  `All software has versions. This is hdu-wiki-qa's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("hdu-wiki-qa v%s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
