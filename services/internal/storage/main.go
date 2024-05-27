package storage

import (
	"services/internal/storage/controllers"
	"services/internal/storage/integrations"
	"services/internal/storage/manager"
	"services/internal/storage/router"
	"services/internal/storage/services"
	"services/pkg/common/response"
	"services/pkg/database"
	"services/pkg/eventbus/rabbitmq"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Run() {
	brokerConfig := rabbitmq.RabbitMQConfig{PrefetchCount: 10, URL: "amqp://rabbitmq:5672", ExchangeName: "Replica.Storage", AcK: true, ExchangePrefix: "IntegrationEvent", ExchangeSuffix: "Storage"}

	// Storage Manager Inıtialized.
	manager, err := manager.NewStorageDiskSystem(manager.StorageManagerConfig{
		DbConfig:     database.DatabaseConfig{Url: "mongodb://root:password@mongo:27017/?authSource=admin", DatabaseName: "test_storage"},
		RootLocation: "./storage",
		Disks: map[string]*manager.StorageDisk{
			"user":  {Location: "user"},
			"guild": {Location: "guild"},
		}})
	if err != nil {
		panic(err)
	}
	// end

	// Initialize Infrastructures and Business Logic.
	storageService, err := services.NewStorageService(manager)

	_, err = initializeBroker(brokerConfig, storageService)

	controller := controllers.NewStorageController(storageService)

	// end

	app := fiber.New(fiber.Config{
		ErrorHandler: response.ResponseWithErrorMessage,
	})

	app.Use(cors.New())
	app.Use(logger.New())

	_ = router.NewRouter(controller, app)

	app.Listen(":5000")

}

func initializeBroker(config rabbitmq.RabbitMQConfig, storage *services.StorageService) (*rabbitmq.RabbitMQ, error) {
	broker, err := rabbitmq.NewRabbitMQ(config)
	if err != nil {
		return nil, err

	}
	// Define Base Integrations with needed Services.
	handlers, err := integrations.NewIntegrations(broker, storage)
	// Define Subscriptions
	if err != nil {
		return nil, err
	}
	_, err = broker.AddSubscription(integrations.UploadImageIntegration{}, handlers.UploadImageIntegrationEventHandler)
	if err != nil {
		return nil, err
	}
	// Defined

	return broker, nil
}
