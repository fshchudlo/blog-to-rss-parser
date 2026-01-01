package utils

import (
	"fmt"
	"net/url"
	"strings"
)

func ResolveRelativeUrl(basePath string, relativePath string) (string, error) {
	if !strings.HasPrefix(relativePath, "/") {
		return relativePath, nil
	}
	baseUrl, err := url.Parse(basePath)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}
	relativeUrl, err := url.Parse(relativePath)
	if err != nil {
		return "", fmt.Errorf("invalid relative URL: %w", err)
	}
	resolved := baseUrl.ResolveReference(relativeUrl)
	return resolved.String(), nil
}
