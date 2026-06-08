package db

import (
	"context"
	"time"

	"community-forum/backend/internal/domain"
	"community-forum/backend/internal/lib"
	"community-forum/backend/internal/models"

	"gorm.io/gorm"
)

type GORMThreadRepository struct {
	db *gorm.DB
}

func NewGORMThreadRepository(db *gorm.DB) *GORMThreadRepository {
	return &GORMThreadRepository{db: db}
}

func (r *GORMThreadRepository) Create(ctx context.Context, t *domain.Thread, tagNames []string) error {
	m := &models.Thread{
		Title:    t.Title,
		Slug:     t.Slug,
		Content:  t.Content,
		Status:   t.Status,
		AuthorID: t.AuthorID,
	}

	if len(tagNames) > 0 {
		var tags []models.Tag
		r.db.WithContext(ctx).Where("name IN ?", tagNames).Find(&tags)
		m.Tags = tags
	}

	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Preload("Author").Preload("Tags").First(m, m.ID).Error; err != nil {
		return err
	}

	domainThread := threadFromModel(m, r.db)
	*t = *domainThread
	return nil
}

func (r *GORMThreadRepository) List(ctx context.Context, page, pageSize int, sort string) ([]domain.Thread, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&models.Thread{}).Where("status = ?", "published").Count(&total)

	offset := (page - 1) * pageSize
	dbQuery := r.db.WithContext(ctx).
		Select("threads.*, "+
			"(SELECT COALESCE(COUNT(*), 0) FROM votes WHERE votes.thread_id = threads.id AND votes.value = 1) AS upvotes, "+
			"(SELECT COALESCE(COUNT(*), 0) FROM votes WHERE votes.thread_id = threads.id AND votes.value = -1) AS downvotes, "+
			"(SELECT COALESCE(COUNT(*), 0) FROM comments WHERE comments.thread_id = threads.id AND comments.parent_id IS NULL) AS replies_count").
		Where("threads.status = ?", "published")

	switch sort {
	case "oldest":
		dbQuery = dbQuery.Order("threads.created_at ASC")
	case "votes":
		dbQuery = dbQuery.Order("(SELECT COALESCE(SUM(CASE WHEN value = 1 THEN 1 WHEN value = -1 THEN -1 ELSE 0 END), 0) FROM votes WHERE votes.thread_id = threads.id) DESC, threads.created_at DESC")
	default:
		dbQuery = dbQuery.Order("threads.created_at DESC")
	}

	var mThreads []models.Thread
	err := dbQuery.Preload("Author").Preload("Tags").
		Offset(offset).Limit(pageSize).
		Find(&mThreads).Error
	if err != nil {
		return nil, 0, err
	}

	threads := make([]domain.Thread, len(mThreads))
	for i := range mThreads {
		threads[i] = *threadFromModel(&mThreads[i], r.db)
	}

	return threads, total, nil
}

func (r *GORMThreadRepository) ListByUser(ctx context.Context, username string, page, pageSize int) ([]domain.Thread, int64, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	r.db.WithContext(ctx).Model(&models.Thread{}).Where("author_id = ?", user.ID).Count(&total)

	offset := (page - 1) * pageSize
	var mThreads []models.Thread
	err := r.db.WithContext(ctx).
		Select("threads.*, "+
			"(SELECT COALESCE(COUNT(*), 0) FROM votes WHERE votes.thread_id = threads.id AND votes.value = 1) AS upvotes, "+
			"(SELECT COALESCE(COUNT(*), 0) FROM votes WHERE votes.thread_id = threads.id AND votes.value = -1) AS downvotes, "+
			"(SELECT COALESCE(COUNT(*), 0) FROM comments WHERE comments.thread_id = threads.id AND comments.parent_id IS NULL) AS replies_count").
		Where("author_id = ?", user.ID).
		Order("threads.created_at DESC").
		Offset(offset).Limit(pageSize).
		Preload("Author").Preload("Tags").
		Find(&mThreads).Error
	if err != nil {
		return nil, 0, err
	}

	threads := make([]domain.Thread, len(mThreads))
	for i := range mThreads {
		threads[i] = *threadFromModel(&mThreads[i], r.db)
	}

	return threads, total, nil
}

func (r *GORMThreadRepository) GetFeatured(ctx context.Context) (*domain.Thread, error) {
	var m models.Thread
	oneWeekAgo := time.Now().Add(-7 * 24 * time.Hour)

	err := r.db.WithContext(ctx).
		Select("threads.*, "+
			"(SELECT COALESCE(COUNT(*), 0) FROM votes WHERE votes.thread_id = threads.id AND votes.value = 1) AS upvotes, "+
			"(SELECT COALESCE(COUNT(*), 0) FROM votes WHERE votes.thread_id = threads.id AND votes.value = -1) AS downvotes, "+
			"(SELECT COALESCE(COUNT(*), 0) FROM comments WHERE comments.thread_id = threads.id AND comments.parent_id IS NULL) AS replies_count").
		Where("threads.status = ? AND threads.created_at >= ?", "published", oneWeekAgo).
		Order("(SELECT COALESCE(SUM(CASE WHEN value = 1 THEN 1 WHEN value = -1 THEN -1 ELSE 0 END), 0) FROM votes WHERE votes.thread_id = threads.id) DESC").
		Preload("Author").
		Preload("Tags").
		First(&m).Error

	if err != nil {
		return nil, err
	}

	return threadFromModel(&m, r.db), nil
}

