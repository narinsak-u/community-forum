package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"community-forum/backend/internal/models"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "community_forum"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_SSLMODE", "disable"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Thread{},
		&models.Comment{},
		&models.Tag{},
		&models.Vote{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	app := fiber.New(
		fiber.Config{
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			},
		},
	)

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Community Forum API"})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	api := app.Group("/api/v1")

	api.Get("/threads", func(c *fiber.Ctx) error {
		var threads []models.Thread
		if result := db.Preload("Author").Preload("Tags").Find(&threads); result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": result.Error.Error(),
			})
		}
		return c.JSON(threads)
	})

	api.Get("/threads/:slug", func(c *fiber.Ctx) error {
		slug := c.Params("slug")
		var thread models.Thread
		if result := db.Preload("Author").Preload("Tags").Preload("Comments.Author").First(&thread, "slug = ?", slug); result.Error != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Thread not found",
			})
		}
		return c.JSON(thread)
	})

	port := getEnv("PORT", "8080")
	log.Printf("Server starting on http://localhost:%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
