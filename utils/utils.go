package utils

import (
	"fmt"
	"os"
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

// containsIgnoreCase kiểm tra chứa không phân biệt hoa thường
func ContainsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
