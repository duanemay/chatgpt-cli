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

// loadOrCreateChatCompletionRequest
// if a sessionFile is provided, and it exists, then it is loaded
// if a sessionFile is provided, and it does not exist, then it is created
// if a sessionFile is not provided, then a new session and sessionFile is created
func loadOrCreateChatCompletionRequest(f *ChatFlags, chatContext *ChatContext) *openai.ChatCompletionRequest {
	var chat *openai.ChatCompletionRequest = nil

	// if a sessionFile is provided, check if it exists
	if f.sessionFile != "" {
		if _, err := os.Stat(f.sessionFile); err == nil {
			// sessionFile provided and exists, load it
			chat = loadSessionFile(f.sessionFile)
			if chatContext.InteractiveSession {
				fmt.Printf("  continuing session from file: %s\n", f.sessionFile)
			}

			// update the chat parameters from the flags
			chat.Model = f.model
			chat.Temperature = f.temperature
			chat.MaxTokens = f.maxTokens
			chat.TopP = f.topP
		}
	} else {
		// no sessionFile provided, create a new sessionFile name
		f.sessionFile = "chatgpt-cli-" + time.Now().UTC().Format(time.RFC3339) + ".json"
	}

	// if a sessionFile was not provided, or it did not exist, create a new session
	if chat == nil {
		chat = &openai.ChatCompletionRequest{
			Model:       f.model,
			Messages:    []openai.ChatCompletionMessage{},
			Temperature: f.temperature,
			MaxTokens:   f.maxTokens,
			TopP:        f.topP,
		}
		if chatContext.InteractiveSession && shouldWriteSession(f) {
			fmt.Printf("  to continue with this session, use: --%s %s\n", FlagSessionFile, f.sessionFile)
		}
	}

	return chat
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
