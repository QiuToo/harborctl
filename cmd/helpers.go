package cmd

import (
	"fmt"
	"strings"
)

// getTagsString extracts tags from the artifact's Tags field and returns them as a formatted string
func getTagsString(tags interface{}) string {
	if tags == nil {
		return ""
	}
	switch v := tags.(type) {
	case []interface{}:
		if len(v) == 0 {
			return ""
		}
		result := make([]string, len(v))
		for i, t := range v {
			if s, ok := t.(string); ok {
				result[i] = s
			}
		}
		return strings.Join(result, ", ")
	case []string:
		if len(v) == 0 {
			return ""
		}
		return strings.Join(v, ", ")
	case string:
		return v
	default:
		return ""
	}
}

// getTagsSlice extracts tags from the artifact's Tags field and returns them as a slice
func getTagsSlice(tags interface{}) []string {
	if tags == nil {
		return nil
	}
	switch v := tags.(type) {
	case []interface{}:
		if len(v) == 0 {
			return nil
		}
		result := make([]string, 0, len(v))
		for _, t := range v {
			if s, ok := t.(string); ok {
				result = append(result, s)
			}
		}
		return result
	case []string:
		return v
	case string:
		if v != "" {
			return []string{v}
		}
		return nil
	default:
		return nil
	}
}

var _ = fmt.Println // avoid unused import warning