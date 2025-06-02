package utils

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	PATTERN_WINDOW_PLAYER_CONFIG  = `<script>window\.playerConfig\s*=\s*(\{.*?\})</script>`
	PATTERN_AUTH_TOKEN_URL        = `exp=(\d+)~acl=([^~]+)~hmac=([a-f0-9]+)/([a-f0-9\-]+)`
	PATTERN_SAFE_FILENAME         = `[<>:"/\\|?*]`
	PATTERN_VIDEO_ID_FROM_VIMEO_1 = `^vimeo\.com/(\d+)`
	PATTERN_VIDEO_ID_FROM_VIMEO_2 = `vimeo\.com/(\d+)`
)

// ExtractVideoIDVimeoFromURL extracts the video ID from a Vimeo URL
func ExtractVideoIDVimeoFromURL(url string) string {
	// Pattern for Vimeo URLs: https://vimeo.com/{video_id}/... or https://vimeo.com/{video_id}
	// The video ID is typically the first number after vimeo.com/

	// Remove protocol and www if present
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "www.")

	// Pattern to match vimeo.com/VIDEO_ID
	re := regexp.MustCompile(PATTERN_VIDEO_ID_FROM_VIMEO_1)
	matches := re.FindStringSubmatch(url)

	if len(matches) >= 2 {
		return matches[1]
	}

	// Alternative pattern for direct video ID extraction
	re2 := regexp.MustCompile(PATTERN_VIDEO_ID_FROM_VIMEO_2)
	matches2 := re2.FindStringSubmatch(url)

	if len(matches2) >= 2 {
		return matches2[1]
	}

	return ""
}

// generateSafeFilename creates a safe filename from title
func GenerateSafeFilename(title string) string {
	// Remove or replace invalid characters
	reg := regexp.MustCompile(PATTERN_SAFE_FILENAME)
	safe := reg.ReplaceAllString(title, "_")

	// Limit length
	if len(safe) > 100 {
		safe = safe[:100]
	}

	// Remove spaces and replace with underscores
	safe = strings.ReplaceAll(safe, " ", "_")

	return safe
}

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
