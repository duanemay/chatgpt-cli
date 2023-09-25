package cmd

import (
	"fmt"
	"github.com/sashabaranov/go-openai"
)

type ImageFlags struct {
	Size              string
	NumberImages      int
	OutputPrefix      string
	CurrentImageCount int
}

func NewImageFlags() *ImageFlags {
	return &ImageFlags{
		CurrentImageCount: 1,
	}
}

func (f *ImageFlags) ValidateFlags() error {
	if f.NumberImages < 1 || f.NumberImages > 10 {
		return fmt.Errorf("NumberImages must be between 1 and 10, inclusive")
	}
	switch f.Size {
	case openai.CreateImageSize256x256, openai.CreateImageSize512x512, openai.CreateImageSize1024x1024:
		// these are fine
	default:
		return fmt.Errorf("Size must be one of 256x256, 512x512, or 1024x1024")
	}
	return nil
}
