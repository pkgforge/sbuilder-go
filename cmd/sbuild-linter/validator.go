package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Validator struct {
	level *ValidationLevel
	node  *yaml.Node
	raw   []byte
	data  map[string]interface{}
	file  string
}

func NewValidator(filePath string) (*Validator, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		Log.Error("Failed to read file", "file", filePath, "error", err)
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var node yaml.Node
	if err := yaml.Unmarshal(data, &node); err != nil {
		Log.Error("Failed to parse YAML", "error", err)
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	var dataMap map[string]interface{}
	if err := yaml.Unmarshal(data, &dataMap); err != nil {
		Log.Error("Failed to unmarshal YAML to map", "error", err)
		return nil, fmt.Errorf("failed to unmarshal YAML to map: %w", err)
	}

	return &Validator{
		node: &node,
		raw:  data,
		data: dataMap,
		file: filePath,
	}, nil
}

func (v *Validator) updateFields() error {
	var node yaml.Node
	if err := yaml.Unmarshal(v.raw, &node); err != nil {
		return fmt.Errorf("failed to unmarshal raw data to node: %w", err)
	}
	v.node = &node

	var data map[string]interface{}
	if err := yaml.Unmarshal(v.raw, &data); err != nil {
		return fmt.Errorf("failed to unmarshal raw data to map: %w", err)
	}
	v.data = data

	return nil
}

func (v *Validator) editNode(path []string, handler func(*yaml.Node) (bool, error)) ([]byte, error) {
	var navigate func(*yaml.Node, []string) (*yaml.Node, error)
	navigate = func(n *yaml.Node, remainingPath []string) (*yaml.Node, error) {
		if len(remainingPath) == 0 {
			return n, nil
		}

		if n.Kind != yaml.MappingNode {
			return nil, fmt.Errorf("expected mapping node at path %v", path)
		}

		for i := 0; i < len(n.Content); i += 2 {
			if n.Content[i].Value == remainingPath[0] {
				return navigate(n.Content[i+1], remainingPath[1:])
			}
		}
		return nil, fmt.Errorf("path %v not found", path)
	}

	targetNode, err := navigate(v.node, path)
	if err != nil {
		return nil, err
	}

	modified, err := handler(targetNode)
	if err != nil {
		return nil, err
	}

	if !modified {
		return nil, nil
	}

	data, err := yaml.Marshal(v.node)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal modified data: %w", err)
	}

	v.raw = data
	if err := v.updateFields(); err != nil {
		return nil, fmt.Errorf("failed to update fields: %w", err)
	}

	return data, nil
}

func (v *Validator) ValidateAll(pkgverFlag bool) (validatedData []byte, warningCount int, err error) {
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
		{"Run Script Validation", v.validateRunScript},
	}

	if pkgverFlag {
		checks = append(checks, struct {
			name string
			fn   func() (validatedData []byte, err error, warn string)
		}{"PkgVer Script Validation", v.validatePkgverScript})
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
