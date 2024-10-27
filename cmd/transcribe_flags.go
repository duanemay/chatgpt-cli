package cmd

type TranscriptionFlags struct {
	model                string
	language             string
	initialSystemMessage string
	inputFiles           []string
}

func NewTranscriptionFlags() *TranscriptionFlags {
	return &TranscriptionFlags{model: "whisper-1"}
}

func (f *TranscriptionFlags) ValidateFlags() error {
	return nil
}
