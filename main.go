package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var mongoRouterHost string

const mongoEnv = "GO_MONGODB_HOSTNAME"

func main() {
	mongoRouterHost = os.Getenv(mongoEnv)

	if mongoRouterHost == "" {
		log.Fatalf("Environment variable %s is not defined.\n", mongoEnv)
	}

	app := fiber.New()

	app.Get("/link-viewer/rest/links/:id", func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))

		if err != nil {
			return c.SendStatus(400)
		}

		record, err := getLinksRecord(id)

		if err != nil {
			// I know it's not a 404, but I don't want to leak information
			return c.SendStatus(404)
		}

		json, _ := json.Marshal(record.Links)

		return c.SendString(string(json))
	})

	app.Post("/link-viewer/rest/store", func(c *fiber.Ctx) error {
		var record LinksRecord

		if err := c.BodyParser(&record); err != nil {
			return c.SendStatus(400)
		}

		id, err := storeLinks(record.Links)

		if err != nil {
			return c.SendStatus(500)
		}

		return c.SendString(id)
	})

	log.Fatal(app.Listen(":3000"))
}
