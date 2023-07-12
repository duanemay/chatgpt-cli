package cmd

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewListModelsCmd(rootFlags *RootFlags) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list-models",
		Short: "lists all models available to your account",
		Long:  "lists all models available to your account",
		RunE:  listModelsCmdRunner(rootFlags),
	}

	_ = cmd.MarkPersistentFlagRequired("apikey")

	return cmd
}

func listModelsCmdRunner(rootFlags *RootFlags) func(cmd *cobra.Command, args []string) error {
	return func(_ *cobra.Command, _ []string) error {
		log.Debugf("listModelsCmd called")
		client, err := setupOpenAIClient(rootFlags.apikey)
		if err != nil {
			log.WithError(err).Error()
			return err
		}
		models, err := client.ListModels(context.Background())
		if err != nil {
			log.WithError(err).Error()
			return err
		}

		for _, model := range models.Models {
			fmt.Printf("%s\n", model.ID)
		}
		return nil
	}
}
