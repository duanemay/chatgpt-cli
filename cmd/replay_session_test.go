package cmd_test

import (
	"github.com/duanemay/chatgpt-cli/cmd"
	"github.com/spf13/cobra"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Replay Session Command", func() {
	var rootCmd *cobra.Command
	commandName := "replay-session"

	BeforeEach(func() {
		rootCmd = cmd.NewRootCmd()
	})

	It("should find command", func() {
		var thisCmd *cobra.Command
		Ω(rootCmd.Commands()).To(ContainElement(HaveField("Use", commandName), &thisCmd))
		Ω(thisCmd.Name()).To(Equal(commandName))
	})

	It("session-file should be required", func() {
		output, _ := ExecuteTest(rootCmd, []string{commandName}, "")
		Ω(output).To(ContainSubstring("Error: required flag(s) \"session-file\" not set\n"))
	})

	It("should output session", func() {
		output, _ := ExecuteTest(rootCmd, []string{commandName, "-v", "--session-file", "test_files/hello.json"}, "")
		Ω(output).To(ContainSubstring("replaySessionCmd called"))
		Ω(output).To(ContainSubstring("user:"))
		Ω(output).To(ContainSubstring("say hello in Japanese\n"))
		Ω(output).To(ContainSubstring("assistant:"))
		Ω(output).To(ContainSubstring("こんにちは\n"))
	})
})
