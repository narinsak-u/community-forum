package handlers

import (
	"community-forum/backend/internal/middleware"
	"community-forum/backend/internal/models"
	"encoding/json"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

type UpdateUserRequest struct {
	Avatar *string  `json:"avatar"`
	Bio    *string  `json:"bio"`
	Stacks []string `json:"stacks"`
}

func userResponse(db *gorm.DB, user models.User) fiber.Map {
	var stacks interface{}
	if user.Stacks == "" {
		stacks = nil
	} else {
		var s []string
		if err := json.Unmarshal([]byte(user.Stacks), &s); err != nil {
			stacks = nil
		} else {
			stacks = s
		}
	}

	return fiber.Map{
		"id":         user.ID,
		"username":   user.Username,
		"avatar":     user.Avatar,
		"bio":        user.Bio,
		"stacks":     stacks,
		"role":       user.Role,
		"created_at": user.CreatedAt,
	}
}

func (h *UserHandler) GetUserHandler(c *fiber.Ctx) error {
	username := c.Params("username")

	var user models.User
	if result := h.DB.Where("username = ?", username).First(&user); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"user": userResponse(h.DB, user),
	})
}

func (h *UserHandler) UpdateUserHandler(c *fiber.Ctx) error {
	username := c.Params("username")

	var user models.User
	if result := h.DB.Where("username = ?", username).First(&user); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	userID := middleware.GetUserID(c)
	if user.ID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden",
		})
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Bio != nil {
		if len(*req.Bio) > 500 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Bio must be at most 500 characters",
			})
		}
		user.Bio = *req.Bio
	}

	if req.Stacks != nil {
		if len(req.Stacks) > 10 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Stacks must have at most 10 items",
			})
		}
		marshaled, err := json.Marshal(req.Stacks)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to process stacks",
			})
		}
		user.Stacks = string(marshaled)
	}

	if req.Avatar != nil {
		user.Avatar = *req.Avatar
	}

	if result := h.DB.Save(&user); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return c.JSON(fiber.Map{
		"user": userResponse(h.DB, user),
	})
}

func (h *UserHandler) GetUserThreadsHandler(c *fiber.Ctx) error {
	username := c.Params("username")

	var user models.User
	if result := h.DB.Where("username = ?", username).First(&user); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "5"))
	if pageSize < 1 {
		pageSize = 5
	}
	if pageSize > 50 {
		pageSize = 50
	}

	var total int64
	h.DB.Model(&models.Thread{}).Where("author_id = ?", user.ID).Count(&total)

	offset := (page - 1) * pageSize
	var threads []models.Thread
	h.DB.Preload("Author").Preload("Tags").
		Where("author_id = ?", user.ID).
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&threads)

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
