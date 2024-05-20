package cmd

import (
	"bufio"
	"context"
	"fmt"
	os2 "github.com/duanemay/chatgpt-cli/pkg/os"
	"github.com/pterm/pterm"
	"github.com/sashabaranov/go-openai"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
	"time"
)

func NewSpeechCmd(rootFlags *RootFlags) *cobra.Command {
	speechFlags := NewSpeechFlags()
	chatContext := NewChatContext()
	var cmd = &cobra.Command{
		Use:     "text-to-speech",
		Aliases: []string{"speech"},
		Short:   "Text to speech, creates an audio file",
		Long:    "Text to speech, creates an audio file",
		RunE:    speechCmdRunner(rootFlags, speechFlags, chatContext),
	}
	cmd.SetContext(context.WithValue(context.Background(), "chatContext", chatContext))

	AddSpeechModelFlag(&speechFlags.ModelStr, cmd.PersistentFlags())
	AddSpeedFlag(&speechFlags.Speed, cmd.PersistentFlags())
	AddVoiceFlag(&speechFlags.VoiceStr, cmd.PersistentFlags())
	AddImageOutputPrefixFlag(&speechFlags.OutputPrefix, "tts-"+time.Now().UTC().Format(time.RFC3339), cmd.PersistentFlags())
	_ = cmd.MarkPersistentFlagRequired("apikey")

	return cmd
}

func speechCmdRunner(rootFlags *RootFlags, speechFlags *SpeechFlags, chatContext *ChatContext) func(cmd *cobra.Command, args []string) error {
	return func(_ *cobra.Command, _ []string) error {
		log.Debugf("speechCmd called")
		err := speechFlags.ValidateFlags()
		if err != nil {
			log.WithError(err).Fatal()
		}

		chatContext.InteractiveSession = detectTerminal()
		if chatContext.InteractiveSession {
			printSpeechBanner(speechFlags)
		}
		client, err := setupOpenAIClient(rootFlags.apikey)
		if err != nil {
			log.WithError(err).Fatal()
		}

		reader := bufio.NewReader(os.Stdin)
		var chatRequestString string
		for {
			if chatContext.InteractiveSession {
				chatRequestString, _ = pterm.DefaultInteractiveTextInput.WithDefaultText("Enter description of the desired image").WithMultiLine().Show() // Text input with multi line enabled
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
				ErrorFmt.Printf("No Image Request to Send, exiting...\n")
				return nil
			}

			if err := sendSpeechMessages(speechFlags, chatContext, client, chatRequestString); err != nil {
				log.WithError(err).Fatal()
			}

			if !chatContext.InteractiveSession {
				break
			}
		}
		return nil
	}
}

func printSpeechBanner(f *SpeechFlags) {
	TitleFmt.Printf("ChatGPT CLI v%s\n", version)
	fmt.Printf("Model: %s, Voice: %s, speed: %0.2f\n", f.Model, f.Voice, f.Speed)
	fmt.Printf("- Press TAB after entering a message to send.\n")
	fmt.Printf("- Press TAB or CTRL+C with a blank message to terminate the session without sending.\n")
}

// sendMessages sends messages to ChatGPT and prints the response
func sendSpeechMessages(f *SpeechFlags, chatContext *ChatContext, client *openai.Client, chatRequestString string) error {
	mySpinner := pterm.DefaultSpinner
	mySpinner.Sequence = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
	mySpinner.RemoveWhenDone = true
	mySpinner.Writer = os.Stderr
	successSpinner, err := mySpinner.Start("Sending to ChatGPT TTS please wait...")

	imageRequest := openai.CreateSpeechRequest{
		Model:          f.Model,
		Input:          chatRequestString,
		ResponseFormat: openai.SpeechResponseFormatMp3,
		Voice:          f.Voice,
		Speed:          f.Speed,
	}
	resp, err := client.CreateSpeech(context.Background(), imageRequest)
	successSpinner.Success()

	if err != nil {
		return err
	}

	fileName := getSpeechFileName(f)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("File creation error: %v\n", err)
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp)
	if err != nil {
		fmt.Printf("File copy error: %v\n", err)
		return err
	}

	fmt.Printf("%s\n", fileName)
	if chatContext.InteractiveSession {
		os2.OpenBrowser(fileName)
	}

	return nil
}

func getSpeechFileName(f *SpeechFlags) string {
	thisImageCount := f.CurrentImageCount
	f.CurrentImageCount = thisImageCount + 1
	filename := fmt.Sprintf("%s-%02d.mp3", f.OutputPrefix, thisImageCount)
	return filename
}
