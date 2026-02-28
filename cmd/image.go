package cmd

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"time"

	os2 "github.com/duanemay/chatgpt-cli/pkg/os"
	"github.com/sashabaranov/go-openai"
	"golang.org/x/image/webp"

	"io"
	"os"
	"strings"

	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewImageCmd(rootFlags *RootFlags) *cobra.Command {
	imageFlags := NewImageFlags()
	chatContext := NewChatContext()
	var cmd = &cobra.Command{
		Use:   "image",
		Short: "Create an image",
		Long:  "Create an image",
		RunE:  imageCmdRunner(rootFlags, imageFlags, chatContext),
	}
	setChatContext(cmd, chatContext)

	AddImageModelFlag(&imageFlags.Model, cmd.PersistentFlags())
	AddNumberImagesFlag(&imageFlags.NumberImages, cmd.PersistentFlags())
	AddImageQualityFlag(&imageFlags.Quality, cmd.PersistentFlags())
	AddImageSizeFlag(&imageFlags.Size, cmd.PersistentFlags())
	AddImageStyleFlag(&imageFlags.Style, cmd.PersistentFlags())
	AddOutputPrefixFlag(&imageFlags.OutputPrefix, "image-"+time.Now().UTC().Format(time.RFC3339), cmd.PersistentFlags())
	_ = cmd.MarkPersistentFlagRequired("apikey")

	return cmd
}

func imageCmdRunner(rootFlags *RootFlags, imageFlags *ImageFlags, chatContext *ChatContext) func(cmd *cobra.Command, args []string) error {
	return func(_ *cobra.Command, _ []string) error {
		log.Debugf("imageCmd called")
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
					if err != nil && !errors.Is(err, io.EOF) {
						log.WithError(err).Fatal()
					} else if errors.Is(err, io.EOF) {
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
	if f.Model == openai.CreateImageModelDallE2 {
		fmt.Printf("model: %s, numberImages: %d, size: %s\n", f.Model, f.NumberImages, f.Size)
	} else {
		fmt.Printf("model: %s, size: %s, style: %s, quality: %s\n", f.Model, f.Size, f.Style, f.Quality)
	}
	fmt.Printf("- Press TAB after entering a message to send.\n")
	fmt.Printf("- Press TAB or CTRL+C with a blank message to terminate the session without sending.\n")
}

// sendMessages sends messages to ChatGPT and prints the response
func sendImageMessages(f *ImageFlags, chatContext *ChatContext, client *openai.Client, chatRequestString string) error {
	mySpinner := pterm.DefaultSpinner
	mySpinner.Sequence = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
	mySpinner.RemoveWhenDone = true
	mySpinner.Writer = os.Stderr
	destination := "DALL-E"

	var imageRequest openai.ImageRequest
	if f.Model == openai.CreateImageModelGptImage1 {
		destination = "GPT-Image-1"
		imageRequest = openai.ImageRequest{
			Prompt:  chatRequestString,
			Model:   f.Model,
			N:       f.NumberImages,
			Quality: f.Quality,
			Size:    f.Size,
		}
	} else {
		imageRequest = openai.ImageRequest{
			Prompt:         chatRequestString,
			Model:          f.Model,
			N:              f.NumberImages,
			Quality:        f.Quality,
			ResponseFormat: openai.CreateImageResponseFormatB64JSON,
			Size:           f.Size,
			Style:          f.Style,
		}
	}
	successSpinner, _ := mySpinner.Start("Sending to " + destination + ", please wait...")
	resp, err := client.CreateImage(context.Background(), imageRequest)
	if err != nil {
		successSpinner.Fail(err.Error())
		return err
	}
	successSpinner.Success()

	for _, data := range resp.Data {
		if err := processImageData(data, f, chatContext); err != nil {
			fmt.Printf("%v\n", err)
			continue
		}
	}
	return nil
}

func getImageFileName(f *ImageFlags) string {
	thisImageCount := f.CurrentImageCount
	f.CurrentImageCount = thisImageCount + 1
	filename := fmt.Sprintf("%s-%02d.png", f.OutputPrefix, thisImageCount)
	return filename
}

func processImageData(data openai.ImageResponseDataInner, f *ImageFlags, chatContext *ChatContext) error {
	imgBytes, err := base64.StdEncoding.DecodeString(data.B64JSON)
	if err != nil {
		return fmt.Errorf("Base64 decode error: %w", err)
	}

	r := bytes.NewReader(imgBytes)
	contentType := http.DetectContentType(imgBytes)
	var imgData image.Image
	switch contentType {
	case "image/png":
		imgData, err = png.Decode(r)
		if err != nil {
			return fmt.Errorf("PNG decode error: %w", err)
		}
	case "image/jpeg":
		imgData, err = jpeg.Decode(r)
		if err != nil {
			return fmt.Errorf("JPEG decode error: %w", err)
		}
	case "image/webp":
		imgData, err = webp.Decode(r)
		if err != nil {
			return fmt.Errorf("WebP decode error: %w", err)
		}
	default:
		return fmt.Errorf("unsupported image content type: %s", contentType)
	}

	fileName := getImageFileName(f)
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("File creation error: %w", err)
	}
	defer func(file *os.File) { _ = file.Close() }(file)

	if err := png.Encode(file, imgData); err != nil {
		return fmt.Errorf("PNG encode error: %w", err)
	}
	fmt.Printf("%s\n", fileName)
	if chatContext.InteractiveSession {
		os2.OpenBrowser(fileName)
	}
	return nil
}
