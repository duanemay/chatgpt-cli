package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"time"
)

// setupOpenAIClient verifies the token exists, and creates a new OpenAI client
func setupOpenAIClient() {
	if apikey == "" {
		log.Fatal("OpenAI API Key not set")
	}
	client = openai.NewClient(apikey)
}

// detectTerminal detects if the CLI is running in a terminal or not
func detectTerminal() bool {
	if terminal.IsTerminal(int(os.Stdin.Fd())) {
		interactiveSession = true
	}
	return interactiveSession
}

// setupSessionFile reads the sessionFile if provided, or creates a new one
func setupSessionFile() {
	if !viper.IsSet(FlagSessionFile) {
		sessionId := time.Now().Unix()
		sessionFile = fmt.Sprintf("chatgpt-cli-%d.json", sessionId)
		if interactiveSession {
			fmt.Printf("  to continue with this session, use: --%s %s\n", FlagSessionFile, sessionFile)
		}
		chat = openai.ChatCompletionRequest{
			Model:       model,
			Messages:    []openai.ChatCompletionMessage{},
			Temperature: temperature,
			MaxTokens:   maxTokens,
			TopP:        topP,
		}
	} else {
		loadSessionFile()
		chat.Model = model
		chat.Temperature = temperature
		chat.MaxTokens = maxTokens
		chat.TopP = topP

		if interactiveSession {
			fmt.Printf("  continuing session from file: %s\n", sessionFile)
		}
	}
}

func loadSessionFile() {
	fileBytes, err := os.ReadFile(sessionFile)
	if err != nil {
		log.WithError(err).Fatal()
	}
	err = json.Unmarshal(fileBytes, &chat)
	if err != nil {
		log.WithError(err).Fatal()
	}
}

// writeSessionFile writes the sessionFile to disk
func writeSessionFile(chat openai.ChatCompletionRequest) {
	objJson, err := json.MarshalIndent(chat, "", "  ")
	if err != nil {
		log.WithError(err).Fatal()
	}
	err = os.WriteFile(sessionFile, objJson, os.ModePerm)
	if err != nil {
		log.WithError(err).Fatal()
	}
}
