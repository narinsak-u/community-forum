package handlers

import (
	"community-forum/backend/internal/middleware"
	"community-forum/backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CommentHandler manages creation and deletion of comments.
type CommentHandler struct {
	DB *gorm.DB
}

func NewCommentHandler(db *gorm.DB) *CommentHandler {
	return &CommentHandler{DB: db}
}

// CreateCommentRequest defines the input for a new comment or reply.
type CreateCommentRequest struct {
	Content  string `json:"content"`
	// ParentID is optional. If provided, the comment is a reply to another comment.
	ParentID *uint  `json:"parentId"`
}

// CreateCommentHandler handles POST /api/threads/:slug/comments
func (h *CommentHandler) CreateCommentHandler(c *fiber.Ctx) error {
	var req CreateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Basic validation for content length.
	if len(req.Content) < 2 || len(req.Content) > 10000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Content must be between 2 and 10000 characters",
		})
	}

	// Verify the thread exists via its slug.
	slug := c.Params("slug")
	var thread models.Thread
	if result := h.DB.Where("slug = ?", slug).First(&thread); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Thread not found",
		})
	}

	// If this is a reply (ParentID is set), perform additional checks.
	if req.ParentID != nil {
		var parent models.Comment
		if result := h.DB.First(&parent, *req.ParentID); result.Error != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Parent comment not found",
			})
		}

		// Ensure the parent comment actually belongs to the same thread.
		if parent.ThreadID != thread.ID {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Parent comment does not belong to this thread",
			})
		}

		// Prevent deeply nested comments (only allow 1 level of nesting).
		if parent.ParentID != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Replies can only be 1 level deep",
			})
		}
	}

	// Get logged-in user ID.
	userID := middleware.GetUserID(c)

	// Create the comment record.
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

	// Preload Author so the frontend knows who wrote the comment.
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

// DeleteCommentHandler handles DELETE /api/comments/:id
func (h *CommentHandler) DeleteCommentHandler(c *fiber.Ctx) error {
	// Parse the comment ID from the URL.
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

	// Security: Only the author or an admin can delete a comment.
	if comment.AuthorID != userID && userRole != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You do not have permission to delete this comment",
		})
	}

	// Clean up: If we delete a comment, we should also delete its replies.
	h.DB.Where("parent_id = ?", comment.ID).Delete(&models.Comment{})

	// Perform the deletion.
	h.DB.Delete(&comment)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Comment deleted",
	})
}
