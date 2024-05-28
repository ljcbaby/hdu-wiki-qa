package cmd

import (
	"github.com/ljcbaby/hdu-wiki-qa/conf"
	"github.com/ljcbaby/hdu-wiki-qa/database"
	"github.com/ljcbaby/hdu-wiki-qa/service"
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Run service",
	Long:  `Run chat api service`,
	Run: func(cmd *cobra.Command, args []string) {
		conf.Init(config)
		database.Connect()
		service.Serve()
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)
}
