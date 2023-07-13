package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewVersionCmd() *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "displays version information",
		Long:  "displays version information",
		Run:   versionCmdRunner,
	}
	return versionCmd
}

func versionCmdRunner(_ *cobra.Command, _ []string) {
	log.Debugf("versionCmd called")
	_, _ = NarratorFmt.Printf("ChatGPT CLI v%s\n", version)
}
