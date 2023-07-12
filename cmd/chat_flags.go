package cmd

type ChatFlags struct {
	model                string
	role                 string
	sessionFile          string
	skipWriteSessionFile bool
	eomMarker            string
	eosMarker            string
	temperature          float32
	maxTokens            int
	topP                 float32
}

func NewChatFlags() *ChatFlags {
	return &ChatFlags{}
}
