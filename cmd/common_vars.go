package cmd

import (
	"github.com/fatih/color"
)

var (
	humanFmt    = color.New(color.FgHiRed).Add(color.Bold)
	aiFmt       = color.New(color.FgHiGreen).Add(color.Bold)
	narratorFmt = color.New(color.FgHiWhite).Add(color.Bold)
)

var version = "0.0.0-dev"
