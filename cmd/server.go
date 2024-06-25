package cmd

import (
	"github.com/KBingsoo/entities/pkg/models"
	"github.com/KBingsoo/orders/internal/gateways/database"
	"github.com/joho/godotenv"
	"github.com/literalog/go-wise/wise"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start webserver",
	RunE: func(cmd *cobra.Command, args []string) error {

		err := godotenv.Load(".env")
		if err != nil {
			return err
		}

		col, err := database.GetCollection("orders")
		if err != nil {
			return err
		}

		_, err = wise.NewMongoSimpleRepository[models.Order](col)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
