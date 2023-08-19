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
)

const (
	defaultMaxTokens     = 0
	defaultModel         = openai.GPT4
	defaultRole          = openai.ChatMessageRoleUser
	defaultTemperature   = 1.0
	defaultTopP          = 1.0
	defaultSystemMessage = ""
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
