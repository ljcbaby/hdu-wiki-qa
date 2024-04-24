package conf

import "errors"

var (
	ErrEmptyHost     error = errors.New("empty host")
	ErrEmptyPort     error = errors.New("empty port")
	ErrEmptyUser     error = errors.New("empty user")
	ErrEmptyPassword error = errors.New("empty password")
	ErrEmptyDatabase error = errors.New("empty database")

	ErrEmptyBaseURL        error = errors.New("empty base url")
	ErrEmptyApiKey         error = errors.New("empty api key")
	ErrEmptyChatModel      error = errors.New("empty chat model")
	ErrEmptyEmbeddingModel error = errors.New("empty embedding model")

	ErrEmptySimilarityThreshold error = errors.New("empty similarity threshold")
	ErrEmptySearchLength        error = errors.New("empty search length")
	ErrEmptySystemPrompt        error = errors.New("empty system prompt")

	ErrEmptyWikiDir     error = errors.New("empty wiki dir")
	ErrEmptyWikiFormat  error = errors.New("empty wiki format")
	ErrEmptyWikiExclude error = errors.New("empty wiki exclude")
)

func checkNecessary() error {
	if Pgsql.Host == "" {
		return ErrEmptyHost
	}
	if Pgsql.Port == 0 {
		return ErrEmptyPort
	}
	if Pgsql.User == "" {
		return ErrEmptyUser
	}
	if Pgsql.Password == "" {
		return ErrEmptyPassword
	}
	if Pgsql.Database == "" {
		return ErrEmptyDatabase
	}
	if Api.BaseURL == "" {
		return ErrEmptyBaseURL
	}
	if Api.ApiKey == "" {
		return ErrEmptyApiKey
	}
	if Api.ChatModel == "" {
		return ErrEmptyChatModel
	}
	if Api.EmbeddingModel == "" {
		return ErrEmptyEmbeddingModel
	}
	if Advanced.SimilarityThreshold == 0 {
		return ErrEmptySimilarityThreshold
	}
	if Advanced.SearchLength == 0 {
		return ErrEmptySearchLength
	}
	if Advanced.SystemPrompt == "" {
		return ErrEmptySystemPrompt
	}
	if Wiki.Dir == "" {
		return ErrEmptyWikiDir
	}
	if len(Wiki.format) == 0 {
		return ErrEmptyWikiFormat
	}
	if len(Wiki.exclude) == 0 {
		return ErrEmptyWikiExclude
	}
	return nil
}
