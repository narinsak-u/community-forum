package handlers

import (
	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/middleware"
	"community-forum/backend/internal/ports"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// UserHandler handles user-specific operations like profile viewing and updates.
type UserHandler struct {
	UserService    ports.UserService
	ThreadService  ports.ThreadService
	SessionManager *middleware.SessionManager
}

// NewUserHandler initializes the handler with its dependencies.
func NewUserHandler(userService ports.UserService, threadService ports.ThreadService, sessionManager *middleware.SessionManager) *UserHandler {
	return &UserHandler{
		UserService:    userService,
		ThreadService:  threadService,
		SessionManager: sessionManager,
	}
}

// UpdateUserRequest uses pointers (e.g., *string) to distinguish between:
// - A field that is missing from the JSON (value is nil)
// - A field that is provided as an empty string (value is "")
type UpdateUserRequest struct {
	Avatar *string  `json:"avatar"`
	Bio    *string  `json:"bio"`
	Stacks []string `json:"stacks"`
}

// userResponse is a helper to format user data for JSON responses consistently.
func userResponse(user *domain.User) fiber.Map {
	return fiber.Map{
		"id":         user.ID,
		"username":   user.Username,
		"avatar":     user.Avatar,
		"bio":        user.Bio,
		"stacks":     user.Stacks,
		"role":       user.Role,
		"created_at": user.CreatedAt,
	}
}

// ListUsersHandler returns a list of all registered users.
func (h *UserHandler) ListUsersHandler(c *fiber.Ctx) error {
	users, err := h.UserService.ListUsers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	items := make([]fiber.Map, 0, len(users))
	for _, u := range users {
		items = append(items, userResponse(u))
	}

	return c.JSON(fiber.Map{
		"users": items,
	})
}

// GetUserHandler retrieves a user's profile by their username.
func (h *UserHandler) GetUserHandler(c *fiber.Ctx) error {
	// c.Params retrieves variables from the URL (e.g., /users/:username)
	username := c.Params("username")

	user, err := h.UserService.GetUserProfile(c.Context(), username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"user": userResponse(user),
	})
}

// UpdateUserHandler allows a user to change their own profile information.
func (h *UserHandler) UpdateUserHandler(c *fiber.Ctx) error {
	username := c.Params("username")

	// Verify the requester is logged in.
	userID := h.SessionManager.GetUserID(c)
	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Prepare an object with the requested updates.
	updates := &domain.User{}
	if req.Bio != nil {
		if len(*req.Bio) > 500 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Bio must be at most 500 characters",
			})
		}
		updates.Bio = *req.Bio
	}
	if req.Avatar != nil {
		updates.Avatar = *req.Avatar
	}
	if req.Stacks != nil {
		if len(req.Stacks) > 10 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Stacks must have at most 10 items",
			})
		}
		updates.Stacks = req.Stacks
	}

	// Security Check: Ensure the user is only updating their OWN profile.
	user, err := h.UserService.GetUserProfile(c.Context(), username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if user.ID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden",
		})
	}

	// Perform the update via the service layer.
	updatedUser, err := h.UserService.UpdateProfile(c.Context(), userID, updates)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return c.JSON(fiber.Map{
		"user": userResponse(updatedUser),
	})
}

// GetUserThreadsHandler returns a paginated list of threads created by a specific user.
func (h *UserHandler) GetUserThreadsHandler(c *fiber.Ctx) error {
	username := c.Params("username")

	// Verify user exists first to return 404
	_, err := h.UserService.GetUserProfile(c.Context(), username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "5"))

	threads, total, err := h.ThreadService.ListByUser(c.Context(), username, page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch threads",
		})
	}

	items := make([]fiber.Map, 0, len(threads))
	for _, t := range threads {
		tags := make([]fiber.Map, 0, len(t.Tags))
		for _, tag := range t.Tags {
			tags = append(tags, fiber.Map{
				"id":    tag.ID,
				"name":  tag.Name,
				"color": tag.Color,
			})
		}

		items = append(items, fiber.Map{
			"id":            t.ID,
			"title":         t.Title,
			"slug":          t.Slug,
			"content":       t.Content,
			"status":        t.Status,
			"view_count":    t.ViewCount,
			"upvotes":       t.Upvotes,
			"downvotes":     t.Downvotes,
			"replies_count": t.RepliesCount,
			"created_at":    t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			"author": fiber.Map{
				"id":       t.Author.ID,
				"username": t.Author.Username,
				"avatar":   t.Author.Avatar,
			},
			"tags": tags,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	if totalPages == 0 {
		totalPages = 1
	}

	return c.JSON(fiber.Map{
		"threads": items,
		"pagination": fiber.Map{
			"page":       page,
			"pageSize":   pageSize,
			"totalItems": total,
			"totalPages": totalPages,
		},
	})
}
