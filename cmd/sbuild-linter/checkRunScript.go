package main

func (v *Validator) validateRunScript() ([]byte, error, string) {
	if xExec, ok := v.data["x_exec"].(map[string]interface{}); ok {
		shell, _ := xExec["shell"].(string)
		if script, ok := xExec["run"].(string); ok {
			_, err, warning := v.validateScript(script, shell)
			if err != nil {
				return nil, err, warning
			}
			return nil, nil, warning
		}
	}
	return nil, nil, ""
}
