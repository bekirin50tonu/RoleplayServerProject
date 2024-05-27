package router

import (
	"services/internal/storage/controllers"

	"github.com/gofiber/fiber/v2"
)

func NewRouter(controller *controllers.StorageController, app *fiber.App) error {

	app.Post("/upload/:disk", controller.SaveFileUpdatedDataHandler)

	app.Get("/image/:disk/:id", controller.GetFileFromParameterHandler)
	return nil
}
