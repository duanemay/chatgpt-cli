package cmd

import (
	"github.com/pterm/pterm"
)

var (
	ErrorFmt = pterm.NewStyle(pterm.FgLightRed, pterm.Bold)
	HumanFmt = pterm.NewStyle(pterm.FgLightRed, pterm.Bold)
	AiFmt    = pterm.NewStyle(pterm.FgLightGreen, pterm.Bold)
	TitleFmt = pterm.NewStyle(pterm.FgLightWhite, pterm.Bold)
)

var version = "0.0.0-dev"
