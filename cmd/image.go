package cmd

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	os2 "github.com/duanemay/chatgpt-cli/pkg/os"
	"github.com/sashabaranov/go-openai"
	"image/png"
	"time"

	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
)

func NewImageCmd(rootFlags *RootFlags) *cobra.Command {
	imageFlags := NewImageFlags()
	chatContext := NewChatContext()
	var cmd = &cobra.Command{
		Use:   "image",
		Short: "Create an image using DALL·E",
		Long:  "Create an image using DALL·E",
		RunE:  imageCmdRunner(rootFlags, imageFlags, chatContext),
	}
	cmd.SetContext(context.WithValue(context.Background(), "chatContext", chatContext))

	AddNumberImagesFlag(&imageFlags.NumberImages, cmd.PersistentFlags())
	AddImageSizeFlag(&imageFlags.Size, cmd.PersistentFlags())
	AddImageOutputPrefixFlag(&imageFlags.OutputPrefix, "dall-e-"+time.Now().UTC().Format(time.RFC3339), cmd.PersistentFlags())
	_ = cmd.MarkPersistentFlagRequired("apikey")

	return cmd
}

func imageCmdRunner(rootFlags *RootFlags, imageFlags *ImageFlags, chatContext *ChatContext) func(cmd *cobra.Command, args []string) error {
	return func(_ *cobra.Command, _ []string) error {
		log.Debugf("chatCmd called")
		err := imageFlags.ValidateFlags()
		if err != nil {
			log.WithError(err).Fatal()
		}

		chatContext.InteractiveSession = detectTerminal()
		if chatContext.InteractiveSession {
			printImageBanner(imageFlags)
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

			if err := sendImageMessages(imageFlags, chatContext, client, chatRequestString); err != nil {
				log.WithError(err).Fatal()
			}

			if !chatContext.InteractiveSession {
				break
			}
		}
		return nil
	}
}

func printImageBanner(f *ImageFlags) {
	TitleFmt.Printf("ChatGPT CLI v%s\n", version)
	fmt.Printf("numberImages: %d, size: %s\n", f.NumberImages, f.Size)
	fmt.Printf("- Press TAB after entering a message to send.\n")
	fmt.Printf("- Press TAB or CTRL+C with a blank message to terminate the session without sending.\n")
}

// sendMessages sends messages to ChatGPT and prints the response
func sendImageMessages(f *ImageFlags, chatContext *ChatContext, client *openai.Client, chatRequestString string) error {
	mySpinner := pterm.DefaultSpinner
	mySpinner.Sequence = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
	mySpinner.RemoveWhenDone = true
	mySpinner.Writer = os.Stderr
	successSpinner, err := mySpinner.Start("Sending to DALL-E, please wait...")

	imageRequest := openai.ImageRequest{
		Prompt:         chatRequestString,
		N:              f.NumberImages,
		Size:           f.Size,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
	}
	resp, err := client.CreateImage(context.Background(), imageRequest)
	successSpinner.Success()

	if err != nil {
		return err
	}

	for _, data := range resp.Data {
		imgBytes, err := base64.StdEncoding.DecodeString(data.B64JSON)
		if err != nil {
			fmt.Printf("Base64 decode error: %v\n", err)
			continue
		}

		r := bytes.NewReader(imgBytes)
		imgData, err := png.Decode(r)
		if err != nil {
			fmt.Printf("PNG decode error: %v\n", err)
			continue
		}

		fileName := getFileName(f)
		file, err := os.Create(fileName)
		if err != nil {
			fmt.Printf("File creation error: %v\n", err)
			continue
		}
		defer file.Close()

		if err := png.Encode(file, imgData); err != nil {
			fmt.Printf("PNG encode error: %v\n", err)
			continue
		}
		fmt.Printf("%s\n", fileName)
		if chatContext.InteractiveSession {
			os2.OpenBrowser(fileName)
		}
	}
	return nil
}

func getFileName(f *ImageFlags) string {
	thisImageCount := f.CurrentImageCount
	f.CurrentImageCount = thisImageCount + 1
	filename := fmt.Sprintf("%s-%02d.png", f.OutputPrefix, thisImageCount)
	return filename
}
