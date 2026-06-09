package handlers

import (
	"errors"
	"math"
	"strconv"

	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/middleware"
	"community-forum/backend/internal/ports"
	"community-forum/backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type ThreadHandler struct {
	ThreadService  ports.ThreadService
	SessionManager *middleware.SessionManager
}

func NewThreadHandler(threadService ports.ThreadService, sessionManager *middleware.SessionManager) *ThreadHandler {
	return &ThreadHandler{
		ThreadService:  threadService,
		SessionManager: sessionManager,
	}
}

type CreateThreadRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
	Status  string   `json:"status"`
}

func (h *ThreadHandler) CreateThreadHandler(c *fiber.Ctx) error {
	var req CreateThreadRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	userID := h.SessionManager.GetUserID(c)

	thread, err := h.ThreadService.Create(c.Context(), req.Title, req.Content, req.Status, req.Tags, userID)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create thread",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(mapThreadToResponse(thread))
}

func (h *ThreadHandler) ListThreadsHandler(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "5"))
	sort := c.Query("sort", "latest")

	threads, total, err := h.ThreadService.List(c.Context(), page, pageSize, sort)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch threads",
		})
	}

	results := make([]fiber.Map, len(threads))
	for i, t := range threads {
		results[i] = mapThreadToResponse(&t)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"threads": results,
		"pagination": fiber.Map{
			"page":       page,
			"pageSize":   pageSize,
			"total":      total,
			"totalPages": totalPages,
		},
	})
}

func (h *ThreadHandler) FeaturedThreadHandler(c *fiber.Ctx) error {
	thread, err := h.ThreadService.GetFeatured(c.Context())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No featured thread found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(mapThreadToResponse(thread))
}

func (h *ThreadHandler) TrendingThreadsHandler(c *fiber.Ctx) error {
	threads, err := h.ThreadService.GetTrending(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch trending threads",
		})
	}

	results := make([]fiber.Map, len(threads))
	for i, t := range threads {
		results[i] = mapThreadToResponse(&t)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"threads": results,
	})
}

func (h *ThreadHandler) GetThreadHandler(c *fiber.Ctx) error {
	slug := c.Params("slug")

	thread, err := h.ThreadService.GetBySlug(c.Context(), slug)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Thread not found",
		})
	}

	if thread.Status != "published" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Thread not found",
		})
	}

	resp := mapThreadToResponse(thread)
	resp["comments"] = serializeComments(thread.Comments)

	return c.Status(fiber.StatusOK).JSON(resp)
}

type UpdateThreadRequest struct {
	Title   *string  `json:"title"`
	Content *string  `json:"content"`
	Status  *string  `json:"status"`
	Tags    []string `json:"tags"`
}

func (h *ThreadHandler) UpdateThreadHandler(c *fiber.Ctx) error {
	slug := c.Params("slug")
	userID := h.SessionManager.GetUserID(c)

	var req UpdateThreadRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	thread, err := h.ThreadService.Update(c.Context(), slug, userID, req.Title, req.Content, req.Status, req.Tags)
	if err != nil {
		if errors.Is(err, usecase.ErrThreadNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Thread not found",
			})
		}
		if errors.Is(err, usecase.ErrPermissionDenied) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You do not have permission to update this thread",
			})
		}
		if errors.Is(err, domain.ErrInvalidInput) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update thread",
		})
	}

	return c.Status(fiber.StatusOK).JSON(mapThreadToResponse(thread))
}

func (h *ThreadHandler) DeleteThreadHandler(c *fiber.Ctx) error {
	slug := c.Params("slug")
	userID := h.SessionManager.GetUserID(c)

	err := h.ThreadService.Delete(c.Context(), slug, userID)
	if err != nil {
		if errors.Is(err, usecase.ErrThreadNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Thread not found",
			})
		}
		if errors.Is(err, usecase.ErrPermissionDenied) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You do not have permission to delete this thread",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete thread",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Thread deleted",
	})
}

func mapThreadToResponse(t *domain.Thread) fiber.Map {
	recentCommenters := make([]fiber.Map, len(t.RecentCommenters))
	for i, c := range t.RecentCommenters {
		recentCommenters[i] = fiber.Map{
			"id":       c.ID,
			"username": c.Username,
			"avatar":   c.Avatar,
		}
	}

	return fiber.Map{
		"id":            t.ID,
		"title":         t.Title,
		"slug":          t.Slug,
		"content":       t.Content,
		"status":        t.Status,
		"view_count":    t.ViewCount,
		"upvotes":       t.Upvotes,
		"downvotes":     t.Downvotes,
		"replies_count": t.RepliesCount,
		"author": fiber.Map{
			"id":       t.Author.ID,
			"username": t.Author.Username,
			"avatar":   t.Author.Avatar,
		},
		"tags":              serializeTags(t.Tags),
		"created_at":        t.CreatedAt,
		"recent_commenters": recentCommenters,
	}
}

func serializeTags(tags []domain.Tag) []fiber.Map {
	result := make([]fiber.Map, len(tags))
	for i, t := range tags {
		result[i] = fiber.Map{
			"id":    t.ID,
			"name":  t.Name,
			"color": t.Color,
		}
	}
	return result
}

func serializeComments(comments []domain.Comment) []fiber.Map {
	results := make([]fiber.Map, len(comments))
	for i, cm := range comments {
		replies := make([]fiber.Map, len(cm.Replies))
		for j, r := range cm.Replies {
			replies[j] = fiber.Map{
				"id":        r.ID,
				"content":   r.Content,
				"upvotes":   r.Upvotes,
				"downvotes": r.Downvotes,
				"author": fiber.Map{
					"id":       r.Author.ID,
					"username": r.Author.Username,
					"avatar":   r.Author.Avatar,
				},
				"created_at": r.CreatedAt,
			}
		}
		results[i] = fiber.Map{
			"id":        cm.ID,
			"content":   cm.Content,
			"upvotes":   cm.Upvotes,
			"downvotes": cm.Downvotes,
			"author": fiber.Map{
				"id":       cm.Author.ID,
				"username": cm.Author.Username,
				"avatar":   cm.Author.Avatar,
			},
			"replies":    replies,
			"created_at": cm.CreatedAt,
		}
	}
	return results
}
