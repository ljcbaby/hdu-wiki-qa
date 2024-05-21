package database

import (
	"fmt"
	"os"

	"github.com/ljcbaby/hdu-wiki-qa/conf"
	"github.com/ljcbaby/hdu-wiki-qa/model"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func Connect() error {
	conf := conf.Pgsql

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		conf.Host, conf.User, conf.Password, conf.Database, conf.Port)

	var err error
	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: &LogrusLogger{logger: logrus.StandardLogger()},
	})
	if err != nil {
		logrus.WithField("module", "database").WithError(err).Fatalf("connect to db failed")
		os.Exit(1)
	}

	if err := DB.AutoMigrate(&model.QAPair{}, &model.FileRecord{}); err != nil {
		logrus.WithField("module", "database").WithError(err).Fatalf("auto migrate failed")
		os.Exit(1)
	}

	return nil
}
