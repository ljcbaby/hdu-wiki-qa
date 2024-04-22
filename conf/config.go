package conf

import (
	_ "embed"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//go:embed config.example.yaml
var configSample []byte

var Config struct {
	pgsql struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"postgresql"`

	api struct {
		BaseURL        string `yaml:"base_url"`
		ApiKey         string `yaml:"api_key"`
		ChatModel      string `yaml:"chatbot_model"`
		EmbeddingModel string `yaml:"embedding_model"`
	} `yaml:"api"`

	advanced struct {
		SimilarityThreshold float64 `yaml:"similarity_threshold"`
		SearchLength        int     `yaml:"search_length"`
		SystemPrompt        string  `yaml:"system_prompt"`
	} `yaml:"model"`

	wiki struct {
		Dir     string   `yaml:"dir"`
		format  []string `yaml:"format"`
		exclude []string `yaml:"exclude"`
	} `yaml:"wiki"`
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
	err := viper.ReadInConfig()

	if err != nil {
		logrus.WithError(err).Fatalf("read config from %s failed", file)
	}

	err = viper.Unmarshal(&Config)
	if err != nil {
		logrus.WithError(err).Fatalf("unmarshal config failed")
	}
}
