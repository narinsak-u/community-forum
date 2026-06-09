// Package main is the entry point for the community forum backend application.
// In Hexagonal Architecture, this file serves as the "Composition Root" or "Main" component,
// where all dependencies are initialized and wired together.
package main

import (
	"log"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	db_adapter "community-forum/backend/internal/adapters/db"
	"community-forum/backend/internal/config"
	"community-forum/backend/internal/handlers"
	"community-forum/backend/internal/middleware"
	"community-forum/backend/internal/usecase"
)

// - load configuration
// - initialize database
// - create Fiber app
// - configure middlewares
// - define routes
func main() {
	cfg := config.Load()
	db := config.InitDB(cfg)
	app := fiber.New(
		fiber.Config{
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			},
		},
	)

	// Step 6: Initialize JWT-based authentication.
	sessionManager := middleware.NewSessionManager(cfg.JWTSecret, cfg.JWTExpiry)

	// Step 7: Add global middlewares.
	// recover: Recovers from panics to prevent the server from crashing.
	// logger: Logs every incoming HTTP request to the console.
	// cors: Configures Cross-Origin Resource Sharing.
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSOrigins,
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
	commentHandler := handlers.NewCommentHandler(commentService, sessionManager)
	voteHandler := handlers.NewVoteHandler(voteService, sessionManager)
	userHandler := handlers.NewUserHandler(userService, threadService, sessionManager)
	tagHandler := handlers.NewTagHandler(tagService, sessionManager)

	// Step 11: Register API routes.
	// Auth routes — rate limited to prevent brute-force attacks
	authLimiter := limiter.New(limiter.Config{
		Max:        10,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		SkipSuccessfulRequests: true,
	})
	api.Post("/auth/signup", authLimiter, authHandler.SignupHandler)
	api.Post("/auth/signin", authLimiter, authHandler.SigninHandler)
	api.Post("/auth/signout", sessionManager.RequireAuth, authHandler.SignoutHandler)
	api.Get("/auth/me", sessionManager.RequireAuth, authHandler.MeHandler)

	// Thread routes
	api.Post("/threads", sessionManager.RequireAuth, threadHandler.CreateThreadHandler)
	api.Get("/threads", threadHandler.ListThreadsHandler)
	api.Get("/threads/featured", threadHandler.FeaturedThreadHandler)
	api.Get("/threads/trending", threadHandler.TrendingThreadsHandler)
	api.Get("/threads/:slug", threadHandler.GetThreadHandler)
	api.Patch("/threads/:slug", sessionManager.RequireAuth, threadHandler.UpdateThreadHandler)
	api.Delete("/threads/:slug", sessionManager.RequireAuth, threadHandler.DeleteThreadHandler)

	// Comment and Vote routes
	api.Post("/threads/:slug/comments", sessionManager.RequireAuth, commentHandler.CreateCommentHandler)
	api.Delete("/comments/:id", sessionManager.RequireAuth, commentHandler.DeleteCommentHandler)
	api.Post("/threads/:slug/vote", sessionManager.RequireAuth, voteHandler.VoteThreadHandler)
	api.Post("/comments/:id/vote", sessionManager.RequireAuth, voteHandler.VoteCommentHandler)

	// User and Tag routes
	api.Get("/users/:username", userHandler.GetUserHandler)
	api.Patch("/users/:username", sessionManager.RequireAuth, userHandler.UpdateUserHandler)
	api.Get("/users/:username/threads", userHandler.GetUserThreadsHandler)
	api.Get("/tags", tagHandler.ListTagsHandler)
	api.Post("/tags", sessionManager.RequireAuth, tagHandler.CreateTagHandler)

	// Step 12: Start the server on the configured port.
	log.Printf("Server starting on http://localhost:%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
