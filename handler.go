package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/joho/godotenv"

	"github.com/lab/tests/app/routes"
)

func initLoadEnv() {

	err := godotenv.Load(".message.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	store := session.New()

	app := fiber.New()
	routes.SetupApiRoutes(app, store)

	initLoadEnv()

	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}

	err := app.Listen(listenAddr)
	if err != nil {
		log.Fatalf("Error while starting Fiber: %v", err)
	}

	log.Printf("Fiber is listening on %s. Go to https://127.0.0.1%s/", listenAddr, listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))

}
