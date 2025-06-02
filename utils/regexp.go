package utils

import (
	"fmt"
	"regexp"
)

const (
	PATTERN_WINDOW_PLAYER_CONFIG = `<script>window\.playerConfig\s*=\s*(\{.*?\})</script>`
	PATTERN_AUTH_TOKEN_URL       = `exp=(\d+)~acl=([^~]+)~hmac=([a-f0-9]+)/([a-f0-9\-]+)`
)

// ExtractPlayerConfigFromHTML extracts the entire window.playerConfig object
func ExtractPlayerConfigFromHTML(htmlContent string) (string, error) {
	// Regex để tìm toàn bộ window.playerConfig
	re := regexp.MustCompile(PATTERN_WINDOW_PLAYER_CONFIG)
	match := re.FindStringSubmatch(htmlContent)

	if len(match) >= 2 {
		return match[1], nil // ← Trả về JSON string
	}

	return "", fmt.Errorf("window.playerConfig not found")
}

// ExtractAuthFromPlayerConfig extracts auth token from parsed config
func ExtractAuthFromPlayerConfig(config map[string]interface{}) (map[string]string, error) {
	// Navigate to the URL
	request, ok := config["request"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("request not found")
	}

	files, ok := request["files"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("files not found")
	}

	dash, ok := files["dash"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("dash not found")
	}

	cdns, ok := dash["cdns"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("cdns not found")
	}

	akfire, ok := cdns["akfire_interconnect_quic"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("akfire_interconnect_quic not found")
	}

	url, ok := akfire["url"].(string)
	if !ok {
		return nil, fmt.Errorf("url not found")
	}

	// Extract auth token từ URL
	re := regexp.MustCompile(PATTERN_AUTH_TOKEN_URL)
	match := re.FindStringSubmatch(url)

	if len(match) >= 5 {
		return map[string]string{
			"exp":       match[1],
			"acl":       match[2],
			"hmac":      match[3],
			"video_id":  match[4],
			"full_auth": fmt.Sprintf("exp=%s~acl=%s~hmac=%s", match[1], match[2], match[3]),
			"full_url":  url,
		}, nil
	}

	return nil, fmt.Errorf("auth token not found in URL")
}
