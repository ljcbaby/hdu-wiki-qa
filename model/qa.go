package model

import (
	"fmt"

	"github.com/oklog/ulid"
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

type QApair struct {
	Id                string          `gorm:"column:id;primaryKey;not null;type:char(26)"` // ULID
	Question          string          `gorm:"column:question;not null;type:varchar"`
	Answer            string          `gorm:"column:answer;not null;type:varchar"`
	QuestionEmbedding pgvector.Vector `gorm:"column:question_embedding;not null;type:vector(1536)"`
	AnswerEmbedding   pgvector.Vector `gorm:"column:answer_embedding;not null;type:vector(1536)"`
	FileID            string          `gorm:"column:file_id;not null;type:char(26)"` // File ULID
}

func (QApair) TableName() string {
	return "qa_items"
}

func (qa *QApair) BeforeSave(db *gorm.DB) error {
	if _, err := ulid.Parse(qa.Id); err != nil {
		return fmt.Errorf("invalid ULID: %s", err)
	}

	if _, err := ulid.Parse(qa.FileID); err != nil {
		return fmt.Errorf("invalid File ULID: %s", err)
	}

	// Check if the foreign key exists
	var file FileRecord
	result := db.First(&file, "id = ?", qa.FileID)
	if result.Error != nil {
		return fmt.Errorf("foreign key does not exist: %s", result.Error)
	}

	return nil
}

type QAResponse struct {
	Question string `json:"Q"`
	Answer   string `json:"A"`
}
