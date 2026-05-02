package handlers

import (
	"community-forum/backend/internal/middleware"
	"community-forum/backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CommentHandler struct {
	DB *gorm.DB
}

func NewCommentHandler(db *gorm.DB) *CommentHandler {
	return &CommentHandler{DB: db}
}

type CreateCommentRequest struct {
	Content  string `json:"content"`
	ParentID *uint  `json:"parentId"`
}

func (h *CommentHandler) CreateCommentHandler(c *fiber.Ctx) error {
	var req CreateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if len(req.Content) < 2 || len(req.Content) > 10000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Content must be between 2 and 10000 characters",
		})
	}

	slug := c.Params("slug")
	var thread models.Thread
	if result := h.DB.Where("slug = ?", slug).First(&thread); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Thread not found",
		})
	}

	if req.ParentID != nil {
		var parent models.Comment
		if result := h.DB.First(&parent, *req.ParentID); result.Error != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Parent comment not found",
			})
		}

		if parent.ThreadID != thread.ID {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Parent comment does not belong to this thread",
			})
		}

		if parent.ParentID != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Replies can only be 1 level deep",
			})
		}
	}

	userID := middleware.GetUserID(c)

	comment := models.Comment{
		Content:  req.Content,
		ThreadID: thread.ID,
		AuthorID: userID,
		ParentID: req.ParentID,
	}

	if result := h.DB.Create(&comment); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create comment",
		})
	}

	h.DB.Preload("Author").First(&comment, comment.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Comment posted",
		"comment": fiber.Map{
			"id":        comment.ID,
			"content":   comment.Content,
			"upvotes":   comment.Upvotes(h.DB),
			"downvotes": comment.Downvotes(h.DB),
			"author": fiber.Map{
				"id":       comment.Author.ID,
				"username": comment.Author.Username,
				"avatar":   comment.Author.Avatar,
			},
			"replies":    []interface{}{},
			"created_at": comment.CreatedAt,
		},
	})
}

func (h *CommentHandler) DeleteCommentHandler(c *fiber.Ctx) error {
	commentID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid comment ID",
		})
	}

	var comment models.Comment
	if result := h.DB.First(&comment, commentID); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Comment not found",
		})
	}

	userID := middleware.GetUserID(c)
	userRole := middleware.GetUserRole(c)

	if comment.AuthorID != userID && userRole != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You do not have permission to delete this comment",
		})
	}

	h.DB.Where("parent_id = ?", comment.ID).Delete(&models.Comment{})

	h.DB.Delete(&comment)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Comment deleted",
	})
}
