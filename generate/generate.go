package generate

import (
	"github.com/ljcbaby/hdu-wiki-qa/model"
	"github.com/sirupsen/logrus"
)

func Init() {
	logrus.WithField("module", "generate").Info("init generate")

	var fileList []string
	listFiles(&fileList)

	var fileRecords []model.FileRecord
	checkFiles(&fileList, &fileRecords)

	logrus.WithField("module", "generate").Info("records: ", fileRecords)
}
