package slug

import (
	"fmt"
	"github.com/ramk42/mini-url/internal/model"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func Generate(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

var (
	validSlug = regexp.MustCompile(fmt.Sprintf("^[%s]{%d}$", regexp.QuoteMeta(charset), model.Sluglength))
)

func Clean(slug string) string {
	cleanedSlug := strings.TrimSpace(slug)
	cleanedSlug = strings.TrimPrefix(cleanedSlug, "/")
	cleanedSlug = strings.TrimSuffix(cleanedSlug, "/")
	return cleanedSlug
}

func Validate(slug string) error {
	const expectedLength = model.Sluglength
	if len(slug) != expectedLength {
		return fmt.Errorf("invalid length: expected %d, got %d", expectedLength, len(slug))
	}

	if !regexp.MustCompile(fmt.Sprintf("^[A-Za-z0-9]{%d}$", expectedLength)).MatchString(slug) {
		return fmt.Errorf("invalid characters detected")
	}

	return nil
}
