package main

import (
	"fmt"
	"strings"
)

func (v *Validator) validateURLs() ([]byte, error, string) {
	urlFields := map[string]bool{
		"homepage":    false,
		"src_url":     false,
		"build_asset": false,
		"desktop":     true,
		"icon":        true,
	}

	for field, singleValue := range urlFields {
		urls, exists := v.data[field]
		if !exists {
			continue
		}

		if singleValue {
			urlStr, ok := urls.(string)
			if !ok {
				return nil, fmt.Errorf("%s must be a string", field), ""
			}

			if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
				return nil, fmt.Errorf("%s contains invalid URL: %s", field, urlStr), ""
			}
		} else {
			urlSlice, ok := urls.([]interface{})
			if !ok {
				return nil, fmt.Errorf("%s must be an array", field), ""
			}

			for _, url := range urlSlice {
				urlStr, ok := url.(string)
				if !ok {
					return nil, fmt.Errorf("invalid URL type in %s", field), ""
				}

				if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
					return nil, fmt.Errorf("%s contains invalid URL: %s", field, urlStr), ""
				}
			}
		}
	}

	return nil, nil, ""
}
