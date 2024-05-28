package generate

import (
	"crypto/rand"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/ljcbaby/hdu-wiki-qa/database"
	"github.com/ljcbaby/hdu-wiki-qa/model"
	"github.com/oklog/ulid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func checkFiles(fileList *[]string, fileRecords *[]model.FileRecord, rmRecords *[]model.FileRecord) {
	db := database.DB

	if err := db.Find(&rmRecords).Error; err != nil {
		logrus.WithField("module", "generate").WithError(err).Fatalf("query file records failed")
	}

	for _, file := range *fileList {
		record := new(model.FileRecord)
		if err := db.Where("file_path = ?", file).First(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logrus.WithField("module", "generate").Infof("file %s not found in database", file)
				fileRecord, err := newFileRecord(file)
				if err != nil {
					logrus.WithField("module", "generate").WithError(err).Fatalf("new file record failed")
					return
				}
				*fileRecords = append(*fileRecords, *fileRecord)
				continue
			}

			logrus.WithField("module", "generate").WithError(err).Fatalf("query file record failed")
		}

		osFile, err := os.Open(file)
		if err != nil {
			logrus.WithField("module", "generate").WithError(err).Fatalf("open file %s failed", file)
			return
		}

		fileContents, err := io.ReadAll(osFile)
		if err != nil {
			logrus.WithField("module", "generate").WithError(err).Fatalf("read file %s failed", file)
			return
		}

		hash := fmt.Sprintf("%x", sha1.Sum(fileContents))
		if hash == record.SHA1 {
			logrus.WithField("module", "generate").Infof("file %s is not modified", file)
		} else {
			logrus.WithField("module", "generate").Infof("file %s is modified", file)
			record.SHA1 = hash
			*fileRecords = append(*fileRecords, *record)
		}

		for i, rmRecord := range *rmRecords {
			if rmRecord.Id == record.Id {
				*rmRecords = append((*rmRecords)[:i], (*rmRecords)[i+1:]...)
				break
			}
		}
	}
}

func newFileRecord(file string) (*model.FileRecord, error) {
	r := new(model.FileRecord)
	r.Id = ulid.MustNew(ulid.Now(), rand.Reader).String()
	r.FilePath = file

	osFile, err := os.Open(file)
	if err != nil {
		logrus.WithField("module", "generate").WithError(err).Fatalf("open file %s failed", file)
		return nil, err
	}

	fileContents, err := io.ReadAll(osFile)
	if err != nil {
		logrus.WithField("module", "generate").WithError(err).Fatalf("read file %s failed", file)
		return nil, err
	}

	hash := fmt.Sprintf("%x", sha1.Sum(fileContents))
	r.SHA1 = hash

	return r, nil
}
