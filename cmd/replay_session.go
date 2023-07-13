package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type ReplaySessionFlags struct {
	sessionFile string
}

func NewReplaySessionFlags() *ReplaySessionFlags {
	return &ReplaySessionFlags{}
}

func NewReplaySessionCmd() *cobra.Command {
	f := NewReplaySessionFlags()
	var cmd = &cobra.Command{
		Use:   "replay-session",
		Short: "Replay a chat session from saved file",
		Long:  "Replay  a chat session from saved file",
		Run:   replaySessionCmdRun(f),
	}

	AddReplaySessionFileFlag(&f.sessionFile, cmd.Flags())
	_ = cmd.MarkFlagRequired(FlagSessionFile)

	return cmd
}

func replaySessionCmdRun(f *ReplaySessionFlags) func(cmd *cobra.Command, args []string) {
	return func(_ *cobra.Command, _ []string) {
		log.Debugf("replaySessionCmd called")

		chatCompletionRequest := loadSessionFile(f.sessionFile)
		for i, message := range chatCompletionRequest.Messages {
			if i%2 == 0 {
				_, _ = HumanFmt.Printf("\n%s:\n", message.Role)
			} else {
				_, _ = AiFmt.Printf("%s:\n", message.Role)
			}
			fmt.Printf("%s\n", message.Content)
		}
		fmt.Println()
	}
}
