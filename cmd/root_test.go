package cmd_test

import (
	"github.com/duanemay/chatgpt-cli/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Root Command", func() {
	var rootCmd *cobra.Command

	BeforeEach(func() {
		rootCmd = cmd.NewRootCmd()
		log.StandardLogger().SetLevel(log.InfoLevel)
	})

	It("should return a command", func() {
		Ω(rootCmd).ToNot(BeNil())
		Ω(rootCmd.Name()).To(Equal("chatgpt-cli"))
	})

	It("should have subcommands", func() {
		var chatCmd *cobra.Command
		Ω(rootCmd.Commands()).To(ContainElement(HaveField("Use", "version")))
		Ω(rootCmd.Commands()).To(ContainElement(HaveField("Use", "list-models")))
		Ω(rootCmd.Commands()).To(ContainElement(HaveField("Use", "replay-session")))
		Ω(rootCmd.Commands()).To(ContainElement(HaveField("Use", "chat"), &chatCmd))
		Ω(chatCmd.Name()).To(Equal("chat"))
	})

	It("should have default flag values", func() {
		Ω(rootCmd.PersistentFlags().GetString("config")).To(Equal(""))
		Ω(rootCmd.PersistentFlags().GetBool("verbose")).To(Equal(false))
		Ω(rootCmd.PersistentFlags().GetString("api-key")).To(Equal(""))
	})

	It("should display help", func() {
		output, _ := ExecuteTest(rootCmd, []string{""}, "")
		Ω(output).To(ContainSubstring("chatgpt-cli is a CLI for ChatGPT"))
		Ω(output).To(ContainSubstring("chatgpt-cli [command]"))
		Ω(output).To(ContainSubstring("chat"))
		Ω(output).To(ContainSubstring("list-models"))
		Ω(output).To(ContainSubstring("replay-session"))
		Ω(output).To(ContainSubstring("version"))
		Ω(output).To(ContainSubstring("Flags:"))
		Ω(output).To(ContainSubstring("--help"))
		Ω(output).To(ContainSubstring("--api-key"))
		Ω(output).To(ContainSubstring("--config"))
		Ω(output).To(ContainSubstring("--verbose "))
	})

	It("should set up default logger", func() {
		Ω(rootCmd.PersistentFlags().GetBool("verbose")).To(Equal(false))
		_, _ = ExecuteTest(rootCmd, []string{""}, "")
		StandardLogger := log.StandardLogger()
		Ω(StandardLogger.Level).To(Equal(log.InfoLevel))
		Ω(StandardLogger.Formatter).To(HaveField("FullTimestamp", true))
		Ω(StandardLogger.Formatter).To(HaveField("DisableLevelTruncation", true))
	})

	It("should set up verbose logger from CLI", func() {
		output, _ := ExecuteTest(rootCmd, []string{"version", "--verbose"}, "")
		Ω(rootCmd.PersistentFlags().GetBool("verbose")).To(Equal(true))
		StandardLogger := log.StandardLogger()
		Ω(StandardLogger.Level).To(Equal(log.DebugLevel))
		Ω(output).To(ContainSubstring("versionCmd called"))
	})

	It("should set up verbose logger from ENV", func() {
		_ = os.Setenv("CHATGPT_VERBOSE", "true")
		defer func() {
			_ = os.Unsetenv("CHATGPT_VERBOSE")
		}()

		output, _ := ExecuteTest(rootCmd, []string{"version"}, "")
		Ω(rootCmd.PersistentFlags().GetBool("verbose")).To(Equal(true))
		StandardLogger := log.StandardLogger()
		Ω(StandardLogger.Level).To(Equal(log.DebugLevel))
		Ω(output).To(ContainSubstring("versionCmd called"))
	})

	It("should have default flag values, after empty config file", func() {
		_, _ = ExecuteTest(rootCmd, []string{"version", "-c", "test_files/empty.properties"}, "")
		Ω(rootCmd.PersistentFlags().GetString("config")).To(Equal("test_files/empty.properties"))
		Ω(rootCmd.PersistentFlags().GetBool("verbose")).To(Equal(false))
		Ω(rootCmd.PersistentFlags().GetString("api-key")).To(Equal(""))
	})

	It("should have set flag values, from CONFIG", func() {
		_, _ = ExecuteTest(rootCmd, []string{"version", "-c", "test_files/rootFlags.properties"}, "")
		Ω(rootCmd.PersistentFlags().GetString("config")).To(Equal("test_files/rootFlags.properties"))
		Ω(rootCmd.PersistentFlags().GetBool("verbose")).To(Equal(true))
		Ω(rootCmd.PersistentFlags().GetString("api-key")).To(Equal("your_api_key"))
	})

	It("should have default flag values, no err, missing config file", func() {
		_, _ = ExecuteTest(rootCmd, []string{"version", "-c", "test_files/missing.properties"}, "")
		Ω(rootCmd.PersistentFlags().GetString("config")).To(Equal("test_files/missing.properties"))
		Ω(rootCmd.PersistentFlags().GetBool("verbose")).To(Equal(false))
		Ω(rootCmd.PersistentFlags().GetString("api-key")).To(Equal(""))
	})

	It("should set up verbose logger from CONFIG", func() {
		output, _ := ExecuteTest(rootCmd, []string{"version", "-c", "test_files/rootFlags.properties"}, "")
		Ω(rootCmd.PersistentFlags().GetBool("verbose")).To(Equal(true))
		StandardLogger := log.StandardLogger()
		Ω(StandardLogger.Level).To(Equal(log.DebugLevel))
		Ω(output).To(ContainSubstring("versionCmd called"))
	})

	Context("should honor priority order of config", func() {
		It("should prefer command line flags over others", func() {
			_ = os.Setenv("CHATGPT_API_KEY", "22")
			defer func() {
				_ = os.Unsetenv("CHATGPT_API_KEY")
			}()

			_, _ = ExecuteTest(rootCmd, []string{"version", "-k", "1", "-c", "test_files/orderFlags.properties"}, "")
			Ω(rootCmd.PersistentFlags().GetString("config")).To(Equal("test_files/orderFlags.properties"))
			Ω(rootCmd.PersistentFlags().GetString("api-key")).To(Equal("1"))
		})

		It("should prefer env vars over config file", func() {
			_ = os.Setenv("CHATGPT_API_KEY", "22")
			defer func() {
				_ = os.Unsetenv("CHATGPT_API_KEY")
			}()

			_, _ = ExecuteTest(rootCmd, []string{"version", "-c", "test_files/orderFlags.properties"}, "")
			Ω(rootCmd.PersistentFlags().GetString("config")).To(Equal("test_files/orderFlags.properties"))
			Ω(rootCmd.PersistentFlags().GetString("api-key")).To(Equal("22"))
		})
		It("should prefer config file over default", func() {
			_, _ = ExecuteTest(rootCmd, []string{"version", "-c", "test_files/orderFlags.properties"}, "")
			Ω(rootCmd.PersistentFlags().GetString("config")).To(Equal("test_files/orderFlags.properties"))
			Ω(rootCmd.PersistentFlags().GetString("api-key")).To(Equal("333"))
		})

	})
})
