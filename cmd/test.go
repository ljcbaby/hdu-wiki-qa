package cmd

import (
	"github.com/ljcbaby/hdu-wiki-qa/conf"
	"github.com/ljcbaby/hdu-wiki-qa/test"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test api",
	Long:  `Test api connection`,
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Test config",
	Long:  `Try to load and prase config file`,
	Run: func(cmd *cobra.Command, args []string) {
		conf.Init(config)
		test.Config(config)
	},
}

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Test chat model",
	Long:  `Test chat model with preset prompts`,
	Run: func(cmd *cobra.Command, args []string) {
		conf.Init(config)
		test.Chat()
	},
}

var embeddingCmd = &cobra.Command{
	Use:   "embedding",
	Short: "Test embedding model",
	Long:  `Test embedding model with preset strings`,
	Run: func(cmd *cobra.Command, args []string) {
		conf.Init(config)
		test.Embedding()
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.AddCommand(configCmd)
	testCmd.AddCommand(chatCmd)
	testCmd.AddCommand(embeddingCmd)
}
