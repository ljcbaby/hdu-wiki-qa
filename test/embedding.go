package test

import (
	"github.com/ljcbaby/hdu-wiki-qa/utils"
	"github.com/sirupsen/logrus"
)

func Embedding() {
	v, err := utils.EmbeddingRequest("尝试生成一段文本的 embedding.")
	if err != nil {
		logrus.WithField("module", "test").WithError(err).Error("embedding request failed")
		return
	}

	logrus.WithField("module", "test").Infof("embedding response: %v", v)
}
