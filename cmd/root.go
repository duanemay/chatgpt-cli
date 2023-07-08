package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:               "chatgpt-cli",
	Short:             "chatgpt-cli is a CLI for ChatGPT",
	Long:              "chatgpt-cli is a CLI for ChatGPT",
	PersistentPreRunE: initializeConfig,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, FlagConfigFile, "c", "", "Config file (default ./.chatgpt-cli then $HOME/.chatgpt-cli)")
	rootCmd.PersistentFlags().StringVarP(&apikey, FlagApiKey, "k", "", "ChatGPT apiKey")
	rootCmd.PersistentFlags().BoolVarP(&verbose, FlagVerbose, "v", false, "verbose logging")

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	})
}

func initializeConfig(cmd *cobra.Command, _ []string) error {
	if configFile != "" {
		viper.SetConfigFile(configFile)
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

	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	return nil
}
