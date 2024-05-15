package test

import (
	"github.com/ljcbaby/hdu-wiki-qa/utils"
	"github.com/sirupsen/logrus"
)

func Chat() {
	s, err := utils.ChatRequest("", "你是什么模型？")
	if err != nil {
		logrus.WithField("module", "test").WithError(err).Error("chat request failed")
		return
	}

	logrus.WithField("module", "test").Infof("chat response: %s", s)
}
