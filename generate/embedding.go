package generate

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/ljcbaby/hdu-wiki-qa/database"
	"github.com/ljcbaby/hdu-wiki-qa/model"
	"github.com/ljcbaby/hdu-wiki-qa/utils"
	"github.com/oklog/ulid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func processFile(file model.FileRecord) error {
	logrus.WithField("module", "generate").WithField("file", file.FilePath).Debugf("process file")

	content, err := os.ReadFile(file.FilePath)
	if err != nil {
		return fmt.Errorf("read file failed: %w", err)
	}

	reFrontMatter := regexp.MustCompile(`(?s)^---\n(.*?)\n---\n\n`)
	str := reFrontMatter.ReplaceAllString(string(content), "")

	texts := sliceText(str, 1000)

	tx := database.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("begin transaction failed: %w", tx.Error)
	}

	if err := tx.Save(&file).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("save file record failed: %w", err)
	}

	if err := tx.Where("file_id = ?", file.Id).Delete(&model.QAPair{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("delete old qa pairs failed: %w", err)
	}

	for _, text := range texts {
		err := generateEmbedding(tx, text, file.Id)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("generate embedding failed: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("commit transaction failed: %w", err)
	}

	return nil
}

func generateEmbedding(tx *gorm.DB, str string, fileId string) error {
	var qas []model.QAResponse

	for max_retry := 5; max_retry > 0; max_retry-- {
		res, err := utils.ChatRequest("你是一个正在学习知识库的机器人，会学习输入的文段并将其转化为如下json格式的问答对输出。尽可能覆盖文段的主要信息。\n\n[{\"Q\": \"string\",\"A\": \"string\"}]",
			str)
		if err != nil {
			return fmt.Errorf("chat request failed: %w", err)
		}

		if strings.HasPrefix(res, "```json") {
			res = strings.TrimPrefix(res, "```json")
			res = strings.TrimSuffix(res, "```")
		}

		err = json.Unmarshal([]byte(res), &qas)
		if err != nil {
			logrus.WithField("module", "generate").Errorf("json unmarshal failed: %s, retry", err)
			continue
		}

		break
	}

	if len(qas) == 0 {
		return fmt.Errorf("no qa pairs generated")
	}

	for _, qa := range qas {
		pair := model.QAPair{
			Id:       ulid.MustNew(ulid.Now(), rand.Reader).String(),
			Question: qa.Question,
			Answer:   qa.Answer,
			FileID:   fileId,
		}

		var err error

		pair.QuestionEmbedding, err = utils.EmbeddingRequest(pair.Question)
		if err != nil {
			return fmt.Errorf("embedding request failed: %w", err)
		}

		pair.AnswerEmbedding, err = utils.EmbeddingRequest(pair.Answer)
		if err != nil {
			return fmt.Errorf("embedding request failed: %w", err)
		}

		if err := tx.Save(&pair).Error; err != nil {
			return fmt.Errorf("save qa pair failed: %w", err)
		}
	}

	return nil
}

func sliceText(text string, length int) []string {
	if len([]rune(text)) <= length {
		return []string{text}
	}

	texts := []string{text}

	for level := 2; level <= 6; level++ {
		cur := []string{}
		flag := false
		for _, t := range texts {
			if len([]rune(t)) <= length {
				cur = append(cur, t)
				continue
			}

			flag = true
			reHeader := regexp.MustCompile(fmt.Sprintf(`(?m)^#{%d} `, level))
			tmp := reHeader.Split(t, -1)
			if len(tmp) == 1 {
				cur = append(cur, t)
				continue
			}

			for i := 1; i < len(tmp); i++ {
				tmp[i] = strings.Repeat("#", level) + " " + tmp[i]
			}

			for i := 0; i < len(tmp)-1; i++ {
				if len([]rune(tmp[i]))+len([]rune(tmp[i+1])) <= length ||
					i == 0 && len([]rune(tmp[i])) <= 50 {
					tmp[i] += tmp[i+1]
					tmp = append(tmp[:i+1], tmp[i+2:]...)
					i--
				}
			}

			cur = append(cur, tmp...)
		}

		if !flag {
			break
		}

		texts = cur
		logrus.WithField("module", "generate").Debugf("level: %d, texts: %d", level, len(texts))
	}

	return texts
}

func cleanFile(file model.FileRecord) error {
	logrus.WithField("module", "generate").WithField("file", file.FilePath).Debugf("clean file")

	tx := database.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("begin transaction failed: %w", tx.Error)
	}

	if err := tx.Where("file_id = ?", file.Id).Delete(&model.QAPair{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("delete qa pairs failed: %w", err)
	}

	if err := tx.Delete(&file).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("delete file record failed: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("commit transaction failed: %w", err)
	}

	return nil
}
