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

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Enter a chat session with ChatGPT",
	Long:  "Enter a chat session with ChatGPT",
	Run:   chatCmdRunner,
}

func init() {
	rootCmd.AddCommand(chatCmd)

	chatCmd.Flags().StringVarP(&model, FlagModel, "m", defaultModel, "ChatGPT Model")
	chatCmd.Flags().StringVarP(&role, FlagRole, "r", defaultRole, "ChatGPT Role")
	chatCmd.Flags().Float32VarP(&temperature, FlagTemperature, "t", defaultTemperature, "temperature, between 0 and 2. Higher values make the output more random")
	chatCmd.Flags().IntVar(&maxTokens, FlagMaxTokens, defaultMaxTokens, "number of tokens to generate = $")
	chatCmd.Flags().Float32Var(&topP, FlagTopP, defaultTopP, "results of the tokens with top_p probability mass")
	chatCmd.Flags().StringVarP(&sessionFile, FlagSessionFile, "s", "", "Continue a session from a file")
	chatCmd.Flags().StringVar(&eomMarker, FlagEomMarker, defaultEomMarker, "Text to enter to mark the end of a message to send to ChatGPT")
	chatCmd.Flags().StringVar(&eosMarker, FlagEosMarker, defaultEosMarker, "Text to enter to end of a session with ChatGPT")

	_ = chatCmd.MarkPersistentFlagRequired(FlagApiKey)
}

func chatCmdRunner(_ *cobra.Command, _ []string) {
	log.Debugf("chatCmd called")
	detectTerminal()
	if interactiveSession {
		_, _ = narratorFmt.Printf("ChatGPT CLI v%s\n", version)
		fmt.Printf("model: %s, role: %s, temp: %0.1f, maxtok: %d, topp: %0.1f\n", model, role, temperature, maxTokens, topP)
		fmt.Printf("- Press CTRL+D or '%s' on a separate line to send.\n", eomMarker)
		fmt.Printf("- Press CTRL+C or enter '%s' on a separate line to terminate the session without sending.\n", eosMarker)
	}
	setupOpenAIClient()
	setupSessionFile()
	reader := bufio.NewReader(os.Stdin)

	for {
		if interactiveSession {
			_, _ = humanFmt.Printf("\nEnter Message")
			fmt.Printf(" (%s to send; %s to exit):\n", eomMarker, eosMarker)
		}
		var lines []string
		for {
			line, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				log.WithError(err).Fatal()
			} else if err == io.EOF {
				break
			}

			trimmedLine := strings.TrimSpace(line)
			if trimmedLine == eomMarker {
				break
			} else if trimmedLine == eosMarker {
				return
			}

			lines = append(lines, line)
		}

		if err := sendMessages(&chat, lines); err != nil {
			log.WithError(err).Fatal()
		}
		writeSessionFile(chat)

		if !interactiveSession {
			break
		}
	}
}

// sendMessages sends messages to ChatGPT and prints the response
func sendMessages(chat *openai.ChatCompletionRequest, lines []string) error {
	if interactiveSession {
		_, _ = narratorFmt.Println("\nSending to ChatGPT, please wait...")
	}

	chat.Messages = append(chat.Messages, openai.ChatCompletionMessage{
		Role:    role,
		Content: strings.Join(lines, "\n"),
	})
	resp, err := client.CreateChatCompletion(context.Background(), *chat)

	if err != nil {
		return err
	}

	for _, choice := range resp.Choices {
		if interactiveSession {
			_, _ = aiFmt.Printf("\nChatGPT response:\n")
		}
		fmt.Printf("%s\n", choice.Message.Content)
	}
	chat.Messages = append(chat.Messages, resp.Choices[0].Message)

	return nil
}
