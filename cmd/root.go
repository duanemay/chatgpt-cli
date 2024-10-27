package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
)

func NewRootCmd() *cobra.Command {
	rootFlags := NewRootFlags()
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	})

	cmds := &cobra.Command{
		Use:               "chatgpt-cli",
		Short:             "chatgpt-cli is a CLI for ChatGPT",
		Long:              "chatgpt-cli is a CLI for ChatGPT",
		PersistentPreRunE: initializeConfig(rootFlags),
		SilenceUsage:      true,
	}
	cmds.ResetFlags()

	cmds.AddCommand(NewImageCmd(rootFlags))
	cmds.AddCommand(NewChatCmd(rootFlags))
	cmds.AddCommand(NewVisionCmd(rootFlags))
	cmds.AddCommand(NewSpeechCmd(rootFlags))
	cmds.AddCommand(NewListModelsCmd(rootFlags))
	cmds.AddCommand(NewReplaySessionCmd())
	cmds.AddCommand(NewVersionCmd())
	cmds.AddCommand(NewTranscriptionCmd(rootFlags))

	AddConfigFileFlag(&rootFlags.configFile, cmds.PersistentFlags())
	AddApiKeyFlag(&rootFlags.apikey, cmds.PersistentFlags())
	AddVerboseFlag(&rootFlags.verbose, cmds.PersistentFlags())

	return cmds
}

func initializeConfig(rootFlags *RootFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		if rootFlags.configFile != "" {
			viper.SetConfigFile(rootFlags.configFile)
		} else {
			viper.SetConfigName(".chatgpt-cli")
			viper.AddConfigPath(".")
			viper.AddConfigPath("$HOME")
		}
		viper.SetConfigType("env")
		viper.SetEnvPrefix("CHATGPT")
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			// It's okay if there isn't a config file
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return err
			}
		}

		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			configName := f.Name
			configName = strings.ReplaceAll(configName, "-", "_")

			if !f.Changed && viper.IsSet(configName) {
				val := viper.Get(configName)
				_ = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			}
		})

		if rootFlags.verbose {
			log.SetLevel(log.DebugLevel)
		}

		return nil
	}
}
