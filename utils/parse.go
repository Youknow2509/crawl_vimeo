package utils

import (
	"encoding/json"
	"fmt"
)

// ParsePlayerConfig parses JSON string to map for easier access
func ParsePlayerConfig(jsonString string) (map[string]interface{}, error) {
    var config map[string]interface{}
    err := json.Unmarshal([]byte(jsonString), &config)
    if err != nil {
        return nil, fmt.Errorf("failed to parse JSON: %v", err)
    }
    return config, nil
}