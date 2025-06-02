package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/youknow2509/crawl_vimeo/consts"
)

// GetCurrentOS detects the current operating system
func GetCurrentOS() int {
	switch runtime.GOOS {
	case "windows":
		return consts.OS_WINDOWS
	case "darwin":
		return consts.OS_MACOS
	case "linux":
		// Check if it's Ubuntu
		if isUbuntu() {
			return consts.OS_UBUNTU
		}
		return consts.OS_LINUX
	default:
		return consts.OS_LINUX // Default to Linux
	}
}

// isUbuntu checks if the current Linux distribution is Ubuntu
func isUbuntu() bool {
	if _, err := os.Stat("/etc/lsb-release"); err == nil {
		content, err := os.ReadFile("/etc/lsb-release")
		if err == nil {
			return strings.Contains(string(content), "Ubuntu")
		}
	}
	return false
}

// GetOSName returns the OS name as string
func GetOSName(osType int) string {
	switch osType {
	case consts.OS_WINDOWS:
		return "Windows"
	case consts.OS_MACOS:
		return "macOS"
	case consts.OS_LINUX:
		return "Linux"
	case consts.OS_UBUNTU:
		return "Ubuntu"
	default:
		return "Unknown"
	}
}

// get url m3u8 in playerConfig
func GetM3U8PathFromPlayerConfig(config map[string]interface{}) (string, error) {
	// request.files.hls.captions
	files, ok := config["request"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("request not found in playerConfig")
	}
	hls, ok := files["files"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("files not found in playerConfig")
	}
	hlsFiles, ok := hls["hls"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("hls not found in playerConfig")
	}
	captions, ok := hlsFiles["captions"].(string)
	if !ok {
		return "", fmt.Errorf("captions not found in playerConfig")
	}
	return captions, nil
}

// ExtractPlayerConfigFromHTML extracts the entire window.playerConfig object
func ExtractPlayerConfigFromHTML(htmlContent string) (string, error) {
    // Regex để tìm toàn bộ window.playerConfig
    pattern := `<script>window\.playerConfig\s*=\s*(\{.*?\})</script>`
    
    re := regexp.MustCompile(pattern)
    match := re.FindStringSubmatch(htmlContent)
    
    if len(match) >= 2 {
        return match[1], nil // ← Trả về JSON string
    }
    
    return "", fmt.Errorf("window.playerConfig not found")
}

// ParsePlayerConfig parses JSON string to map for easier access
func ParsePlayerConfig(jsonString string) (map[string]interface{}, error) {
    var config map[string]interface{}
    err := json.Unmarshal([]byte(jsonString), &config)
    if err != nil {
        return nil, fmt.Errorf("failed to parse JSON: %v", err)
    }
    return config, nil
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
    pattern := `exp=(\d+)~acl=([^~]+)~hmac=([a-f0-9]+)/([a-f0-9\-]+)`
    re := regexp.MustCompile(pattern)
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


// containsIgnoreCase kiểm tra chứa không phân biệt hoa thường
func ContainsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
