package handlers

import (
	"community-forum/backend/internal/middleware"
	"community-forum/backend/internal/ports"

	"github.com/gofiber/fiber/v2"
)

// VoteHandler handles upvoting and downvoting for threads and comments.
type VoteHandler struct {
	VoteService    ports.VoteService
	SessionManager *middleware.SessionManager
}

func NewVoteHandler(voteService ports.VoteService, sessionManager *middleware.SessionManager) *VoteHandler {
	return &VoteHandler{
		VoteService:    voteService,
		SessionManager: sessionManager,
	}
}

// VoteRequest expects a value: 1 (upvote), -1 (downvote), or 0 (remove vote).
type VoteRequest struct {
	Value int8 `json:"value"`
}

// VoteThreadHandler handles POST /api/threads/:slug/vote
func (h *VoteHandler) VoteThreadHandler(c *fiber.Ctx) error {
	var req VoteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	slug := c.Params("slug")
	userID := h.SessionManager.GetUserID(c)

	upvotes, downvotes, err := h.VoteService.VoteThread(c.Context(), slug, userID, req.Value)
	if err != nil {
		if err.Error() == "Value must be -1, 0, or 1" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err.Error() == "Thread not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Thread not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process vote",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "Vote recorded",
		"upvotes":   upvotes,
		"downvotes": downvotes,
	})
}

// VoteCommentHandler handles POST /api/comments/:id/vote
func (h *VoteHandler) VoteCommentHandler(c *fiber.Ctx) error {
	var req VoteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	commentID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid comment ID",
		})
	}

	userID := h.SessionManager.GetUserID(c)

	upvotes, downvotes, err := h.VoteService.VoteComment(c.Context(), uint(commentID), userID, req.Value)
	if err != nil {
		if err.Error() == "Value must be -1, 0, or 1" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err.Error() == "Comment not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Comment not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process vote",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "Vote recorded",
		"upvotes":   upvotes,
		"downvotes": downvotes,
	})
}
