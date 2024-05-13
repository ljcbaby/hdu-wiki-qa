package generate

import "github.com/sirupsen/logrus"

func Init() {
	logrus.Info("init generate")

	var fileList []string
	listFiles(&fileList)

	logrus.Info("file list: ", fileList)
}
