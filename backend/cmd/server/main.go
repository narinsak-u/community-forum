// Package main is the entry point for the community forum backend application.
// In Hexagonal Architecture, this file serves as the "Composition Root" or "Main" component,
// where all dependencies are initialized and wired together.
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

	db_adapter "community-forum/backend/internal/adapters/db"
	"community-forum/backend/internal/handlers"
	"community-forum/backend/internal/middleware"
	"community-forum/backend/internal/models"
	"community-forum/backend/internal/usecase"
)

func main() {
	// Step 1: Load environment variables from a .env file if it exists.
	// godotenv.Load() returns an error if the file is missing, which we handle by logging.
	// This is a common Go idiom: check for error immediately after the function call.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Step 2: Construct the Database Source Name (DSN) using environment variables.
	// We use a helper function getEnv to provide default values if variables are missing.
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "community_forum"),
		getEnv("DB_PORT", "5433"),
		getEnv("DB_SSLMODE", "disable"),
	)

	// Step 3: Connect to the PostgreSQL database using GORM.
	// GORM is an ORM library that simplifies database interactions.
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// log.Fatalf logs the error and then terminates the program with exit code 1.
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Step 4: Perform database migrations.
	// AutoMigrate creates or updates database tables based on the provided GORM models.
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

	// Step 5: Initialize the Fiber web framework.
	// Fiber is a high-performance web framework inspired by Express (Node.js).
	app := fiber.New(
		fiber.Config{
			// Custom global error handler to return errors as JSON.
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			},
		},
	)

	// Step 6: Initialize the session store for authentication.
	sessionManager := middleware.NewSessionManager()

	// Step 7: Add global middlewares.
	// recover: Recovers from panics to prevent the server from crashing.
	// logger: Logs every incoming HTTP request to the console.
	// cors: Configures Cross-Origin Resource Sharing.
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:8080,http://127.0.0.1:8080",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
	}))

	// Step 8: Define basic health check routes.
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Community Forum API"})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Step 9: Group routes under /api/v1 for versioning.
	api := app.Group("/api/v1")

	// Step 10: Wire up the Hexagonal Architecture layers.
	// 1. Adapters (Database): Repository implementation.
	userRepo := db_adapter.NewGORMUserRepository(db)
	threadRepo := db_adapter.NewGORMThreadRepository(db)
	commentRepo := db_adapter.NewGORMCommentRepository(db)
	voteRepo := db_adapter.NewGORMVoteRepository(db)
	tagRepo := db_adapter.NewGORMTagRepository(db)

	// 2. Use Cases (Services): Business logic implementation, receiving the repository via Dependency Injection.
	authService := usecase.NewAuthService(userRepo)
	userService := usecase.NewUserService(userRepo)
	threadService := usecase.NewThreadService(threadRepo)
	commentService := usecase.NewCommentService(commentRepo, threadRepo)
	voteService := usecase.NewVoteService(voteRepo, threadRepo, commentRepo)
	tagService := usecase.NewTagService(tagRepo)

	// 3. Handlers (Controllers): HTTP layer, receiving services via Dependency Injection.
	authHandler := handlers.NewAuthHandler(authService, sessionManager)
	threadHandler := handlers.NewThreadHandler(threadService, sessionManager)

	// Step 11: Register API routes.
	// Auth routes
	api.Post("/auth/signup", authHandler.SignupHandler)
	api.Post("/auth/signin", authHandler.SigninHandler)
	api.Post("/auth/signout", sessionManager.RequireAuth, authHandler.SignoutHandler)
	api.Get("/auth/me", sessionManager.RequireAuth, authHandler.MeHandler)

	// Thread routes
	api.Post("/threads", sessionManager.RequireAuth, threadHandler.CreateThreadHandler)
	api.Get("/threads", threadHandler.ListThreadsHandler)
	api.Get("/threads/featured", threadHandler.FeaturedThreadHandler)
	api.Get("/threads/trending", threadHandler.TrendingThreadsHandler)
	api.Get("/threads/:slug", sessionManager.RequireAuth, threadHandler.GetThreadHandler)
	api.Patch("/threads/:slug", sessionManager.RequireAuth, threadHandler.UpdateThreadHandler)
	api.Delete("/threads/:slug", sessionManager.RequireAuth, threadHandler.DeleteThreadHandler)

	commentHandler := handlers.NewCommentHandler(commentService, sessionManager)
	voteHandler := handlers.NewVoteHandler(voteService, sessionManager)

	// Comment and Vote routes
	api.Post("/threads/:slug/comments", sessionManager.RequireAuth, commentHandler.CreateCommentHandler)
	api.Delete("/comments/:id", sessionManager.RequireAuth, commentHandler.DeleteCommentHandler)

	api.Post("/threads/:slug/vote", sessionManager.RequireAuth, voteHandler.VoteThreadHandler)
	api.Post("/comments/:id/vote", sessionManager.RequireAuth, voteHandler.VoteCommentHandler)

	userHandler := handlers.NewUserHandler(userService, threadService, sessionManager)
	tagHandler := handlers.NewTagHandler(tagService, sessionManager)

	// User and Tag routes
	api.Get("/users/:username", userHandler.GetUserHandler)
	api.Patch("/users/:username", sessionManager.RequireAuth, userHandler.UpdateUserHandler)
	api.Get("/users/:username/threads", userHandler.GetUserThreadsHandler)

	api.Get("/tags", tagHandler.ListTagsHandler)
	api.Post("/tags", sessionManager.RequireAuth, tagHandler.CreateTagHandler)

	// Step 12: Start the server on the configured port.
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on http://localhost:%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// getEnv is a helper function to read environment variables or return a default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
