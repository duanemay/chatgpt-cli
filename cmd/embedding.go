package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pterm/pterm"
	"github.com/sashabaranov/go-openai"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewEmbeddingCmd(rootFlags *RootFlags) *cobra.Command {
	embeddingFlags := NewEmbeddingFlags()
	chatContext := NewChatContext()
	var cmd = &cobra.Command{
		Use:     "embedding",
		Aliases: []string{"embed"},
		Short:   "Generate embeddings for input text",
		Long:    "Generate embeddings for input text using OpenAI's embedding models",
		RunE:    embeddingCmdRunner(rootFlags, embeddingFlags, chatContext),
	}
	setChatContext(cmd, chatContext)

	AddEmbeddingModelFlag(&embeddingFlags.ModelStr, cmd.PersistentFlags())
	AddDimensionsFlag(&embeddingFlags.Dimensions, cmd.PersistentFlags())
	_ = cmd.MarkPersistentFlagRequired("apikey")

	return cmd
}

func embeddingCmdRunner(rootFlags *RootFlags, embeddingFlags *EmbeddingFlags, chatContext *ChatContext) func(cmd *cobra.Command, args []string) error {
	return func(_ *cobra.Command, _ []string) error {
		log.Debugf("embeddingCmd called")
		err := embeddingFlags.ValidateFlags()
		if err != nil {
			log.WithError(err).Fatal()
		}

		chatContext.InteractiveSession = detectTerminal()
		if chatContext.InteractiveSession {
			printEmbeddingBanner(embeddingFlags)
		}
		client, err := setupOpenAIClient(rootFlags.apikey)
		if err != nil {
			log.WithError(err).Fatal()
		}

		reader := bufio.NewReader(os.Stdin)
		var inputText string
		for {
			if chatContext.InteractiveSession {
				inputText, _ = pterm.DefaultInteractiveTextInput.WithDefaultText("Enter text to generate embeddings for").WithMultiLine().Show()
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
				inputText = strings.Join(lines, "\n")
			}
			if len(inputText) == 0 {
				ErrorFmt.Printf("No text to embed, exiting...\n")
				return nil
			}

			if err := sendEmbeddingRequest(embeddingFlags, client, inputText); err != nil {
				log.WithError(err).Fatal()
			}

			if !chatContext.InteractiveSession {
				break
			}
		}
		return nil
	}
}

func printEmbeddingBanner(f *EmbeddingFlags) {
	TitleFmt.Printf("ChatGPT CLI v%s\n", version)
	if f.Dimensions > 0 {
		fmt.Printf("Model: %s, Dimensions: %d\n", f.Model, f.Dimensions)
	} else {
		fmt.Printf("Model: %s\n", f.Model)
	}
	fmt.Printf("- Press TAB after entering text to send.\n")
	fmt.Printf("- Press TAB or CTRL+C with a blank message to terminate the session without sending.\n")
}

// sendEmbeddingRequest sends text to the OpenAI embeddings API and prints the response
func sendEmbeddingRequest(f *EmbeddingFlags, client *openai.Client, inputText string) error {
	mySpinner := pterm.DefaultSpinner
	mySpinner.Sequence = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
	mySpinner.RemoveWhenDone = true
	mySpinner.Writer = os.Stderr
	successSpinner, _ := mySpinner.Start("Sending to OpenAI Embeddings API, please wait...")

	embeddingReq := openai.EmbeddingRequest{
		Input: inputText,
		Model: f.Model,
	}
	if f.Dimensions > 0 {
		embeddingReq.Dimensions = f.Dimensions
	}

	resp, err := client.CreateEmbeddings(context.Background(), embeddingReq)
	if err != nil {
		successSpinner.Fail(err.Error())
		return err
	}
	successSpinner.Success()

	output := EmbeddingOutput{
		Model:      string(resp.Model),
		Dimensions: len(resp.Data[0].Embedding),
		Data:       make([]EmbeddingData, len(resp.Data)),
		Usage: EmbeddingUsage{
			PromptTokens: resp.Usage.PromptTokens,
			TotalTokens:  resp.Usage.TotalTokens,
		},
	}
	for i, d := range resp.Data {
		output.Data[i] = EmbeddingData{
			Index:     d.Index,
			Embedding: d.Embedding,
		}
	}

	jsonOutput, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsonOutput))

	return nil
}

// EmbeddingOutput is the structured output for embedding results
type EmbeddingOutput struct {
	Model      string          `json:"model"`
	Dimensions int             `json:"dimensions"`
	Data       []EmbeddingData `json:"data"`
	Usage      EmbeddingUsage  `json:"usage"`
}

// EmbeddingData holds a single embedding vector
type EmbeddingData struct {
	Index     int       `json:"index"`
	Embedding []float32 `json:"embedding"`
}

// EmbeddingUsage holds token usage information
type EmbeddingUsage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}
