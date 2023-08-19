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
