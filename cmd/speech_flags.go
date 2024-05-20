package cmd

import (
	"fmt"
	"github.com/sashabaranov/go-openai"
)

type SpeechFlags struct {
	ModelStr string
	Model    openai.SpeechModel

	VoiceStr string
	Voice    openai.SpeechVoice

	Speed             float64
	OutputPrefix      string
	CurrentImageCount int
}

func NewSpeechFlags() *SpeechFlags {
	return &SpeechFlags{
		CurrentImageCount: 1,
	}
}

func (f *SpeechFlags) ValidateFlags() error {
	if f.Speed < 0.25 || f.Speed > 4.0 {
		return fmt.Errorf("speed must be between 0.25 and 4.0, inclusive")
	}

	switch f.ModelStr {
	case string(openai.TTSModel1), string(openai.TTSModel1HD), string(openai.TTSModelCanary):
		// these are fine
		f.Model = openai.SpeechModel(f.ModelStr)

	default:
		return fmt.Errorf("model must be one of tts-1, tts-1-hd, or canary-tts")
	}

	switch f.VoiceStr {
	case string(openai.VoiceAlloy),
		string(openai.VoiceEcho),
		string(openai.VoiceFable),
		string(openai.VoiceOnyx),
		string(openai.VoiceNova),
		string(openai.VoiceShimmer):
		// these are fine
		f.Voice = openai.SpeechVoice(f.VoiceStr)

	default:
		return fmt.Errorf("voice must be one of alloy, echo, fable, onyx, nova, or shimmer")
	}
	return nil
}
