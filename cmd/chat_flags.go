package cmd

type ChatFlags struct {
	model                string
	role                 string
	initialSystemMessage string
	sessionFile          string
	skipWriteSessionFile bool
	temperature          float32
	maxTokens            int
	topP                 float32
}

func NewChatFlags() *ChatFlags {
	return &ChatFlags{}
}

func ChatFlagsFromVisionFlags(f *VisionFlags) *ChatFlags {
	return &ChatFlags{
		model:                f.model,
		role:                 f.role,
		initialSystemMessage: f.initialSystemMessage,
		sessionFile:          f.sessionFile,
		skipWriteSessionFile: f.skipWriteSessionFile,

		temperature: defaultTemperature,
		maxTokens:   defaultMaxTokens,
		topP:        defaultTopP,
	}
}
