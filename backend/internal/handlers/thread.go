package handlers

import (
	"community-forum/backend/internal/lib"
	"community-forum/backend/internal/middleware"
	"community-forum/backend/internal/models"
	"math"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ThreadHandler struct {
	DB *gorm.DB
}

func NewThreadHandler(db *gorm.DB) *ThreadHandler {
	return &ThreadHandler{DB: db}
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

	if len(req.Title) < 5 || len(req.Title) > 255 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title must be between 5 and 255 characters",
		})
	}

	if len(req.Content) < 10 || len(req.Content) > 50000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Content must be between 10 and 50000 characters",
		})
	}

	if req.Status == "" {
		req.Status = "draft"
	}
	if req.Status != "draft" && req.Status != "published" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Status must be draft or published",
		})
	}

	slug, err := lib.GenerateUniqueSlug(req.Title, h.DB, "threads", "slug")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate slug",
		})
	}

	userID := middleware.GetUserID(c)

	var tags []models.Tag
	if len(req.Tags) > 0 {
		h.DB.Where("name IN ?", req.Tags).Find(&tags)
	}

	thread := models.Thread{
		Title:    req.Title,
		Slug:     slug,
		Content:  req.Content,
		Status:   req.Status,
		AuthorID: userID,
		Tags:     tags,
	}

	if result := h.DB.Create(&thread); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create thread",
		})
	}

	h.DB.Preload("Author").Preload("Tags").First(&thread, thread.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":            thread.ID,
		"title":         thread.Title,
		"slug":          thread.Slug,
		"content":       thread.Content,
		"status":        thread.Status,
		"view_count":    thread.ViewCount,
		"upvotes":       thread.Upvotes(h.DB),
		"downvotes":     thread.Downvotes(h.DB),
		"replies_count": thread.RepliesCount(h.DB),
		"author": fiber.Map{
			"id":       thread.Author.ID,
			"username": thread.Author.Username,
			"avatar":   thread.Author.Avatar,
		},
		"tags":       serializeTags(thread.Tags),
		"created_at": thread.CreatedAt,
	})
}

func (h *ThreadHandler) ListThreadsHandler(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "5"))
	sort := c.Query("sort", "latest")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 5
	}
	if pageSize > 50 {
		pageSize = 50
	}

	query := h.DB.Model(&models.Thread{}).Where("status = ?", "published")

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize

	var threads []models.Thread
	dbQuery := h.DB.Where("status = ?", "published")

	switch sort {
	case "oldest":
		dbQuery = dbQuery.Order("created_at ASC")
	case "votes":
		dbQuery = dbQuery.Order("(SELECT COALESCE(SUM(CASE WHEN value = 1 THEN 1 WHEN value = -1 THEN -1 ELSE 0 END), 0) FROM votes WHERE votes.thread_id = threads.id) DESC, created_at DESC")
	default:
		dbQuery = dbQuery.Order("created_at DESC")
	}

	dbQuery = dbQuery.Preload("Author").Preload("Tags").
		Offset(offset).Limit(pageSize).
		Find(&threads)

	if dbQuery.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch threads",
		})
	}

	type threadResult struct {
		ID           uint         `json:"id"`
		Title        string       `json:"title"`
		Slug         string       `json:"slug"`
		Content      string       `json:"content"`
		Status       string       `json:"status"`
		ViewCount    uint         `json:"view_count"`
		Upvotes      int64        `json:"upvotes"`
		Downvotes    int64        `json:"downvotes"`
		RepliesCount int64        `json:"replies_count"`
		Author       models.User  `json:"author"`
		Tags         []models.Tag `json:"tags"`
		CreatedAt    time.Time    `json:"created_at"`
	}

	results := make([]threadResult, len(threads))
	for i, t := range threads {
		results[i] = threadResult{
			ID:           t.ID,
			Title:        t.Title,
			Slug:         t.Slug,
			Content:      t.Content,
			Status:       t.Status,
			ViewCount:    t.ViewCount,
			Upvotes:      t.Upvotes(h.DB),
			Downvotes:    t.Downvotes(h.DB),
			RepliesCount: t.RepliesCount(h.DB),
			Author:       t.Author,
			Tags:         t.Tags,
			CreatedAt:    t.CreatedAt,
		}
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
	var thread models.Thread

	oneWeekAgo := time.Now().Add(-7 * 24 * time.Hour)

	err := h.DB.
		Where("status = ? AND created_at >= ?", "published", oneWeekAgo).
		Joins("LEFT JOIN votes ON votes.thread_id = threads.id").
		Select("threads.*, COALESCE(SUM(CASE WHEN votes.value = 1 THEN 1 WHEN votes.value = -1 THEN -1 ELSE 0 END), 0) as score").
		Group("threads.id").
		Order("score DESC").
		Preload("Author").
		Preload("Tags").
		First(&thread).Error

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No featured thread found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":            thread.ID,
		"title":         thread.Title,
		"slug":          thread.Slug,
		"content":       thread.Content,
		"status":        thread.Status,
		"view_count":    thread.ViewCount,
		"upvotes":       thread.Upvotes(h.DB),
		"downvotes":     thread.Downvotes(h.DB),
		"replies_count": thread.RepliesCount(h.DB),
		"author": fiber.Map{
			"id":       thread.Author.ID,
			"username": thread.Author.Username,
			"avatar":   thread.Author.Avatar,
		},
		"tags":       serializeTags(thread.Tags),
		"created_at": thread.CreatedAt,
	})
}

