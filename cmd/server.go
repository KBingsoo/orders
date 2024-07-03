package cmd

import (
	"os"

	"github.com/KBingsoo/entities/pkg/models"
	"github.com/KBingsoo/orders/internal/domain/orders"
	"github.com/KBingsoo/orders/internal/gateways/database"
	"github.com/KBingsoo/orders/internal/gateways/pubsub"
	"github.com/KBingsoo/orders/internal/gateways/web"
	"github.com/joho/godotenv"
	"github.com/literalog/go-wise/wise"
	"github.com/streadway/amqp"

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

		repository, err := wise.NewMongoSimpleRepository[models.Order](col)
		if err != nil {
			return err
		}

		connection, err := amqp.Dial(os.Getenv("RABBIT_URI"))
		if err != nil {
			return err
		}
		defer connection.Close()

		cardConsumer, err := pubsub.NewCardConsumer(connection)
		if err != nil {
			return err
		}

		cardProducer, err := pubsub.NewCardProducer(connection)
		if err != nil {
			return err
		}

		service := orders.NewManager(repository, cardProducer, cardConsumer)

		handler := orders.NewHandler(service)

		server := web.NewServer(handler)

		errCh := make(chan error)

		go func() {
			errCh <- server.Run(8080)
		}()

		go func() {
			errCh <- service.Consume()
		}()

		return <-errCh

	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
