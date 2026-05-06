package handlers

import (
	"community-forum/backend/internal/middleware"
	"community-forum/backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// VoteHandler handles upvoting and downvoting for threads and comments.
type VoteHandler struct {
	DB *gorm.DB
}

func NewVoteHandler(db *gorm.DB) *VoteHandler {
	return &VoteHandler{DB: db}
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

	// Validate allowed values.
	if req.Value != -1 && req.Value != 0 && req.Value != 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Value must be -1, 0, or 1",
		})
	}

	// Find the target thread.
	slug := c.Params("slug")
	var thread models.Thread
	if result := h.DB.Where("slug = ?", slug).First(&thread); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Thread not found",
		})
	}

	userID := middleware.GetUserID(c)

	// If value is 0, it means the user wants to retract their vote.
	if req.Value == 0 {
		h.DB.Where("user_id = ? AND thread_id = ?", userID, thread.ID).Delete(&models.Vote{})
	} else {
		// Check if a vote already exists from this user for this thread.
		var existingVote models.Vote
		result := h.DB.Where("user_id = ? AND thread_id = ?", userID, thread.ID).First(&existingVote)

		if result.RowsAffected > 0 {
			// Update the existing vote if they change from up to down (or vice versa).
			h.DB.Model(&existingVote).Update("value", req.Value)
		} else {
			// Create a brand new vote record.
			threadID := thread.ID
			vote := models.Vote{
				UserID:   userID,
				ThreadID: &threadID,
				Value:    req.Value,
			}
			h.DB.Create(&vote)
		}
	}

	// Refresh thread data to get the new vote counts.
	h.DB.Where("slug = ?", slug).First(&thread)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "Vote recorded",
		"upvotes":   thread.Upvotes(h.DB),
		"downvotes": thread.Downvotes(h.DB),
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

	if req.Value != -1 && req.Value != 0 && req.Value != 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Value must be -1, 0, or 1",
		})
	}

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

	// Handle removing a vote.
	if req.Value == 0 {
		h.DB.Where("user_id = ? AND comment_id = ?", userID, comment.ID).Delete(&models.Vote{})
	} else {
		// Handle creating or updating a vote.
		var existingVote models.Vote
		result := h.DB.Where("user_id = ? AND comment_id = ?", userID, comment.ID).First(&existingVote)

		if result.RowsAffected > 0 {
			h.DB.Model(&existingVote).Update("value", req.Value)
		} else {
			commentIDUint := comment.ID
			vote := models.Vote{
				UserID:    userID,
				CommentID: &commentIDUint,
				Value:     req.Value,
			}
			h.DB.Create(&vote)
		}
	}

	// Refresh comment to get updated counts.
	h.DB.First(&comment, comment.ID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "Vote recorded",
		"upvotes":   comment.Upvotes(h.DB),
		"downvotes": comment.Downvotes(h.DB),
	})
}
