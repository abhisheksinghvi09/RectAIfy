package config

import "errors"

var (
	ErrMissingOpenAIKey = errors.New("OPENAI_API_KEY environment variable is required")
)
