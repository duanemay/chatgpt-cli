package cmd

import (
	"bufio"
	"context"
	"fmt"
	"github.com/pterm/pterm"
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
	setChatContext(cmd, chatContext)

	AddModelFlag(&chatFlags.model, cmd.PersistentFlags())
	AddRoleFlag(&chatFlags.role, cmd.PersistentFlags())
	AddSessionFileFlag(&chatFlags.sessionFile, cmd.PersistentFlags())
	AddSkipWriteSessionFileFlag(&chatFlags.skipWriteSessionFile, cmd.PersistentFlags())
	AddInitialSystemMessageFlag(&chatFlags.initialSystemMessage, cmd.PersistentFlags())
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
		if chatFlags.initialSystemMessage != "" {
			chatCompletionRequest.Messages = append(chatCompletionRequest.Messages, openai.ChatCompletionMessage{
				Role:    "system",
				Content: chatFlags.initialSystemMessage,
			})
		}

		reader := bufio.NewReader(os.Stdin)
		var chatRequestString string
		for {
			if chatContext.InteractiveSession {
				chatRequestString, _ = pterm.DefaultInteractiveTextInput.WithDefaultText("Enter Message").WithMultiLine().Show() // Text input with multi line enabled
			} else {
				var lines []string
				for {
					line, err := reader.ReadString('\n')
					log.WithError(err).Debugf("readString returned")
					if err != nil && err != io.EOF {
						log.WithError(err).Fatal()
					} else if err == io.EOF {
						break
					}

					lines = append(lines, line)
				}
				chatRequestString = strings.Join(lines, "\n")
			}
			if len(chatRequestString) == 0 {
				ErrorFmt.Printf("No Message to Send, exiting...\n")
				return nil
			}

			if err := sendChatMessages(chatFlags, chatContext, chatCompletionRequest, client, chatRequestString); err != nil {
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

func printBanner(f *ChatFlags) {
	TitleFmt.Printf("ChatGPT CLI v%s\n", version)
	fmt.Printf("model: %s, role: %s, temp: %0.1f, maxtok: %d, topp: %0.1f\n", f.model, f.role, f.temperature, f.maxTokens, f.topP)
	fmt.Printf("- Press TAB after entering a message to send.\n")
	fmt.Printf("- Press TAB or CTRL+C with a blank message to terminate the session without sending.\n")
}

// sendMessages sends messages to ChatGPT and prints the response
func sendChatMessages(f *ChatFlags, chatContext *ChatContext, chatCompletionRequest *openai.ChatCompletionRequest, client *openai.Client, chatRequestString string) error {
	mySpinner := pterm.DefaultSpinner
	mySpinner.Sequence = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
	mySpinner.RemoveWhenDone = true
	mySpinner.SetWriter(os.Stderr)
	successSpinner, _ := mySpinner.Start("Sending to ChatGPT, please wait...")

	chatCompletionRequest.Messages = append(chatCompletionRequest.Messages, openai.ChatCompletionMessage{
		Role:    f.role,
		Content: chatRequestString,
	})
	resp, err := client.CreateChatCompletion(context.Background(), *chatCompletionRequest)
	if err != nil {
		successSpinner.Fail(err.Error())
		return err
	}
	successSpinner.Success()

	for _, choice := range resp.Choices {
		if chatContext.InteractiveSession {
			AiFmt.Printf("\nChatGPT response:\n")
		}
		fmt.Printf("%s\n", choice.Message.Content)
	}
	chatCompletionRequest.Messages = append(chatCompletionRequest.Messages, resp.Choices[0].Message)

	return nil
}