func (h *ThreadHandler) TrendingThreadsHandler(c *fiber.Ctx) error {
	var threads []models.Thread

	err := h.DB.
		Where("status = ?", "published").
		Joins("LEFT JOIN votes ON votes.thread_id = threads.id").
		Select("threads.*, COALESCE(SUM(CASE WHEN votes.value = 1 THEN 1 WHEN votes.value = -1 THEN -1 ELSE 0 END), 0) as score").
		Group("threads.id").
		Order("score DESC").
		Limit(3).
		Preload("Author").
		Preload("Tags").
		Find(&threads).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch trending threads",
		})
	}

	type trendingResult struct {
		ID           uint         `json:"id"`
		Title        string       `json:"title"`
		Slug         string       `json:"slug"`
		Content      string       `json:"content"`
		Status       string       `json:"status"`
		ViewCount    uint         `json:"view_count"`
		Upvotes      int64        `json:"upvotes"`
		Downvotes    int64        `json:"downvotes"`
		RepliesCount int64        `json:"replies_count"`
		Author       models.User  `json:"author"`
		Tags         []models.Tag `json:"tags"`
		CreatedAt    time.Time    `json:"created_at"`
	}

	results := make([]trendingResult, len(threads))
	for i, t := range threads {
		results[i] = trendingResult{
			ID:           t.ID,
			Title:        t.Title,
			Slug:         t.Slug,
			Content:      t.Content,
			Status:       t.Status,
			ViewCount:    t.ViewCount,
			Upvotes:      t.Upvotes(h.DB),
			Downvotes:    t.Downvotes(h.DB),
			RepliesCount: t.RepliesCount(h.DB),
			Author:       t.Author,
			Tags:         t.Tags,
			CreatedAt:    t.CreatedAt,
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"threads": results,
	})
}

