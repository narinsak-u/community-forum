package lib

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

// Regex patterns used to clean up strings for URLs.
var nonAlphaNumeric = regexp.MustCompile(`[^a-z0-9]+`) // Matches anything that is NOT a lowercase letter or number
var multiHyphen = regexp.MustCompile(`-{2,}`)          // Matches two or more hyphens in a row

// GenerateSlug converts a title (like "Hello World!") into a URL-friendly slug (like "hello-world").
func GenerateSlug(title string) string {
	slug := strings.ToLower(title)
	// Replace all non-alphanumeric characters with a hyphen
	slug = nonAlphaNumeric.ReplaceAllString(slug, "-")
	// Replace multiple hyphens with a single one
	slug = multiHyphen.ReplaceAllString(slug, "-")
	// Remove leading and trailing hyphens
	slug = strings.Trim(slug, "-")
	return slug
}

// GenerateUniqueSlug ensures that the generated slug isn't already used in the database.
// If "my-thread" exists, it will try "my-thread-1", "my-thread-2", etc.
func GenerateUniqueSlug(title string, db *gorm.DB, tableName string, column string) (string, error) {
	base := GenerateSlug(title)
	if base == "" {
		base = "thread" // Fallback if the title was only special characters
	}

	slug := base
	// Try adding a simple counter suffix (up to 10 attempts)
	for i := 1; i <= 10; i++ {
		var count int64
		// Check if a row with this slug already exists in the specified table and column
		if err := db.Table(tableName).Where(column+" = ?", slug).Count(&count).Error; err != nil {
			return "", fmt.Errorf("failed to check slug uniqueness: %w", err)
		}

		// If count is 0, the slug is available!
		if count == 0 {
			return slug, nil
		}

		// If taken, try appending the counter
		slug = fmt.Sprintf("%s-%d", base, i)
	}

	// If still not unique after 10 tries, append a random 4-digit number
	suffix := rand.Intn(9000) + 1000
	slug = fmt.Sprintf("%s-%d", base, suffix)

	var count int64
	if err := db.Table(tableName).Where(column+" = ?", slug).Count(&count).Error; err != nil {
		return "", fmt.Errorf("failed to check slug uniqueness: %w", err)
	}
	if count == 0 {
		return slug, nil
	}

	// Final fallback: append a larger random number
	return fmt.Sprintf("%s-%d", base, rand.Intn(99999-10000)+10000), nil
}
