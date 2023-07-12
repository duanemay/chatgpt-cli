package cmd

type ChatContext struct {
	InteractiveSession bool
}

func NewChatContext() *ChatContext {
	return &ChatContext{InteractiveSession: false}
}
