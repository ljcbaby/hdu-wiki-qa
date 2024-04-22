package cmd

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hdu-wiki-qa",
	Short: "hdu-wiki-qa is a GPT chat for hduers with hdu-wiki knowledge",
}

var (
	config string
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().StringVarP(&config, "config", "c", "config.yaml", "config file path")
}
