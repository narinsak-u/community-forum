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

	"community-forum/backend/internal/handlers"
	"community-forum/backend/internal/middleware"
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
		&models.Session{},
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

	middleware.InitSessionStore()

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

	authHandler := handlers.NewAuthHandler(db)
	threadHandler := handlers.NewThreadHandler(db)

	api.Post("/auth/signup", authHandler.SignupHandler)
	api.Post("/auth/signin", authHandler.SigninHandler)
	api.Post("/auth/signout", middleware.RequireAuth, authHandler.SignoutHandler)
	api.Get("/auth/me", middleware.RequireAuth, authHandler.MeHandler)

	api.Post("/threads", middleware.RequireAuth, threadHandler.CreateThreadHandler)
	api.Get("/threads", threadHandler.ListThreadsHandler)
	api.Get("/threads/featured", threadHandler.FeaturedThreadHandler)
	api.Get("/threads/trending", threadHandler.TrendingThreadsHandler)
	api.Get("/threads/:slug", middleware.RequireAuth, threadHandler.GetThreadHandler)
	api.Patch("/threads/:slug", middleware.RequireAuth, threadHandler.UpdateThreadHandler)
	api.Delete("/threads/:slug", middleware.RequireAuth, threadHandler.DeleteThreadHandler)

	commentHandler := handlers.NewCommentHandler(db)
	voteHandler := handlers.NewVoteHandler(db)

	api.Post("/threads/:slug/comments", middleware.RequireAuth, commentHandler.CreateCommentHandler)
	api.Delete("/comments/:id", middleware.RequireAuth, commentHandler.DeleteCommentHandler)

	api.Post("/threads/:slug/vote", middleware.RequireAuth, voteHandler.VoteThreadHandler)
	api.Post("/comments/:id/vote", middleware.RequireAuth, voteHandler.VoteCommentHandler)

	userHandler := handlers.NewUserHandler(db)
	tagHandler := handlers.NewTagHandler(db)

	api.Get("/users/:username", userHandler.GetUserHandler)
	api.Patch("/users/:username", middleware.RequireAuth, userHandler.UpdateUserHandler)
	api.Get("/users/:username/threads", userHandler.GetUserThreadsHandler)

	api.Get("/tags", tagHandler.ListTagsHandler)
	api.Post("/tags", middleware.RequireAuth, tagHandler.CreateTagHandler)

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
