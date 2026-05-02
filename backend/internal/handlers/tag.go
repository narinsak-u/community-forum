package handlers

import (
	"community-forum/backend/internal/middleware"
	"community-forum/backend/internal/models"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TagHandler struct {
	DB *gorm.DB
}

func NewTagHandler(db *gorm.DB) *TagHandler {
	return &TagHandler{DB: db}
}

type CreateTagRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

var hexColorRegex = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

func (h *TagHandler) ListTagsHandler(c *fiber.Ctx) error {
	var tags []models.Tag
	if result := h.DB.Order("name ASC").Find(&tags); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch tags",
		})
	}

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

func (h *TagHandler) CreateTagHandler(c *fiber.Ctx) error {
	userRole := middleware.GetUserRole(c)
	if userRole != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin access required",
		})
	}

	var req CreateTagRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Color = strings.TrimSpace(req.Color)

	if len(req.Name) < 3 || len(req.Name) > 50 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tag name must be between 3 and 50 characters",
		})
	}

	if req.Color != "" && !hexColorRegex.MatchString(req.Color) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid color format, must be a valid hex color (e.g. #6366f1)",
		})
	}

	if req.Color == "" {
		req.Color = "#6366f1"
	}

	var existing models.Tag
	if h.DB.Where("LOWER(name) = ?", strings.ToLower(req.Name)).First(&existing).RowsAffected > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Tag with this name already exists",
		})
	}

	tag := models.Tag{
		Name:  req.Name,
		Color: req.Color,
	}

	if result := h.DB.Create(&tag); result.Error != nil {
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
