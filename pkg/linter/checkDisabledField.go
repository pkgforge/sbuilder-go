package linter

import (
	"fmt"
)

func (v *Validator) validateDisabledField() ([]byte, error, string) {
	disabled, exists := v.data["_disabled"]
	if !exists {
		return nil, fmt.Errorf("_disabled field does not exist"), ""
	}

	if _, ok := disabled.(bool); !ok {
		return nil, fmt.Errorf("_disabled field must be a boolean"), ""
	}

	return nil, nil, ""
}
