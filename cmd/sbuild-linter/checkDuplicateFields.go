package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"strings"
)

func (v *Validator) validateNoDuplicates() ([]byte, error, string) {
	var checkDupes func(interface{}, string, map[string]bool) (error, []string)
	checkDupes = func(i interface{}, path string, seen map[string]bool) (error, []string) {
		var warnings []string
		switch x := i.(type) {
		case map[string]interface{}:
			localSeen := make(map[string]bool)
			for k, value := range x {
				if localSeen[k] {
					warnings = append(warnings, fmt.Sprintf("%s at %s", k, path))
				}
				localSeen[k] = true

				if seen[k] {
					warnings = append(warnings, fmt.Sprintf("%s at %s", k, path))
				}
				seen[k] = true

				if err, subWarnings := checkDupes(value, path+"."+k, seen); err != nil {
					return err, subWarnings
				} else if len(subWarnings) > 0 {
					warnings = append(warnings, subWarnings...)
				}
			}
		case []interface{}:
			listSeen := make(map[string]bool)
			for i, value := range x {
				strValue, ok := value.(string)
				if !ok {
					continue
				}
				if listSeen[strValue] {
					warnings = append(warnings, fmt.Sprintf("\"%s\" at %s[%d]", strValue, path, i))
				}
				listSeen[strValue] = true
			}
		}
		return nil, warnings
	}

	seen := make(map[string]bool)
	_, warnings := checkDupes(v.data, "root", seen)
	if len(warnings) > 0 {
		v.data = removeDuplicates(v.data).(map[string]interface{})
		data, err := v.editNode([]string{}, func(node *yaml.Node) (bool, error) {
			return true, nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed to edit node: %w", err), ""
		}
		return data, nil, strings.Join(warnings, ", ")
	}
	return nil, nil, ""
}
