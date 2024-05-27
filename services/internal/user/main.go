package user

import (
	"log"
	"services/internal/user/controllers"
	"services/internal/user/integrations"
	"services/internal/user/repositories"
	"services/internal/user/router"
	"services/internal/user/services"
	"services/pkg/common/response"
	"services/pkg/database"
	"services/pkg/eventbus/rabbitmq"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"

	_ "services/internal/user/docs"
)

// @title User App Service
// @version 1.0
// @description This is a Documentation from User Service
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /api/user
// @securityDefinitions.apikey BearerAuth
// @in header
// @name X-Token
func Run() {
	brokerConfig := rabbitmq.RabbitMQConfig{PrefetchCount: 10, URL: "amqp://rabbitmq:5672", ExchangeName: "Replica.User", AcK: true, ExchangePrefix: "IntegrationEvent", ExchangeSuffix: "User"}

	databaseConfig := database.DatabaseConfig{Url: "mongodb://root:password@mongo:27017/?authSource=admin", DatabaseName: "test_user"}

	db, err := database.New(databaseConfig)

	if err != nil {
		panic(err)
	}
	defer db.Close()
	service_user, err := initializeUserService(db)
	if err != nil {
		panic(err)
	}
	service_account, err := initializeAccountService(db)
	if err != nil {
		panic(err)
	}
	service_session, err := initializeSessionService(db)
	if err != nil {
		panic(err)
	}

	service_socket := services.NewSocketService()

	broker, err := initializeBroker(brokerConfig, service_user)
	if err != nil {
		panic(err)
	}
	defer broker.Close()

	controller, err := controllers.NewUserController(service_user, service_session, service_account, service_socket)
	if err != nil {
		panic(err)
	}

	_, err = initializeApp(controller)
	if err != nil {
		panic(err)
	}

}

func initializeBroker(config rabbitmq.RabbitMQConfig, user *services.UserService) (*rabbitmq.RabbitMQ, error) {
	broker, err := rabbitmq.NewRabbitMQ(config)
	if err != nil {
		return nil, err

	}
	// Define Base Integrations with needed Services.
	handlers, err := integrations.NewIntegrations(broker, user)
	// Define Subscriptions
	if err != nil {
		return nil, err
	}
	_, err = broker.AddSubscription(integrations.GetUserIntegrationEvent{}, handlers.GetUserIntegrationEventHandler)
	if err != nil {
		return nil, err
	}
	// Defined

	return broker, nil
}

func initializeUserService(db *database.Database) (*services.UserService, error) {

	repository_user, err := repositories.NewUserRepository(db.Database, "users")
	if err != nil {
		return nil, err
	}

	service, err := services.NewUserService(*repository_user)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func initializeAccountService(db *database.Database) (*services.AccountService, error) {
	repository_account, err := repositories.NewAccountRepository(db.Database, "accounts")
	if err != nil {
		return nil, err
	}

	service, err := services.NewAccountService(*repository_account)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func initializeSessionService(db *database.Database) (*services.SessionService, error) {
	repository_session, err := repositories.NewSessionRepository(db.Database, "sessions")
	if err != nil {
		return nil, err
	}

	service, err := services.NewSessionService(repository_session)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func initializeApp(controller *controllers.UserController) (any, error) {
	app := fiber.New(fiber.Config{
		ErrorHandler: response.ResponseWithErrorMessage,
	})

	app.Use(cors.New())
	app.Use(logger.New(logger.Config{
		DisableColors: false,
	}))

	app.Get("/swagger/*", swagger.New(swagger.Config{
		ConfigURL: "http://localhost:3000/api/user/swagger/doc.json",
	}))

	app, err := router.Initialize(app, controller)
	if err != nil {
		log.Fatal(err)
	}

	app.Listen("0.0.0.0:5000")

	//sigCh := make(chan os.Signal, 1)
	//signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	return nil, nil
}
