package cmd

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/sashabaranov/go-openai"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
)

func NewVisionCmd(rootFlags *RootFlags) *cobra.Command {
	visionFlags := NewVisionFlags()
	chatContext := NewChatContext()
	var cmd = &cobra.Command{
		Use:   "vision",
		Short: "Answer questions based on images",
		Long:  "Answer questions based on images",
		RunE:  visionCmdRunner(rootFlags, visionFlags, chatContext),
	}
	setChatContext(cmd, chatContext)

	AddModelFlag(&visionFlags.model, cmd.PersistentFlags())
	AddDetailFlag(&visionFlags.DetailStr, cmd.PersistentFlags())
	AddRoleFlag(&visionFlags.role, cmd.PersistentFlags())
	AddInputFileFlag(&visionFlags.inputFiles, cmd.PersistentFlags())
	AddSessionFileFlag(&visionFlags.sessionFile, cmd.PersistentFlags())
	AddSkipWriteSessionFileFlag(&visionFlags.skipWriteSessionFile, cmd.PersistentFlags())
	AddInitialSystemMessageFlag(&visionFlags.initialSystemMessage, cmd.PersistentFlags())
	_ = cmd.MarkPersistentFlagRequired(FlagApiKey)
	_ = cmd.MarkPersistentFlagRequired(FlagInputFile)

	return cmd
}

func visionCmdRunner(rootFlags *RootFlags, visionFlags *VisionFlags, chatContext *ChatContext) func(cmd *cobra.Command, args []string) error {
	return func(_ *cobra.Command, _ []string) error {
		log.Debugf("visionCmd called")
		err := visionFlags.ValidateFlags()
		if err != nil {
			log.WithError(err).Fatal()
		}

		chatContext.InteractiveSession = detectTerminal()
		if chatContext.InteractiveSession {
			printVisionBanner(visionFlags)
		}
		client, err := setupOpenAIClient(rootFlags.apikey)
		if err != nil {
			log.WithError(err).Fatal()
		}

		chatFlags := ChatFlagsFromVisionFlags(visionFlags)
		chatCompletionRequest := loadOrCreateChatCompletionRequest(chatFlags, chatContext)
		if visionFlags.initialSystemMessage != "" {
			chatCompletionRequest.Messages = append(chatCompletionRequest.Messages, openai.ChatCompletionMessage{
				Role:    "system",
				Content: visionFlags.initialSystemMessage,
			})
		}

		reader := bufio.NewReader(os.Stdin)
		var chatRequestString string
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

		if err := sendVisionMessages(visionFlags, chatContext, chatCompletionRequest, client, chatRequestString); err != nil {
			log.WithError(err).Fatal()
		}
		return nil
	}
}

func printVisionBanner(f *VisionFlags) {
	TitleFmt.Printf("ChatGPT CLI v%s\n", version)
	fmt.Printf("Model: %s, detail: %s\n", f.model, f.Detail)
	fmt.Printf("- Press TAB after entering a message to send.\n")
	fmt.Printf("- Press TAB or CTRL+C with a blank message to terminate the session without sending.\n")
}

// sendVisionMessages sends messages to ChatGPT and prints the response
func sendVisionMessages(f *VisionFlags, chatContext *ChatContext, chatCompletionRequest *openai.ChatCompletionRequest, client *openai.Client, chatRequestString string) error {
	mySpinner := pterm.DefaultSpinner
	mySpinner.Sequence = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
	mySpinner.RemoveWhenDone = true
	mySpinner.Writer = os.Stderr
	successSpinner, _ := mySpinner.Start("Sending to ChatGPT, please wait...")

	content := []openai.ChatMessagePart{
		{
			Type: openai.ChatMessagePartTypeText,
			Text: chatRequestString,
		},
	}

	for _, file := range f.inputFiles {
		image, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", file, err)
		}

		// Encode the image to base64
		encodedImage := base64.StdEncoding.EncodeToString(image)

		content = append(content, openai.ChatMessagePart{
			Type: openai.ChatMessagePartTypeImageURL,
			ImageURL: &openai.ChatMessageImageURL{
				URL:    "data:image/jpeg;base64," + encodedImage,
				Detail: f.Detail,
			},
		})
	}

	chatCompletionRequest.Messages = append(chatCompletionRequest.Messages, openai.ChatCompletionMessage{
		Role:         f.role,
		MultiContent: content,
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
