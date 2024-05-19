package cmd

import (
	"fmt"
	"github.com/sashabaranov/go-openai"
)

type ImageFlags struct {
	Model             string
	Size              string
	Quality           string
	Style             string
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
	switch f.Model {
	case openai.CreateImageModelDallE2:
		return f.ValidateDalle2Flags()
	case openai.CreateImageModelDallE3:
		return f.ValidateDalle3Flags()
	default:
		return fmt.Errorf("model must be one of 'dall-e-2' or 'dall-e-3'")
	}
}

func (f *ImageFlags) ValidateDalle2Flags() error {
	if f.NumberImages < 1 || f.NumberImages > 10 {
		return fmt.Errorf("NumberImages must be between 1 and 10, inclusive")
	}
	switch f.Size {
	case openai.CreateImageSize256x256, openai.CreateImageSize512x512, openai.CreateImageSize1024x1024:
		// these are fine
	default:
		return fmt.Errorf("size must be one of 256x256, 512x512, or 1024x1024, for DALL-E 2")
	}
	return nil
}

func (f *ImageFlags) ValidateDalle3Flags() error {
	if f.NumberImages < 1 || f.NumberImages > 1 {
		return fmt.Errorf("NumberImages must be 1, for DALL-E 3")
	}
	switch f.Size {
	case openai.CreateImageSize1024x1024, openai.CreateImageSize1024x1792, openai.CreateImageSize1792x1024:
		// these are fine
	default:
		return fmt.Errorf("size must be one of 1024x1024, 1792x1024, or 1024x1792, for DALL-E 3")
	}
	switch f.Quality {
	case openai.CreateImageQualityStandard, openai.CreateImageQualityHD:
		// these are fine
	default:
		return fmt.Errorf("quality must be one of standard or hd, for DALL-E 3")
	}
	switch f.Style {
	case openai.CreateImageStyleNatural, openai.CreateImageStyleVivid:
		// these are fine
	default:
		return fmt.Errorf("style must be one of vivid or natural, for DALL-E 3")
	}
	return nil
}
