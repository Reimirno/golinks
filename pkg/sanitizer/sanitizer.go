package sanitizer

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/reimirno/golinks/pkg/types"
)

// We always sanitize input before persisting it
// So, not much sanitization is needed when reading it back
func SanitizeInput(m types.Mapper, pair *types.PathUrlPair) error {
	pair.Mapper = m.GetName()
	canonicalPath, err := CanonicalizePath(pair.Path)
	if err != nil {
		return err
	}
	pair.Path = canonicalPath
	canonicalUrl, err := CanonicalizeUrl(pair.Url)
	if err != nil {
		return err
	}
	pair.Url = canonicalUrl
	pair.UseCount = 0
	return nil
}

// Since we always sanitize input before persisting it
// We don't need to do much when reading it back
// This function is really just a no-op
func SanitizeOutput(m types.Mapper, pair *types.PathUrlPair) {
	if m.Readonly() {
		pair.UseCount = 0
	}
}

// Make sure mapIn is already assigned to m
func SanitizeInputMap(m types.Mapper, mapIn *types.PathUrlPairMap) error {
	clone := make(types.PathUrlPairMap)
	for key, pair := range *mapIn {
		canonicalPath, err := CanonicalizePath(key)
		if err != nil {
			return err
		}
		err = SanitizeInput(m, pair)
		if err != nil {
			return err
		}
		clone[canonicalPath] = pair
	}
	*mapIn = clone
	return nil
}

// CanonicalizePath does a few things:
// Processes path:
// - trims leading and trailing slashes
// - removes underscore, hyphen and dot in string
// - ensures path begins with a slash
// - replaces multiple consecutive slashes with a single slash
// Validates path:
// - ensures path is not "/" or "/d"
// - ensures path are all properly escaped using url.Parse
func CanonicalizePath(path string) (string, error) {
	// process path
	path = strings.Trim(path, "/")
	path = regexp.MustCompile("[_.-]").ReplaceAllString(path, "")
	path = regexp.MustCompile("/+").ReplaceAllString(path, "/")
	path = "/" + path

	// validate path
	if path == "" || path == "/" || path == "/d" || strings.HasPrefix(path, "/d/") {
		return "", ErrInvalidPath(path, "path is reserved")
	}
	urlParsed, err := url.Parse(path)
	if err != nil {
		return "", ErrInvalidPath(path, fmt.Sprintf("path is invalid url: %s", err.Error()))
	}

	return urlParsed.String(), nil
}

// CanonicalizeUrl trims spaces from the url
// It does not do too much, because we just blindly send it to user
// It is up to the user to ensure the url is correct
func CanonicalizeUrl(url string) (string, error) {
	url = strings.Trim(url, " ")
	return url, nil
}
