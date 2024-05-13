package model

import (
	"fmt"
	"time"

	"github.com/oklog/ulid"
	"gorm.io/gorm"
)

type FileRecord struct {
	Id           string `gorm:"column:id;primaryKey;not null;type:char(26)"` // ULID
	FilePath     string `gorm:"column:file_path;not null;type:varchar"`      // File path
	SHA1         string `gorm:"column:sha1;not null;type:char(40)"`          // File SHA1
	LastModified int64  `gorm:"column:last_modified;not null;type:bigint"`   // modify time
}

func (FileRecord) TableName() string {
	return "files"
}

func (f *FileRecord) BeforeSave(*gorm.DB) error {
	if _, err := ulid.Parse(f.Id); err != nil {
		return fmt.Errorf("invalid ULID: %s", err)
	}

	if len(f.SHA1) != 40 {
		return fmt.Errorf("invalid SHA1: %s", f.SHA1)
	}

	f.LastModified = time.Now().Unix()
	return nil
}