func (h *ThreadHandler) GetThreadHandler(c *fiber.Ctx) error {
	slug := c.Params("slug")

	var thread models.Thread
	if result := h.DB.Where("slug = ?", slug).First(&thread); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Thread not found",
		})
	}

	h.DB.Model(&models.Thread{}).Where("id = ?", thread.ID).Update("view_count", gorm.Expr("view_count + 1"))

	h.DB.Preload("Author").Preload("Tags").
		Preload("Comments", "parent_id IS NULL").
		Preload("Comments.Replies").
		Preload("Comments.Author").
		Preload("Comments.Replies.Author").
		First(&thread, thread.ID)

	type commentResult struct {
		ID        uint            `json:"id"`
		Content   string          `json:"content"`
		Upvotes   int64           `json:"upvotes"`
		Downvotes int64           `json:"downvotes"`
		Author    fiber.Map       `json:"author"`
		Replies   []commentResult `json:"replies"`
		CreatedAt time.Time       `json:"created_at"`
	}

	serializeComments := func(comments []models.Comment) []commentResult {
		results := make([]commentResult, len(comments))
		for i, cm := range comments {
			replies := make([]commentResult, len(cm.Replies))
			for j, r := range cm.Replies {
				replies[j] = commentResult{
					ID:        r.ID,
					Content:   r.Content,
					Upvotes:   r.Upvotes(h.DB),
					Downvotes: r.Downvotes(h.DB),
					Author: fiber.Map{
						"id":       r.Author.ID,
						"username": r.Author.Username,
						"avatar":   r.Author.Avatar,
					},
					CreatedAt: r.CreatedAt,
				}
			}
			results[i] = commentResult{
				ID:        cm.ID,
				Content:   cm.Content,
				Upvotes:   cm.Upvotes(h.DB),
				Downvotes: cm.Downvotes(h.DB),
				Author: fiber.Map{
					"id":       cm.Author.ID,
					"username": cm.Author.Username,
					"avatar":   cm.Author.Avatar,
				},
				Replies:   replies,
				CreatedAt: cm.CreatedAt,
			}
		}
		return results
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":            thread.ID,
		"title":         thread.Title,
		"slug":          thread.Slug,
		"content":       thread.Content,
		"status":        thread.Status,
		"view_count":    thread.ViewCount + 1,
		"upvotes":       thread.Upvotes(h.DB),
		"downvotes":     thread.Downvotes(h.DB),
		"replies_count": thread.RepliesCount(h.DB),
		"author": fiber.Map{
			"id":       thread.Author.ID,
			"username": thread.Author.Username,
			"avatar":   thread.Author.Avatar,
		},
		"tags":       serializeTags(thread.Tags),
		"comments":   serializeComments(thread.Comments),
		"created_at": thread.CreatedAt,
	})
}

type UpdateThreadRequest struct {
	Title   *string  `json:"title"`
	Content *string  `json:"content"`
	Status  *string  `json:"status"`
	Tags    []string `json:"tags"`
}

func (h *ThreadHandler) UpdateThreadHandler(c *fiber.Ctx) error {
	slug := c.Params("slug")

	var thread models.Thread
	if result := h.DB.Where("slug = ?", slug).First(&thread); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Thread not found",
		})
	}

	userID := middleware.GetUserID(c)
	if thread.AuthorID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You do not have permission to update this thread",
		})
	}

	var req UpdateThreadRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Title != nil {
		if len(*req.Title) < 5 || len(*req.Title) > 255 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Title must be between 5 and 255 characters",
			})
		}
		newSlug, err := lib.GenerateUniqueSlug(*req.Title, h.DB, "threads", "slug")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to regenerate slug",
			})
		}
		thread.Title = *req.Title
		thread.Slug = newSlug
	}

	if req.Content != nil {
		if len(*req.Content) < 10 || len(*req.Content) > 50000 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Content must be between 10 and 50000 characters",
			})
		}
		thread.Content = *req.Content
	}

	if req.Status != nil {
		if *req.Status != "draft" && *req.Status != "published" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Status must be draft or published",
			})
		}
		thread.Status = *req.Status
	}

	if req.Tags != nil {
		var tags []models.Tag
		if len(req.Tags) > 0 {
			h.DB.Where("name IN ?", req.Tags).Find(&tags)
		}
		thread.Tags = tags
	}

	if result := h.DB.Save(&thread); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update thread",
		})
	}

	h.DB.Preload("Author").Preload("Tags").First(&thread, thread.ID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":            thread.ID,
		"title":         thread.Title,
		"slug":          thread.Slug,
		"content":       thread.Content,
		"status":        thread.Status,
		"view_count":    thread.ViewCount,
		"upvotes":       thread.Upvotes(h.DB),
		"downvotes":     thread.Downvotes(h.DB),
		"replies_count": thread.RepliesCount(h.DB),
		"author": fiber.Map{
			"id":       thread.Author.ID,
			"username": thread.Author.Username,
			"avatar":   thread.Author.Avatar,
		},
		"tags":       serializeTags(thread.Tags),
		"created_at": thread.CreatedAt,
	})
}

func (h *ThreadHandler) DeleteThreadHandler(c *fiber.Ctx) error {
	slug := c.Params("slug")

	var thread models.Thread
	if result := h.DB.Where("slug = ?", slug).First(&thread); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Thread not found",
		})
	}

	userID := middleware.GetUserID(c)
	if thread.AuthorID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You do not have permission to delete this thread",
		})
	}

	h.DB.Delete(&thread)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Thread deleted",
	})
}

func serializeTags(tags []models.Tag) []fiber.Map {
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
