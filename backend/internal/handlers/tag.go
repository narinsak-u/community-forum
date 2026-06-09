package handlers

import (
	"community-forum/backend/internal/middleware"
	"community-forum/backend/internal/ports"

	"github.com/gofiber/fiber/v2"
)

// TagHandler handles operations related to thread categories (tags).
type TagHandler struct {
	TagService     ports.TagService
	SessionManager *middleware.SessionManager
}

func NewTagHandler(tagService ports.TagService, sessionManager *middleware.SessionManager) *TagHandler {
	return &TagHandler{
		TagService:     tagService,
		SessionManager: sessionManager,
	}
}

// CreateTagRequest defines the data needed to create a new tag.
type CreateTagRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"` // Optional hex color code
}

// ListTagsHandler retrieves all available tags, sorted alphabetically.
func (h *TagHandler) ListTagsHandler(c *fiber.Ctx) error {
	tags, err := h.TagService.ListTags(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch tags",
		})
	}

	// Format the tags for the JSON response.
	items := make([]fiber.Map, 0, len(tags))
	for _, tag := range tags {
		items = append(items, fiber.Map{
			"id":    tag.ID,
			"name":  tag.Name,
			"color": tag.Color,
		})
	}

	return c.JSON(fiber.Map{
		"tags": items,
	})
}

// CreateTagHandler allows an admin to create a new tag.
func (h *TagHandler) CreateTagHandler(c *fiber.Ctx) error {
	var req CreateTagRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	userRole := h.SessionManager.GetUserRole(c)

	tag, err := h.TagService.CreateTag(c.Context(), req.Name, req.Color, userRole)
	if err != nil {
		if err.Error() == "Admin access required" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err.Error() == "Tag name must be between 3 and 50 characters" ||
			err.Error() == "Invalid color format, must be a valid hex color (e.g. #6366f1)" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err.Error() == "Tag with this name already exists" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create tag",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"tag": fiber.Map{
			"id":    tag.ID,
			"name":  tag.Name,
			"color": tag.Color,
		},
	})
}
