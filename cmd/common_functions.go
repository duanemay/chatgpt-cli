package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"time"
)

// setupOpenAIClient verifies the token exists, and creates a new OpenAI client
func setupOpenAIClient(apikey string) (*openai.Client, error) {
	if apikey == "" {
		return nil, errors.Errorf("OpenAI API Key not set")
	}
	client := openai.NewClient(apikey)
	return client, nil
}

// detectTerminal detects if the CLI is running in a terminal or not
func detectTerminal() bool {
	if terminal.IsTerminal(int(os.Stdin.Fd())) {
		return true
	}
	return false
}

// loadOrCreateChatCompletionRequest reads the sessionFile if provided, or creates a new one
func loadOrCreateChatCompletionRequest(f *ChatFlags, chatContext *ChatContext) (chat *openai.ChatCompletionRequest) {
	if f.sessionFile == "" {
		f.sessionFile = "chatgpt-cli-" + time.Now().UTC().Format(time.RFC3339) + ".json"
		if chatContext.InteractiveSession && shouldWriteSession(f) {
			fmt.Printf("  to continue with this session, use: --%s %s\n", FlagSessionFile, f.sessionFile)
		}
		chat = &openai.ChatCompletionRequest{
			Model:       f.model,
			Messages:    []openai.ChatCompletionMessage{},
			Temperature: f.temperature,
			MaxTokens:   f.maxTokens,
			TopP:        f.topP,
		}
	} else {
		chat = loadSessionFile(f.sessionFile)
		chat.Model = f.model
		chat.Temperature = f.temperature
		chat.MaxTokens = f.maxTokens
		chat.TopP = f.topP

		if chatContext.InteractiveSession {
			fmt.Printf("  continuing session from file: %s\n", f.sessionFile)
		}
	}
	return
}

// shouldWriteSession determines if the sessionFile should be written to disk
func shouldWriteSession(f *ChatFlags) bool {
	return !f.skipWriteSessionFile
}

// loadSessionFile loads the sessionFile from disk
func loadSessionFile(sessionFile string) (chat *openai.ChatCompletionRequest) {
	fileBytes, err := os.ReadFile(sessionFile)
	if err != nil {
		log.WithError(err).Fatal()
	}
	err = json.Unmarshal(fileBytes, &chat)
	if err != nil {
		log.WithError(err).Fatal()
	}
	return
}

// writeSessionFile writes the sessionFile to disk
func writeSessionFile(f *ChatFlags, chat *openai.ChatCompletionRequest) {
	objJson, err := json.MarshalIndent(chat, "", "  ")
	if err != nil {
		log.WithError(err).Fatal()
	}
	err = os.WriteFile(f.sessionFile, objJson, os.ModePerm)
	if err != nil {
		log.WithError(err).Fatal()
	}
}
