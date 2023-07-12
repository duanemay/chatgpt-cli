package cmd_test

import (
	"github.com/duanemay/chatgpt-cli/cmd"
	"github.com/spf13/cobra"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Version Command", func() {
	var rootCmd *cobra.Command
	commandName := "version"

	BeforeEach(func() {
		rootCmd = cmd.NewRootCmd()
	})

	It("should find command", func() {
		var thisCmd *cobra.Command
		Ω(rootCmd.Commands()).To(ContainElement(HaveField("Use", commandName), &thisCmd))
		Ω(thisCmd.Name()).To(Equal(commandName))
	})

	It("should output version", func() {
		output, _ := ExecuteTest(rootCmd, []string{commandName, "-v"}, "")
		Ω(output).To(ContainSubstring("versionCmd called"))
		Ω(output).To(ContainSubstring("ChatGPT CLI v0.0.0-dev"))
	})
})
