package handlers

import (
	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/middleware"
	"community-forum/backend/internal/models"
	"community-forum/backend/internal/ports"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// UserHandler handles user-specific operations like profile viewing and updates.
type UserHandler struct {
	// UserService follows the Hexagonal Architecture (Port).
	UserService ports.UserService
	// DB is used directly here for some methods (like GetUserThreadsHandler).
	// This indicates the handler is still in transition to a full Hexagonal structure.
	DB          *gorm.DB
}

// NewUserHandler initializes the handler with its dependencies.
func NewUserHandler(userService ports.UserService, db *gorm.DB) *UserHandler {
	return &UserHandler{UserService: userService, DB: db}
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
	userID := middleware.GetUserID(c)
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

	// Step 1: Find the user by username.
	var user models.User
	if result := h.DB.Where("username = ?", username).First(&user); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Step 2: Handle pagination queries (e.g., ?page=1&pageSize=5).
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "5"))
	if pageSize < 1 {
		pageSize = 5
	}
	if pageSize > 50 {
		pageSize = 50 // Cap page size for performance
	}

	// Step 3: Count total threads for pagination metadata.
	var total int64
	h.DB.Model(&models.Thread{}).Where("author_id = ?", user.ID).Count(&total)

	// Step 4: Fetch the specific "slice" of threads for the current page.
	offset := (page - 1) * pageSize
	var threads []models.Thread
	// Preload fetches related data (Author and Tags) in as few queries as possible.
	h.DB.Preload("Author").Preload("Tags").
		Where("author_id = ?", user.ID).
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&threads)

	// Step 5: Format the database models into a JSON-friendly structure.
	type threadItem struct {
		ID           uint        `json:"id"`
		Title        string      `json:"title"`
		Slug         string      `json:"slug"`
		Content      string      `json:"content"`
		Status       string      `json:"status"`
		ViewCount    uint        `json:"view_count"`
		Upvotes      int64       `json:"upvotes"`
		Downvotes    int64       `json:"downvotes"`
		RepliesCount int64       `json:"replies_count"`
		CreatedAt    string      `json:"created_at"`
		Author       fiber.Map   `json:"author"`
		Tags         []fiber.Map `json:"tags"`
	}

	items := make([]threadItem, 0, len(threads))
	for _, t := range threads {
		tags := make([]fiber.Map, 0, len(t.Tags))
		for _, tag := range t.Tags {
			tags = append(tags, fiber.Map{
				"id":    tag.ID,
				"name":  tag.Name,
				"color": tag.Color,
			})
		}

		items = append(items, threadItem{
			ID:           t.ID,
			Title:        t.Title,
			Slug:         t.Slug,
			Content:      t.Content,
			Status:       t.Status,
			ViewCount:    t.ViewCount,
			Upvotes:      t.Upvotes(h.DB),
			Downvotes:    t.Downvotes(h.DB),
			RepliesCount: t.RepliesCount(h.DB),
			CreatedAt:    t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			Author: fiber.Map{
				"id":       t.Author.ID,
				"username": t.Author.Username,
				"avatar":   t.Author.Avatar,
			},
			Tags: tags,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

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
