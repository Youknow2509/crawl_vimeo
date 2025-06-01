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

// GetOSSystem trả về hằng số OS theo hệ điều hành hiện tại
func GetOSSystem() int {
	switch runtime.GOOS {
	case "darwin":
		return consts.OS_MACOS
	case "linux":
		// Nếu muốn phân biệt Ubuntu, cần kiểm tra file /etc/os-release
		if isUbuntu() {
			return consts.OS_UBUNTU
		}
		return consts.OS_LINUX
	case "windows":
		return consts.OS_WINDOWS
	default:
		return -1 // Không xác định
	}
}

// isUbuntu kiểm tra xem hệ điều hành có phải Ubuntu không
func isUbuntu() bool {
	// Đơn giản nhất là kiểm tra file /etc/os-release (chỉ trên Linux)
	// Nếu không phải Linux, trả về false
	if runtime.GOOS != "linux" {
		return false
	}
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return false
	}
	return containsUbuntu(string(data))
}

// containsUbuntu kiểm tra chuỗi có chứa "ubuntu" không
func containsUbuntu(s string) bool {
	return (len(s) > 0) && (containsIgnoreCase(s, "ubuntu"))
}

// containsIgnoreCase kiểm tra chứa không phân biệt hoa thường
func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
