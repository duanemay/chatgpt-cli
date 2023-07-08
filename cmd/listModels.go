package cmd

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var listModelsCmd = &cobra.Command{
	Use:   "list-models",
	Short: "lists all models available to your account",
	Long:  "lists all models available to your account",
	Run:   listModelsCmdRunner,
}

func init() {
	rootCmd.AddCommand(listModelsCmd)
	_ = chatCmd.MarkPersistentFlagRequired("apikey")
}

func listModelsCmdRunner(cmd *cobra.Command, args []string) {
	log.Debugf("listModelsCmd called")
	setupOpenAIClient()

	models, err := client.ListModels(context.Background())
	if err != nil {
		log.WithError(err).Fatal()
	}

	for _, model := range models.Models {
		fmt.Printf("%s\n", model.ID)
	}
}
