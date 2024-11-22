package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// validateScript handles the common validation logic for shell scripts
func (v *Validator) validateScript(script string, shell string) (string, error, string) {
	var shebang string
	if shell != "" && shell != "sh" {
		shebang = fmt.Sprintf("#!/usr/bin/env %s", shell)
	} else {
		shell = "sh"
		shebang = "#!/bin/sh"
	}

	// Check if interpreter exists
	if _, err := exec.LookPath(shell); err != nil {
		return "", fmt.Errorf("interpreter (%s) not found in $PATH", shell), ""
	}

	fullScript := shebang + "\n" + script

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "shellcheck-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err), ""
	}
	defer os.Remove(tmpFile.Name())

	// Write script to temp file
	if _, err := tmpFile.WriteString(fullScript); err != nil {
		return "", fmt.Errorf("failed to write script: %w", err), ""
	}
	if err := tmpFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close temp file: %w", err), ""
	}

	cmdError := exec.Command("shellcheck", "--severity=error", tmpFile.Name())
	if output, err := cmdError.CombinedOutput(); err != nil {
		fmt.Println(fmt.Errorf("shellcheck with --severity=error failed: %s", string(output)))
		cmdWarning := exec.Command("shellcheck", "--severity=warning", tmpFile.Name())
		if output, err := cmdWarning.CombinedOutput(); err != nil {
			return "", nil, fmt.Sprintf("shellcheck warnings: %s", output)
		}
	}

	return shell, nil, ""
}

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

func (v *Validator) validatePkgverScript() ([]byte, error, string) {
	if xExec, ok := v.data["x_exec"].(map[string]interface{}); ok {
		shell, _ := xExec["shell"].(string)

		if script, ok := xExec["pkgver"].(string); ok {
			validShell, err, warning := v.validateScript(script, shell)
			if err != nil {
				return nil, err, warning
			}

			// Execute pkgver script
			cmd := exec.Command(validShell, "-c", script)
			output, err := cmd.CombinedOutput()
			if err != nil {
				return nil, fmt.Errorf("failed to execute pkgver script: %w", err), warning
			}

			// Write output to pkgver file
			pkgverFile := filepath.Base(v.file) + ".pkgver"
			if err := os.WriteFile(pkgverFile, output, 0644); err != nil {
				return nil, fmt.Errorf("failed to write pkgver file: %w", err), warning
			}

			return nil, nil, warning
		}
	}
	return nil, nil, ""
}
