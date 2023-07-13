package cmd

import (
	"github.com/fatih/color"
)

var (
	HumanFmt    = color.New(color.FgHiRed).Add(color.Bold)
	AiFmt       = color.New(color.FgHiGreen).Add(color.Bold)
	NarratorFmt = color.New(color.FgHiWhite).Add(color.Bold)
)

var version = "0.0.0-dev"
