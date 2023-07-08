package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var replaySessionCmd = &cobra.Command{
	Use:   "replay-session",
	Short: "Replay a chat session from saved file",
	Long:  "Replay  a chat session from saved file",
	Run:   replaySessionCmdRunner,
}

func init() {
	rootCmd.AddCommand(replaySessionCmd)
	replaySessionCmd.Flags().StringVarP(&sessionFile, FlagSessionFile, "s", "", "File to replay a Session from")
	_ = replaySessionCmd.MarkFlagRequired(FlagSessionFile)
}

func replaySessionCmdRunner(cmd *cobra.Command, args []string) {
	log.Debugf("replaySessionCmd called")

	loadSessionFile()
	for i, message := range chat.Messages {
		if i%2 == 0 {
			humanFmt.Printf("\n%s:\n", message.Role)
		} else {
			aiFmt.Printf("%s:\n", message.Role)
		}
		fmt.Printf("%s\n", message.Content)
	}
	fmt.Println()
}
