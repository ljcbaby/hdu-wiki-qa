package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ljcbaby/hdu-wiki-qa/conf"
	"github.com/ljcbaby/hdu-wiki-qa/model"
	"github.com/pgvector/pgvector-go"
	"github.com/sirupsen/logrus"
)

func ChatRequest(system string, user string) (string, error) {
	if user == "" {
		return "", fmt.Errorf("user input cannot be empty")
	}

	config := conf.Api

	payload := model.ChatRequest{
		Model: config.ChatModel,
	}

	if system != "" {
		message := model.Message{
			Role:    "system",
			Content: system,
		}
		payload.Messages = append(payload.Messages, message)
	}

	message := model.Message{
		Role:    "user",
		Content: user,
	}
	payload.Messages = append(payload.Messages, message)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %v", err)
	}

	logrus.WithField("module", "utils").Debugf("payload: %s", payloadBytes)

	req, err := http.NewRequest("POST", config.BaseURL+"/v1/chat/completions", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to API Endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		logrus.WithField("module", "utils").Warn("chat rate limit exceeded, waiting for 30 seconds")
		time.Sleep(30 * time.Second)
		resp, err = client.Do(req)
		if err != nil {
			return "", fmt.Errorf("failed to send request to API Endpoint: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests {
			logrus.WithField("module", "utils").Warn("chat rate limit exceeded, waiting for 1 minute")
			time.Sleep(1 * time.Minute)
			resp, err = client.Do(req)
			if err != nil {
				return "", fmt.Errorf("failed to send request to API Endpoint: %v", err)
			}
			defer resp.Body.Close()
		}
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	logrus.WithField("module", "utils").Debugf("response: %s", respBody)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	var response model.ChatResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return response.Choices[0].Message.Content, nil
}

func EmbeddingRequest(input string) (pgvector.Vector, error) {
	config := conf.Api

	payload := map[string]interface{}{
		"model": config.EmbeddingModel,
		"input": input,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return pgvector.Vector{}, fmt.Errorf("failed to marshal payload: %v", err)
	}

	logrus.WithField("module", "utils").Debugf("payload: %s", payloadBytes)

	req, err := http.NewRequest("POST", config.BaseURL+"/v1/embeddings", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return pgvector.Vector{}, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+config.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return pgvector.Vector{}, fmt.Errorf("failed to send request to API Endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		logrus.WithField("module", "utils").Warn("embedding rate limit exceeded, waiting for 30 seconds")
		time.Sleep(30 * time.Second)
		resp, err = client.Do(req)
		if err != nil {
			return pgvector.Vector{}, fmt.Errorf("failed to send request to API Endpoint: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests {
			logrus.WithField("module", "utils").Warn("embedding rate limit exceeded, waiting for 1 minute")
			time.Sleep(1 * time.Minute)
			resp, err = client.Do(req)
			if err != nil {
				return pgvector.Vector{}, fmt.Errorf("failed to send request to API Endpoint: %v", err)
			}
			defer resp.Body.Close()
		}
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return pgvector.Vector{}, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		logrus.WithField("module", "utils").Debugf("response: %s", respBody)
		return pgvector.Vector{}, fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	var response model.EmbeddingResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logrus.WithField("module", "utils").Debugf("response: %s", respBody)
		return pgvector.Vector{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return response.Data[0].Embedding.Vector, nil
}
