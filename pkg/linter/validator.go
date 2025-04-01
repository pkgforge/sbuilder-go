package linter

import (
	"fmt"
	"os"

	"github.com/pkgforge/sbuilder-go/pkg/logger"
	"github.com/goccy/go-yaml"
)

type Validator struct {
	level *ValidationLevel
	data  map[string]interface{}
	raw   []byte
	file  string
}

func NewValidator(filePath string) (*Validator, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		logger.Log.Error("Failed to read file", "file", filePath, "error", err)
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var dataMap map[string]interface{}
	if err := yaml.Unmarshal(data, &dataMap); err != nil {
		logger.Log.Error("Failed to unmarshal YAML to map", "error", err)
		return nil, fmt.Errorf("failed to unmarshal YAML to map: %w", err)
	}

	return &Validator{
		level: &ValidationLevel{logger: logger.Log},
		data:  dataMap,
		raw:   data,
		file:  filePath,
	}, nil
}

func (v *Validator) updateFields() error {
	var dataMap map[string]interface{}
	if err := yaml.Unmarshal(v.raw, &dataMap); err != nil {
		return fmt.Errorf("failed to unmarshal raw data to map: %w", err)
	}
	v.data = dataMap
	return nil
}

func (v *Validator) updateField(path []string, value interface{}) ([]byte, error) {
	// Handle root-level updates (when path is nil or empty)
	if path == nil || len(path) == 0 {
		v.data = value.(map[string]interface{})
	} else {
		// Start from the root of the data
		current := v.data

		// Traverse the path, creating nested maps as needed
		for i, key := range path {
			// If we're at the last element, set the value
			if i == len(path)-1 {
				current[key] = value
				break
			}

			// For non-last elements, ensure we have a map to traverse into
			if next, ok := current[key].(map[string]interface{}); ok {
				current = next
			} else {
				// Create a new map if the path doesn't exist
				newMap := make(map[string]interface{})
				current[key] = newMap
				current = newMap
			}
		}
	}

	// Marshal the updated data back to YAML
	data, err := yaml.Marshal(v.data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal modified data: %w", err)
	}

	v.raw = data
	if err := v.updateFields(); err != nil {
		return nil, fmt.Errorf("failed to update fields: %w", err)
	}

	return data, nil
}

func (v *Validator) ValidateAll(pkgverFlag, noShellcheckFlag bool) (validatedData []byte, warningCount int, err error) {
	warningCount = 0

	checks := []struct {
		name string
		fn   func() (validatedData []byte, err error, warn string)
	}{
		{"YAML Validation", v.validateYAML},
		{"Shebang Check", v.checkShebang},
		{"_disabled Check", v.validateDisabledField},
		{"Single Value Fields Check", v.validateSingleValueFields},
		{"Empty Fields Check", v.validateNoEmptyFields},
		{"Duplicate Fields Check", v.validateNoDuplicates},
		{"Enforced Fields Check", v.validateEnforcedFields},
		{"Required Fields Check", v.validateRequiredFields},
		{"Categories Validation", v.validateCategories},
		{"URL Fields Validation", v.validateURLs},
		{"PKG Id Validation", v.validatePkgID},
		{"Run Script Validation", func() ([]byte, error, string) { return v.validateRunScript(noShellcheckFlag) }},
	}

	if pkgverFlag {
		checks = append(checks, struct {
			name string
			fn   func() (validatedData []byte, err error, warn string)
		}{"PkgVer Script Validation", func() ([]byte, error, string) { return v.validatePkgverScript(noShellcheckFlag) }})
	}

	for _, check := range checks {
		data, err, warn := check.fn()
		if err != nil {
			v.level.LogError(check.name, err)
			return nil, warningCount, fmt.Errorf("%s: %w", check.name, err)
		}

		if warn != "" {
			v.level.LogWarn(check.name, warn)
			warningCount++
		} else {
			v.level.LogSuccess(check.name, "")
		}

		if data != nil {
			v.raw = data
			if err := v.updateFields(); err != nil {
				return nil, warningCount, fmt.Errorf("failed to update fields: %w", err)
			}
			validatedData = data
		}
	}

	if validatedData == nil {
		return v.raw, warningCount, nil
	}

	return validatedData, warningCount, nil
}
