package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/pterm/pterm"
	"github.com/sashabaranov/go-openai"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewTranscriptionCmd(rootFlags *RootFlags) *cobra.Command {
	transcriptionFlags := NewTranscriptionFlags()
	chatContext := NewChatContext()
	var cmd = &cobra.Command{
		Use:   "transcribe",
		Short: "Transcribe audio to text",
		Long:  "Transcribe audio to text",
		RunE:  transcribeCmdRunner(rootFlags, transcriptionFlags, chatContext),
	}
	setChatContext(cmd, chatContext)

	AddModelFlag(&transcriptionFlags.model, cmd.PersistentFlags())
	AddLanguageFlag(&transcriptionFlags.language, cmd.PersistentFlags())
	AddInputFileFlag(&transcriptionFlags.inputFiles, cmd.PersistentFlags())
	AddInitialSystemMessageFlag(&transcriptionFlags.initialSystemMessage, cmd.PersistentFlags())
	_ = cmd.MarkPersistentFlagRequired(FlagApiKey)
	_ = cmd.MarkPersistentFlagRequired(FlagInputFile)

	return cmd
}

func transcribeCmdRunner(rootFlags *RootFlags, transcriptionFlags *TranscriptionFlags, chatContext *ChatContext) func(cmd *cobra.Command, args []string) error {
	return func(_ *cobra.Command, _ []string) error {
		log.Debugf("transcribeCmd called")
		err := transcriptionFlags.ValidateFlags()
		if err != nil {
			log.WithError(err).Fatal()
		}

		chatContext.InteractiveSession = detectTerminal()
		if chatContext.InteractiveSession {
			printTranscriptionBanner(transcriptionFlags)
		}
		client, err := setupOpenAIClient(rootFlags.apikey)
		if err != nil {
			log.WithError(err).Fatal()
		}

		if err := sendTranscriptionMessages(transcriptionFlags, chatContext, client); err != nil {
			log.WithError(err).Fatal()
		}
		return nil
	}
}

func printTranscriptionBanner(f *TranscriptionFlags) {
	TitleFmt.Printf("ChatGPT CLI v%s\n", version)
	fmt.Printf("Transcription Model: %s\n", f.model)
	fmt.Printf("- Press TAB after entering a message to send.\n")
	fmt.Printf("- Press TAB or CTRL+C with a blank message to terminate the session without sending.\n")
}

// sendVisionMessages sends messages to ChatGPT and prints the response
func sendTranscriptionMessages(f *TranscriptionFlags, chatContext *ChatContext, client *openai.Client) error {
	mySpinner := pterm.DefaultSpinner
	mySpinner.Sequence = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
	mySpinner.RemoveWhenDone = true
	mySpinner.Writer = os.Stderr

	for _, file := range f.inputFiles {
		successSpinner, _ := mySpinner.Start("Sending to ChatGPT, please wait...")

		resp, err := client.CreateTranscription(
			context.Background(),
			openai.AudioRequest{
				Prompt:   f.initialSystemMessage,
				Language: f.language,
				Model:    openai.Whisper1,
				FilePath: file,
			},
		)
		if err != nil {
			successSpinner.Fail(err.Error())
			return err
		}
		successSpinner.Success()

		if chatContext.InteractiveSession {
			AiFmt.Printf("\nChatGPT response:\n")
		}
		fmt.Printf("%s\n", resp.Text)
	}
	return nil
}
