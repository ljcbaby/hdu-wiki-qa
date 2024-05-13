package cmd

import (
	"github.com/ljcbaby/hdu-wiki-qa/conf"
	"github.com/ljcbaby/hdu-wiki-qa/database"
	"github.com/ljcbaby/hdu-wiki-qa/generate"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate embeddings data",
	Long:  `Generate embeddings data from wiki data`,
	Run: func(cmd *cobra.Command, args []string) {
		conf.Init(config)
		database.Connect()
		generate.Init()
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
