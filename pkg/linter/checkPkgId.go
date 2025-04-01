package linter

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

func (v *Validator) validatePkgID() ([]byte, error, string) {
	// Check if pkg_id exists and is not empty
	if pkgID, exists := v.data["pkg_id"].(string); exists && pkgID != "" {
		if !regexp.MustCompile(`^[a-zA-Z0-9\+\-_\.]+$`).MatchString(pkgID) {
			return nil, fmt.Errorf("pkg_id can only contain alphabets, digits, and the following special characters: + - _ ."), ""
		}
		return nil, nil, ""
	}

	// Get src_url to generate pkg_id
	var srcURL string
	if urls, ok := v.data["src_url"].([]interface{}); ok && len(urls) > 0 {
		if firstURL, ok := urls[0].(string); ok {
			srcURL = firstURL
		} else {
			return nil, fmt.Errorf("src_url is not a valid string"), ""
		}
	} else if singleURL, ok := v.data["src_url"].(string); ok {
		srcURL = singleURL
	} else {
		return nil, fmt.Errorf("src_url is missing or invalid"), ""
	}

	// Generate pkg_id from src_url
	parsedURL, err := url.Parse(srcURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err), ""
	}

	path := parsedURL.Host + parsedURL.Path
	components := strings.Split(path, "/")
	var filteredComponents []string
	for _, comp := range components {
		if comp != "" {
			filteredComponents = append(filteredComponents, comp)
		}
	}
	pkgID := strings.Join(filteredComponents, ".")

	// Replace special characters with dots
	var result strings.Builder
	for _, char := range pkgID {
		if regexp.MustCompile(`[a-zA-Z0-9\+\-_\.]`).MatchString(string(char)) {
			result.WriteRune(char)
		} else {
			result.WriteRune('.')
		}
	}
	pkgID = strings.TrimRight(result.String(), ".")

	// Validate the generated pkg_id
	if !regexp.MustCompile(`^[a-zA-Z0-9\+\-_\.]+$`).MatchString(pkgID) {
		return nil, fmt.Errorf("generated pkg_id can only contain alphabets, digits, and the following special characters: + - _ ."), ""
	}

	// Update the pkg_id field
	updatedData, err := v.updateField([]string{"pkg_id"}, pkgID)
	if err != nil {
		return nil, fmt.Errorf("failed to update field: %w", err), ""
	}

	return updatedData, nil, ""
}
