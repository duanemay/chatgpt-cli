package cmd

import (
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type EmbeddingFlags struct {
	ModelStr   string
	Model      openai.EmbeddingModel
	Dimensions int
}

func NewEmbeddingFlags() *EmbeddingFlags {
	return &EmbeddingFlags{}
}

func (f *EmbeddingFlags) ValidateFlags() error {
	switch f.ModelStr {
	case string(openai.SmallEmbedding3), string(openai.LargeEmbedding3), string(openai.AdaEmbeddingV2):
		f.Model = openai.EmbeddingModel(f.ModelStr)
	default:
		return fmt.Errorf("model must be one of text-embedding-3-small, text-embedding-3-large, or text-embedding-ada-002")
	}

	if f.Dimensions < 0 {
		return fmt.Errorf("dimensions must be a non-negative integer")
	}

	return nil
}
