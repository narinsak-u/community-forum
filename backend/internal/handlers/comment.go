package handlers

import (
	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/middleware"
	"community-forum/backend/internal/ports"

	"github.com/gofiber/fiber/v2"
)

// CommentHandler manages creation and deletion of comments.
type CommentHandler struct {
	CommentService ports.CommentService
	SessionManager *middleware.SessionManager
}

func NewCommentHandler(commentService ports.CommentService, sessionManager *middleware.SessionManager) *CommentHandler {
	return &CommentHandler{
		CommentService: commentService,
		SessionManager: sessionManager,
	}
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

	slug := c.Params("slug")
	userID := h.SessionManager.GetUserID(c)

	comment, err := h.CommentService.Create(c.Context(), slug, req.Content, req.ParentID, userID)
	if err != nil {
		if err.Error() == "Thread not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Thread not found",
			})
		}
		if err.Error() == "Content must be between 2 and 10000 characters" ||
			err.Error() == "Parent comment not found" ||
			err.Error() == "Parent comment does not belong to this thread" ||
			err.Error() == "Replies can only be 1 level deep" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create comment",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Comment posted",
		"comment": mapCommentToResponse(comment),
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

	userID := h.SessionManager.GetUserID(c)
	userRole := h.SessionManager.GetUserRole(c)

	err = h.CommentService.Delete(c.Context(), uint(commentID), userID, userRole)
	if err != nil {
		if err.Error() == "Comment not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Comment not found",
			})
		}
		if err.Error() == "You do not have permission to delete this comment" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete comment",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Comment deleted",
	})
}

func mapCommentToResponse(c *domain.Comment) fiber.Map {
	return fiber.Map{
		"id":        c.ID,
		"content":   c.Content,
		"upvotes":   c.Upvotes,
		"downvotes": c.Downvotes,
		"author": fiber.Map{
			"id":       c.Author.ID,
			"username": c.Author.Username,
			"avatar":   c.Author.Avatar,
		},
		"replies":    []interface{}{},
		"created_at": c.CreatedAt,
	}
}