func (r *GORMThreadRepository) GetTrending(ctx context.Context) ([]domain.Thread, error) {
	var mThreads []models.Thread

	err := r.db.WithContext(ctx).
		Select("threads.*, "+
			"(SELECT COALESCE(COUNT(*), 0) FROM votes WHERE votes.thread_id = threads.id AND votes.value = 1) AS upvotes, "+
			"(SELECT COALESCE(COUNT(*), 0) FROM votes WHERE votes.thread_id = threads.id AND votes.value = -1) AS downvotes, "+
			"(SELECT COALESCE(COUNT(*), 0) FROM comments WHERE comments.thread_id = threads.id AND comments.parent_id IS NULL) AS replies_count").
		Where("status = ?", "published").
		Order("(SELECT COALESCE(SUM(CASE WHEN value = 1 THEN 1 WHEN value = -1 THEN -1 ELSE 0 END), 0) FROM votes WHERE votes.thread_id = threads.id) DESC").
		Limit(3).
		Preload("Author").
		Preload("Tags").
		Find(&mThreads).Error

	if err != nil {
		return nil, err
	}

	threads := make([]domain.Thread, len(mThreads))
	for i := range mThreads {
		threads[i] = *threadFromModel(&mThreads[i], r.db)
	}

	return threads, nil
}

func (r *GORMThreadRepository) GetBySlug(ctx context.Context, slug string) (*domain.Thread, error) {
	var m models.Thread
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&m).Error; err != nil {
		return nil, err
	}

	// Fetch everything: Author, Tags, and Nested Comments.
	r.db.WithContext(ctx).Preload("Author").Preload("Tags").
		Preload("Comments", "parent_id IS NULL").
		Preload("Comments.Replies").
		Preload("Comments.Author").
		Preload("Comments.Replies.Author").
		First(&m, m.ID)

	return threadFromModel(&m, r.db), nil
}

func (r *GORMThreadRepository) Update(ctx context.Context, t *domain.Thread, tagNames []string) error {
	var m models.Thread
	if err := r.db.WithContext(ctx).Where("id = ?", t.ID).First(&m).Error; err != nil {
		return err
	}

	m.Title = t.Title
	m.Slug = t.Slug
	m.Content = t.Content
	m.Status = t.Status

	if tagNames != nil {
		var tags []models.Tag
		if len(tagNames) > 0 {
			r.db.WithContext(ctx).Where("name IN ?", tagNames).Find(&tags)
		}
		// Clear existing tags and set new ones
		r.db.WithContext(ctx).Model(&m).Association("Tags").Replace(&tags)
	}

	if err := r.db.WithContext(ctx).Save(&m).Error; err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Preload("Author").Preload("Tags").First(&m, m.ID).Error; err != nil {
		return err
	}
	domainThread := threadFromModel(&m, r.db)
	*t = *domainThread

	return nil
}

func (r *GORMThreadRepository) Delete(ctx context.Context, t *domain.Thread) error {
	return r.db.WithContext(ctx).Delete(&models.Thread{}, t.ID).Error
}

func (r *GORMThreadRepository) IncrementViewCount(ctx context.Context, threadID uint) error {
	return r.db.WithContext(ctx).Model(&models.Thread{}).Where("id = ?", threadID).Update("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *GORMThreadRepository) GenerateUniqueSlug(ctx context.Context, title string) (string, error) {
	return lib.GenerateUniqueSlug(title, r.db, "threads", "slug")
}

func threadFromModel(m *models.Thread, db *gorm.DB) *domain.Thread {
	if m == nil {
		return nil
	}
	t := &domain.Thread{
		ID:           m.ID,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		Title:        m.Title,
		Slug:         m.Slug,
		Content:      m.Content,
		Status:       m.Status,
		ViewCount:    m.ViewCount,
		AuthorID:     m.AuthorID,
		Upvotes:      m.Upvotes,
		Downvotes:    m.Downvotes,
		RepliesCount: m.RepliesCount,
	}

	if m.Author.ID != 0 {
		t.Author = domain.User{
			ID:       m.Author.ID,
			Username: m.Author.Username,
			Avatar:   m.Author.Avatar,
		}
	}

	if len(m.Tags) > 0 {
		t.Tags = make([]domain.Tag, len(m.Tags))
		for i, mt := range m.Tags {
			t.Tags[i] = domain.Tag{
				ID:    mt.ID,
				Name:  mt.Name,
				Color: mt.Color,
			}
		}
	}

	if len(m.Comments) > 0 {
		t.Comments = make([]domain.Comment, len(m.Comments))
		for i, mc := range m.Comments {
			t.Comments[i] = *commentFromModel(&mc, db)
		}
	}

	return t
}

func commentFromModel(m *models.Comment, db *gorm.DB) *domain.Comment {
	if m == nil {
		return nil
	}
	c := &domain.Comment{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Content:   m.Content,
		ThreadID:  m.ThreadID,
		AuthorID:  m.AuthorID,
		Upvotes:   m.Upvotes(db),
		Downvotes: m.Downvotes(db),
	}

	if m.ParentID != nil {
		c.ParentID = m.ParentID
	}

	if m.Author.ID != 0 {
		c.Author = domain.User{
			ID:       m.Author.ID,
			Username: m.Author.Username,
			Avatar:   m.Author.Avatar,
		}
	}

	if len(m.Replies) > 0 {
		c.Replies = make([]domain.Comment, len(m.Replies))
		for i, mr := range m.Replies {
			c.Replies[i] = *commentFromModel(&mr, db)
		}
	}

	return c
}
