package cmd_test

import (
	"bytes"
	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"os"
	"sync"
)

func ExecuteTest(cmd *cobra.Command, args []string, input string) (output string, err error) {
	// Store the original stdout and stderr
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
	}()

	// Create a pipe to capture stdout and stderr, redirect everything to the pipe
	reader, writer, err := os.Pipe()
	os.Stdout = writer
	os.Stderr = writer
	pterm.SetDefaultOutput(writer)
	cmd.SetOut(writer)
	cmd.SetErr(writer)
	log.SetOutput(writer)

	// copy pipe output to a buffer, so we can read it later
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		_, _ = io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()

	// handle stdin
	in := bytes.NewBufferString(input)
	cmd.SetIn(in)

	// execute the command
	cmd.SetArgs(args)
	err = cmd.Execute()
	_ = writer.Close()

	// read the output
	output = <-out
	return
}
