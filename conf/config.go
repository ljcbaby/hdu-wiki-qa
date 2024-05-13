package conf

import (
	_ "embed"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//go:embed config.example.yaml
var configSample []byte

var Pgsql struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

var Api struct {
	BaseURL        string
	ApiKey         string
	ChatModel      string
	EmbeddingModel string
}

var Advanced struct {
	SimilarityThreshold float64
	SearchLength        int
	SystemPrompt        string
}

var Wiki struct {
	Dir     string
	Format  []string
	Exclude []string
}

func Init(file string) {
	if _, err := os.Stat(file); err != nil {
		if !os.IsNotExist(err) {
			logrus.WithError(err).Fatalf("get stat of %s failed", file)
		}
		logrus.Infof("config not existed, creating at %s", file)
		created, err := os.Create(file)
		if err != nil {
			logrus.WithError(err).Fatalf("create config at %s failed", file)
		}
		if _, err := created.Write(configSample); err != nil {
			logrus.WithError(err).Fatalf("write config at %s failed", file)
		}
	}

	viper.SetConfigFile(file)
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		logrus.WithError(err).Fatalf("read config from %s failed", file)
	}

	setConfig()
	if err := checkNecessary(); err != nil {
		logrus.WithError(err).Fatalf("set config failed")
		os.Exit(1)
	}
}

func setConfig() error {
	Pgsql.Host = viper.GetString("postgresql.host")
	Pgsql.Port = viper.GetInt("postgresql.port")
	Pgsql.User = viper.GetString("postgresql.username")
	Pgsql.Password = viper.GetString("postgresql.password")
	Pgsql.Database = viper.GetString("postgresql.database")

	Api.BaseURL = viper.GetString("api.base_url")
	Api.ApiKey = viper.GetString("api.api_key")
	Api.ChatModel = viper.GetString("api.chat_model")
	Api.EmbeddingModel = viper.GetString("api.embedding_model")

	Advanced.SimilarityThreshold = viper.GetFloat64("model.similarity_threshold")
	Advanced.SearchLength = viper.GetInt("model.search_length")
	Advanced.SystemPrompt = viper.GetString("model.system_prompt")

	Wiki.Dir = viper.GetString("wiki.dir")
	Wiki.Format = viper.GetStringSlice("wiki.format")
	Wiki.Exclude = viper.GetStringSlice("wiki.exclude")

	return nil
}
