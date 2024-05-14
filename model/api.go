package model

import (
	"encoding/json"

	"github.com/pgvector/pgvector-go"
)

type ChatResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int                    `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChoiceItem           `json:"choices"`
	Usage   map[string]interface{} `json:"usage"`
}

type ChoiceItem struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type EmbeddingResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Index     int       `json:"index"`
		Embedding Embedding `json:"embedding"`
	} `json:"data"`
	Model string                 `json:"model"`
	Usage map[string]interface{} `json:"usage"`
}

type Embedding struct {
	pgvector.Vector
}

func (e *Embedding) UnmarshalJSON(data []byte) error {
	var arr []float32
	err := json.Unmarshal(data, &arr)
	if err != nil {
		return err
	}
	e.Vector = pgvector.NewVector(arr)
	return nil
}
