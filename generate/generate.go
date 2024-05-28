package generate

import (
	"github.com/ljcbaby/hdu-wiki-qa/model"
	"github.com/sirupsen/logrus"
)

func Init() {
	logrus.WithField("module", "generate").Infof("init generate")

	var fileList []string
	listFiles(&fileList)

	var fileRecords []model.FileRecord
	var rmRecords []model.FileRecord
	checkFiles(&fileList, &fileRecords, &rmRecords)

	logrus.WithField("module", "generate").Debugf("records: %v", fileRecords)
	logrus.WithField("module", "generate").Debugf("rm records: %v", rmRecords)

	for _, record := range fileRecords {
		err := processFile(record)
		if err != nil {
			logrus.WithField("module", "generate").WithField("file", record.FilePath).WithError(err).Error("generate embedding failed")
		}
	}

	for _, record := range rmRecords {
		err := cleanFile(record)
		if err != nil {
			logrus.WithField("module", "generate").WithField("file", record.FilePath).WithError(err).Error("remove file failed")
		}
	}
}
