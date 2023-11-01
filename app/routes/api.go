package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/lab/tests/app/controllers"
)

func SetupApiRoutes(app *fiber.App, store *session.Store) {

	// * Single Image Upload
	singleImgUpload := app.Group("/api/v1")
	singleImgUpload.Post("/image-upload", controllers.SingleUploadImage)
	singleImgUpload.Get("/image-get", controllers.ListBlobsInContainer)

}
