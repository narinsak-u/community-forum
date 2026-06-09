package handlers

import (
	"errors"

	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/middleware"
	"community-forum/backend/internal/ports"
	"community-forum/backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

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

	slug := c.Params("slug")
	userID := h.SessionManager.GetUserID(c)

	comment, err := h.CommentService.Create(c.Context(), slug, req.Content, req.ParentID, userID)
	if err != nil {
		if errors.Is(err, usecase.ErrThreadNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Thread not found",
			})
		}
		if errors.Is(err, domain.ErrInvalidInput) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if errors.Is(err, usecase.ErrCommentNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Parent comment not found",
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

func (h *CommentHandler) DeleteCommentHandler(c *fiber.Ctx) error {
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
		if errors.Is(err, usecase.ErrCommentNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Comment not found",
			})
		}
		if errors.Is(err, usecase.ErrCommentForbidden) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You do not have permission to delete this comment",
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
