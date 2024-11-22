package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"net/url"
	"regexp"
	"strings"
)

func (v *Validator) validatePkgID() ([]byte, error, string) {
	data := v.data
	// Check if pkg_id exists and is not empty
	if pkgID, exists := data["pkg_id"].(string); exists && pkgID != "" {
		// Validate existing pkg_id format
		if !regexp.MustCompile(`^[a-zA-Z0-9\+\-_\.]+$`).MatchString(pkgID) {
			return nil, fmt.Errorf("pkg_id can only contain alphabets, digits, and the following special characters: + - _ ."), ""
		}
		return nil, nil, ""
	}

	// Get src_url to generate pkg_id
	var srcURL string
	if urls, ok := data["src_url"].([]interface{}); ok && len(urls) > 0 {
		if firstURL, ok := urls[0].(string); ok {
			srcURL = firstURL
		} else {
			return nil, fmt.Errorf("src_url is not a valid string"), ""
		}
	} else if singleURL, ok := data["src_url"].(string); ok {
		srcURL = singleURL
	} else {
		return nil, fmt.Errorf("src_url is missing or invalid"), ""
	}

	// Generate pkg_id from src_url
	parsedURL, err := url.Parse(srcURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err), ""
	}

	// Remove scheme and join components with dots
	path := parsedURL.Host + parsedURL.Path
	components := strings.Split(path, "/")
	var filteredComponents []string
	for _, comp := range components {
		if comp != "" {
			filteredComponents = append(filteredComponents, comp)
		}
	}
	pkgID := strings.Join(filteredComponents, ".")

	// Replace special characters with dots (except - and _)
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

	// Parse the existing YAML to preserve structure
	var root yaml.Node
	if err := yaml.Unmarshal(v.raw, &root); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err), ""
	}

	// Find the pkg field and insert pkg_id after it
	if len(root.Content) > 0 && root.Content[0].Kind == yaml.MappingNode {
		mapping := root.Content[0]
		for i := 0; i < len(mapping.Content); i += 2 {
			if mapping.Content[i].Value == "pkg" {
				// Create pkg_id nodes
				pkgIDKey := &yaml.Node{
					Kind:  yaml.ScalarNode,
					Tag:   "!!str",
					Value: "pkg_id",
				}
				pkgIDValue := &yaml.Node{
					Kind:  yaml.ScalarNode,
					Tag:   "!!str",
					Value: pkgID,
					Style: yaml.DoubleQuotedStyle,
				}

				// Insert pkg_id nodes right after pkg
				newContent := make([]*yaml.Node, 0, len(mapping.Content)+2)
				newContent = append(newContent, mapping.Content[:i+2]...)
				newContent = append(newContent, pkgIDKey, pkgIDValue)
				newContent = append(newContent, mapping.Content[i+2:]...)
				mapping.Content = newContent

				// Marshal the modified YAML
				modifiedYAML, err := yaml.Marshal(&root)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal YAML: %w", err), ""
				}
				return modifiedYAML, nil, ""
			}
		}
		return nil, fmt.Errorf("pkg field not found"), ""
	}

	return nil, fmt.Errorf("invalid YAML structure"), ""
}
