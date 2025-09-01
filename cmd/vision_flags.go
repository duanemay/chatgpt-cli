package cmd

import (
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type VisionFlags struct {
	model                string
	DetailStr            string
	Detail               openai.ImageURLDetail
	role                 string
	initialSystemMessage string
	inputFiles           []string

	skipWriteSessionFile bool
	sessionFile          string
}

func NewVisionFlags() *VisionFlags {
	return &VisionFlags{}
}

func (f *VisionFlags) ValidateFlags() error {
	switch f.DetailStr {
	case string(openai.ImageURLDetailAuto), string(openai.ImageURLDetailHigh), string(openai.ImageURLDetailLow):
		// these are fine, convert to ImageURLDetail
		f.Detail = openai.ImageURLDetail(f.DetailStr)
	default:
		return fmt.Errorf("detail must be one of 'auto', 'high', or 'low'")
	}
	return nil
}
