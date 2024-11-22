package main

import (
	"os"
	"path/filepath"
)

func removeDuplicates(data interface{}) interface{} {
	switch x := data.(type) {
	case map[string]interface{}:
		newMap := make(map[string]interface{})
		for k, v := range x {
			newMap[k] = removeDuplicates(v)
		}
		return newMap
	case []interface{}:
		seen := make(map[string]bool)
		var newList []interface{}
		for _, v := range x {
			strValue, ok := v.(string)
			if !ok || !seen[strValue] {
				newList = append(newList, v)
				seen[strValue] = true
			}
		}
		return newList
	default:
		return data
	}
}

func writeDataToNewFile(originalFile string, data []byte) error {
	newFile := filepath.Base(originalFile) + ".validated"
	err := os.WriteFile(newFile, data, 0644)
	if err != nil {
		Log.Error("Failed to write processed data to new file", "file", newFile, "error", err)
		return err
	}
	Log.Info("Processed data written to new file", "file", newFile)
	return nil
}
