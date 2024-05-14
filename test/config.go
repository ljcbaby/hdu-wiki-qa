package test

import (
	"github.com/ljcbaby/hdu-wiki-qa/conf"
	"github.com/sirupsen/logrus"
)

func Config(path string) {
	logrus.WithField("module", "test").Debugf("load config from %s", path)

	pgsql := conf.Pgsql
	logrus.WithField("module", "test").Infof("pgsql config: %v", pgsql)

	api := conf.Api
	logrus.WithField("module", "test").Infof("api config: %v", api)

	advanced := conf.Advanced
	logrus.WithField("module", "test").Infof("advanced config: %v", advanced)

	wiki := conf.Wiki
	logrus.WithField("module", "test").Infof("wiki config: %v", wiki)
}
