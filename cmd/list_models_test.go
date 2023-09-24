package cmd_test

import (
	"github.com/duanemay/chatgpt-cli/cmd"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

var _ = Describe("List Models Command", func() {
	var rootCmd *cobra.Command
	commandName := "list-models"

	BeforeEach(func() {
		rootCmd = cmd.NewRootCmd()
	})

	It("should find command", func() {
		var thisCmd *cobra.Command
		Ω(rootCmd.Commands()).To(ContainElement(HaveField("Use", commandName), &thisCmd))
		Ω(thisCmd.Name()).To(Equal(commandName))
	})

	It("api-key should be required", func() {
		output, _ := ExecuteTest(rootCmd, []string{commandName, "-c", "test_files/empty.properties"}, "")
		Ω(output).To(ContainSubstring("OpenAI API Key not set\n"))
	})

	It("should list models", func() {
		output, _ := ExecuteTest(rootCmd, []string{commandName, "-v"}, "")
		Ω(output).To(ContainSubstring("gpt-4"))
		Ω(output).To(ContainSubstring("gpt-3.5-turbo"))
	})
})
