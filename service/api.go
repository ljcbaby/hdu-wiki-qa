package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ljcbaby/hdu-wiki-qa/conf"
	"github.com/ljcbaby/hdu-wiki-qa/database"
	"github.com/ljcbaby/hdu-wiki-qa/model"
	"github.com/ljcbaby/hdu-wiki-qa/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

func Chat(c *gin.Context) {
	var req model.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithField("module", "service").Errorf("Failed to bind json: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	q := req.Messages[len(req.Messages)-1].Content
	str, err := getKnowledge(q)
	if err != nil {
		if err.Error() == "ESCAPE" {
			c.JSON(http.StatusOK, model.ApiResponse{
				RawContent: "未查到相关知识。",
				Reply:      "抱歉，我好像还不会回答这个问题。",
			})
			return
		}
		logrus.WithField("module", "service").Errorf("Failed to get knowledge: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	config := conf.Advanced

	res, err := utils.ChatRequest(config.SystemPrompt, str)
	if err != nil {
		logrus.WithField("module", "service").Errorf("Failed to chat: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, model.ApiResponse{
		RawContent: str,
		Reply:      res,
	})
}

func getKnowledge(q string) (string, error) {
	config := conf.Advanced
	db := database.DB

	v, err := utils.EmbeddingRequest(q)
	if err != nil {
		return "", fmt.Errorf("get embedding: %v", err)
	}

	var QA4Q []model.QAPair
	var QA4A []model.QAPair

	err = db.Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "cosine_distance(question_embedding,?)", Vars: []interface{}{v.String()}},
	}).Clauses(clause.Where{
		Exprs: []clause.Expression{clause.Expr{SQL: "cosine_distance(question_embedding,?) < ?", Vars: []interface{}{v.String(), config.SimilarityThreshold}}},
	}).Limit(config.SearchLength).Find(&QA4Q).Error
	if err != nil {
		return "", fmt.Errorf("get Question Knowledge: %v", err)
	}

	for _, qa := range QA4Q {
		logrus.WithField("module", "service").Debugf("QA4Q ID: %v", qa.Id)
	}

	err = db.Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "cosine_distance(answer_embedding,?)", Vars: []interface{}{v.String()}},
	}).Clauses(clause.Where{
		Exprs: []clause.Expression{clause.Expr{SQL: "cosine_distance(answer_embedding,?) < ?", Vars: []interface{}{v.String(), config.SimilarityThreshold}}},
	}).Limit(config.SearchLength).Find(&QA4A).Error
	if err != nil {
		return "", fmt.Errorf("get Answer Knowledge: %v", err)
	}

	for _, qa := range QA4A {
		logrus.WithField("module", "service").Debugf("QA4A ID: %v", qa.Id)
	}

	if len(QA4Q) == 0 && len(QA4A) == 0 {
		return "", fmt.Errorf("ESCAPE")
	}

	str := "用户的问题为：" + q + "\n\n"
	str += "以下是参考问答及其来源：\n"

	for _, qa := range QA4Q {
		str += "Q: " + qa.Question + "\n"
		str += "A: " + qa.Answer + "\n"

		f := model.FileRecord{
			Id: qa.FileID,
		}
		db.First(&f)

		str += "来源：杭电导航/" + f.FilePath[:len(f.FilePath)-3] + "\n\n"
	}

	for _, qa := range QA4A {
		str += "Q: " + qa.Question + "\n"
		str += "A: " + qa.Answer + "\n"

		f := model.FileRecord{
			Id: qa.FileID,
		}
		db.First(&f)

		str += "来源：杭电导航/" + f.FilePath[:len(f.FilePath)-3] + "\n\n"
	}

	return str, nil
}
