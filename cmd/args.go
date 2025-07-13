package cmd

import (
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/pflag"
)

const (
	FlagApiKey               = "api-key"
	FlagConfigFile           = "config"
	FlagInitialSystemMessage = "system-message"
	FlagMaxTokens            = "max-tokens"
	FlagModel                = "model"
	FlagSkipWriteSessionFile = "skip-write-session"
	FlagRole                 = "role"
	FlagSessionFile          = "session-file"
	FlagTemperature          = "temperature"
	FlagTopP                 = "top-p"
	FlagVerbose              = "verbose"
	FlagImageSize            = "size"
	FlagNumberImages         = "number"
	FlagImageModel           = "model"
	FlagImageQuality         = "quality"
	FlagImageStyle           = "style"
	FlagOutputPrefix         = "output-prefix"
	FlagDetail               = "detail"
	FlagVoice                = "voice"
	FlagSpeed                = "speed"
	FlagInputFile            = "file"
	FlagLanguage             = "language"
)

const (
	defaultMaxTokens     = 0
	defaultModel         = openai.GPT4oMini
	defaultRole          = openai.ChatMessageRoleUser
	defaultTemperature   = 1.0
	defaultTopP          = 1.0
	defaultSystemMessage = ""
	defaultNumberImages  = 1
	defaultImageModel    = openai.CreateImageModelGptImage1
	defaultImageQuality  = openai.CreateImageQualityHigh
	defaultImageStyle    = openai.CreateImageStyleVivid
	defaultImageSize     = openai.CreateImageSize1024x1024
	defaultDetail        = string(openai.ImageURLDetailAuto)
	defaultSpeed         = 1.0
	defaultVoice         = string(openai.VoiceAlloy)
	defaultSpeechModel   = string(openai.TTSModel1)
)

// AddConfigFileFlag initialises the ConfigFile flag.
func AddConfigFileFlag(str *string, flags *pflag.FlagSet) {
	flags.StringVarP(str, FlagConfigFile, "c", "",
		"Config file (default ./.chatgpt-cli then $HOME/.chatgpt-cli)",
	)
}

func AddVerboseFlag(b *bool, flags *pflag.FlagSet) {
	flags.BoolVarP(b, FlagVerbose, "v", false, "verbose logging")
}

func AddApiKeyFlag(str *string, flags *pflag.FlagSet) {
	flags.StringVarP(str, FlagApiKey, "k", "", "ChatGPT apiKey")
}

func AddModelFlag(str *string, flags *pflag.FlagSet) {
	flags.StringVarP(str, FlagModel, "m", defaultModel, "ChatGPT Model")
}

func AddSpeechModelFlag(str *string, flags *pflag.FlagSet) {
	flags.StringVarP(str, FlagModel, "m", defaultSpeechModel, "ChatGPT Text to Speech Model")
}

func AddDetailFlag(str *string, flags *pflag.FlagSet) {
	flags.StringVarP(str, FlagDetail, "d", defaultDetail, "control how the model processes the image and generates its textual understanding. Must be one of low, high, or auto")
}

func AddVoiceFlag(str *string, flags *pflag.FlagSet) {
	flags.StringVar(str, FlagVoice, defaultVoice, "control the voice used. Must be one of alloy, echo, fable, onyx, nova, or shimmer")
}

func AddSpeedFlag(f *float64, flags *pflag.FlagSet) {
	flags.Float64VarP(f, FlagSpeed, "s", defaultSpeed, "control speed of audio. Must be between 0.25 and 4.0")
}

func AddInputFileFlag(str *[]string, flags *pflag.FlagSet) {
	flags.StringArrayVarP(str, FlagInputFile, "f", nil, "Input image files, maybe specified more than once")
}

func AddRoleFlag(str *string, flags *pflag.FlagSet) {
	flags.StringVarP(str, FlagRole, "r", defaultRole, "ChatGPT Role")
}

func AddSessionFileFlag(str *string, flags *pflag.FlagSet) {
	flags.StringVarP(str, FlagSessionFile, "s", "", "Continue a session from file")
}

func AddReplaySessionFileFlag(str *string, flags *pflag.FlagSet) {
	flags.StringVarP(str, FlagSessionFile, "s", "", "Replay a session from file")
}

func AddSkipWriteSessionFileFlag(b *bool, flags *pflag.FlagSet) {
	flags.BoolVar(b, FlagSkipWriteSessionFile, false, "Do not write or update session file")
}

func AddInitialSystemMessageFlag(str *string, flags *pflag.FlagSet) {
	flags.StringVar(str, FlagInitialSystemMessage, defaultSystemMessage, "Initial System message sent to ChatGPT")
}

func AddTemperatureFlag(f *float32, flags *pflag.FlagSet) {
	flags.Float32Var(f, FlagTemperature, defaultTemperature, "Temperature, between 0 and 2. Higher values make the output more random")
}

func AddMaxTokensFlag(i *int, flags *pflag.FlagSet) {
	flags.IntVar(i, FlagMaxTokens, defaultMaxTokens, "Maximum number of tokens to generate (default 0, no limit)")
}

func AddTopPFlag(f *float32, flags *pflag.FlagSet) {
	flags.Float32Var(f, FlagTopP, defaultTopP, "TopP, between 0 and 1. tokens with top_p probability mass")
}

func AddNumberImagesFlag(n *int, flags *pflag.FlagSet) {
	flags.IntVarP(n, FlagNumberImages, "n", defaultNumberImages, "Number of images to generate, between 1 and 10, for DALL-E 2")
}

func AddImageModelFlag(str *string, flags *pflag.FlagSet) {
	flags.StringVarP(str, FlagImageModel, "m", defaultImageModel, "Model to use for image generation. Must be one of 'gpt-image-1', 'dall-e-3', 'dall-e-2")
}

func AddImageStyleFlag(str *string, flags *pflag.FlagSet) {
	flags.StringVar(str, FlagImageStyle, defaultImageStyle, "Style to use for image generation. Must be one of 'vivid' or 'natural' (for DALL-E 3 only)")
}

func AddImageQualityFlag(str *string, flags *pflag.FlagSet) {
	flags.StringVar(str, FlagImageQuality, defaultImageQuality, "Quality to use for image generation. Must be one of 'standard' or 'hd' (for GPT-Image-1 and DALL-E 3)")
}

func AddImageSizeFlag(str *string, flags *pflag.FlagSet) {
	flags.StringVarP(str, FlagImageSize, "s", defaultImageSize, "Size of the generated images. Must be one of: 256x256, 512x512, or 1024x1024 for DALL-E 2; 1024x1024, 1792x1024, or 1024x1792 for DALL-E 3; or 1024x1024, 1024x1536, or 1536x1024 for GPT-Image-1")
}

func AddImageOutputPrefixFlag(str *string, defaultName string, flags *pflag.FlagSet) {
	flags.StringVarP(str, FlagOutputPrefix, "o", defaultName, "Prefix used for the output file names")
}

func AddLanguageFlag(str *string, flags *pflag.FlagSet) {
	flags.StringVarP(str, FlagLanguage, "l", "", "language of the input audio")
}
