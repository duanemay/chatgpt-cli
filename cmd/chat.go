package cmd

import (
	"bufio"
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
)

func NewChatCmd(rootFlags *RootFlags) *cobra.Command {
	chatFlags := NewChatFlags()
	chatContext := NewChatContext()
	var cmd = &cobra.Command{
		Use:   "chat",
		Short: "Enter a chat session with ChatGPT",
		Long:  "Enter a chat session with ChatGPT",
		RunE:  chatCmdRun(rootFlags, chatFlags, chatContext),
	}
	cmd.SetContext(context.WithValue(context.Background(), "chatContext", chatContext))

	AddModelFlag(&chatFlags.model, cmd.PersistentFlags())
	AddRoleFlag(&chatFlags.role, cmd.PersistentFlags())
	AddSessionFileFlag(&chatFlags.sessionFile, cmd.PersistentFlags())
	AddSkipWriteSessionFileFlag(&chatFlags.skipWriteSessionFile, cmd.PersistentFlags())
	AddEomMarkerFlag(&chatFlags.eomMarker, cmd.PersistentFlags())
	AddEosMarkerFlag(&chatFlags.eosMarker, cmd.PersistentFlags())
	AddTemperatureFlag(&chatFlags.temperature, cmd.PersistentFlags())
	AddMaxTokensFlag(&chatFlags.maxTokens, cmd.PersistentFlags())
	AddTopPFlag(&chatFlags.topP, cmd.PersistentFlags())
	_ = cmd.MarkPersistentFlagRequired(FlagApiKey)

	return cmd
}

func chatCmdRun(rootFlags *RootFlags, chatFlags *ChatFlags, chatContext *ChatContext) func(cmd *cobra.Command, args []string) error {
	return func(_ *cobra.Command, _ []string) error {
		log.Debugf("chatCmd called")
		chatContext.InteractiveSession = detectTerminal()
		if chatContext.InteractiveSession {
			printBanner(chatFlags)
		}
		client, err := setupOpenAIClient(rootFlags.apikey)
		if err != nil {
			log.WithError(err).Fatal()
		}

		chatCompletionRequest := loadOrCreateChatCompletionRequest(chatFlags, chatContext)
		reader := bufio.NewReader(os.Stdin)

		for {
			if chatContext.InteractiveSession {
				printPrompt(chatFlags.eomMarker, chatFlags.eosMarker)
			}
			var lines []string
			for {
				line, err := reader.ReadString('\n')
				log.WithError(err).Debugf("readString returned")
				if err != nil && err != io.EOF {
					log.WithError(err).Fatal()
				} else if err == io.EOF {
					break
				}

				trimmedLine := strings.TrimSpace(line)
				if trimmedLine == chatFlags.eomMarker {
					break
				} else if trimmedLine == chatFlags.eosMarker {
					return nil
				}

				lines = append(lines, line)
			}
			if len(lines) == 0 {
				log.Warning("No Message to Send")
				continue
			}

			if err := sendMessages(chatFlags, chatContext, chatCompletionRequest, client, lines); err != nil {
				log.WithError(err).Fatal()
			}

			if shouldWriteSession(chatFlags) {
				writeSessionFile(chatFlags, chatCompletionRequest)
			}

			if !chatContext.InteractiveSession {
				break
			}
		}
		return nil
	}
}

func printPrompt(eomMarker, eosMarker string) {
	_, _ = HumanFmt.Printf("\nEnter Message")
	fmt.Printf(" (%s to send; %s to exit):\n", eomMarker, eosMarker)
}

func printBanner(f *ChatFlags) {
	_, _ = NarratorFmt.Printf("ChatGPT CLI v%s\n", version)
	fmt.Printf("model: %s, role: %s, temp: %0.1f, maxtok: %d, topp: %0.1f\n", f.model, f.role, f.temperature, f.maxTokens, f.topP)
	fmt.Printf("- Press CTRL+D or '%s' on a separate line to send.\n", f.eomMarker)
	fmt.Printf("- Press CTRL+C or enter '%s' on a separate line to terminate the session without sending.\n", f.eosMarker)
}

// sendMessages sends messages to ChatGPT and prints the response
func sendMessages(f *ChatFlags, chatContext *ChatContext, chatCompletionRequest *openai.ChatCompletionRequest, client *openai.Client, lines []string) error {
	if chatContext.InteractiveSession {
		_, _ = narratorFmt.Println("\nSending to ChatGPT, please wait...")
	}

	chatCompletionRequest.Messages = append(chatCompletionRequest.Messages, openai.ChatCompletionMessage{
		Role:    f.role,
		Content: strings.Join(lines, "\n"),
	})
	resp, err := client.CreateChatCompletion(context.Background(), *chatCompletionRequest)
	if err != nil {
		return err
	}

	for _, choice := range resp.Choices {
		if chatContext.InteractiveSession {
			_, _ = AiFmt.Printf("\nChatGPT response:\n")
		}
		fmt.Printf("%s\n", choice.Message.Content)
	}
	chatCompletionRequest.Messages = append(chatCompletionRequest.Messages, resp.Choices[0].Message)

	return nil
}
