package lib

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

var nonAlphaNumeric = regexp.MustCompile(`[^a-z0-9]+`)
var multiHyphen = regexp.MustCompile(`-{2,}`)

func GenerateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = nonAlphaNumeric.ReplaceAllString(slug, "-")
	slug = multiHyphen.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}

func GenerateUniqueSlug(title string, db *gorm.DB, tableName string, column string) (string, error) {
	base := GenerateSlug(title)
	if base == "" {
		base = "thread"
	}

	slug := base
	for i := 1; i <= 10; i++ {
		var count int64
		if err := db.Table(tableName).Where(column+" = ?", slug).Count(&count).Error; err != nil {
			return "", fmt.Errorf("failed to check slug uniqueness: %w", err)
		}
		if count == 0 {
			return slug, nil
		}
		slug = fmt.Sprintf("%s-%d", base, i)
	}

	suffix := rand.Intn(9000) + 1000
	slug = fmt.Sprintf("%s-%d", base, suffix)

	var count int64
	if err := db.Table(tableName).Where(column+" = ?", slug).Count(&count).Error; err != nil {
		return "", fmt.Errorf("failed to check slug uniqueness: %w", err)
	}
	if count == 0 {
		return slug, nil
	}

	return fmt.Sprintf("%s-%d", base, rand.Intn(99999-10000)+10000), nil
}
