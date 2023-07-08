package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "displays version information",
	Long:  "displays version information",
	Run:   versionCmdRunner,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func versionCmdRunner(cmd *cobra.Command, args []string) {
	log.Debugf("versionCmd called")
	narratorFmt.Printf("ChatGPT CLI v%s\n", version)
}
