package utils

import (
	"os"
	"runtime"
	"strings"

	"github.com/youknow2509/crawl_vimeo/consts"
)

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
