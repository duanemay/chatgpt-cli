package cmd

import (
	"github.com/fatih/color"
	"github.com/sashabaranov/go-openai"
)

var version = "0.0.0-dev"

const FlagConfigFile = "config"

var configFile string

const FlagApiKey = "api-key"

var apikey string

const FlagVerbose = "verbose"

var verbose bool

const defaultModel = openai.GPT4
const FlagModel = "model"

var model string

const defaultRole = openai.ChatMessageRoleUser
const FlagRole = "role"

var role string

const defaultTemperature = 1.0
const FlagTemperature = "temperature"

var temperature float32

const defaultMaxTokens = 0
const FlagMaxTokens = "max-tokens"

var maxTokens int

const defaultTopP = 1.0
const FlagTopP = "top-p"

var topP float32

const FlagSessionFile = "session-file"

var sessionFile string

const defaultEomMarker = "\\s"
const FlagEomMarker = "eom"

var eomMarker string

const defaultEosMarker = "\\q"
const FlagEosMarker = "eos"

var eosMarker string

var interactiveSession = false

var client *openai.Client
var chat openai.ChatCompletionRequest

var humanFmt = color.New(color.FgHiRed).Add(color.Bold)
var aiFmt = color.New(color.FgHiGreen).Add(color.Bold)
var narratorFmt = color.New(color.FgHiWhite).Add(color.Bold)
