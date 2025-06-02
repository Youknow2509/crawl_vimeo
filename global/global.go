package global

import (
	"github.com/youknow2509/crawl_vimeo/services"
	"github.com/youknow2509/crawl_vimeo/utils"
	"google.golang.org/api/youtube/v3"
)

var (
	// OS System
	OS_SYSTEM   = 0               // Default OS system type
	OS_EXECUTOR *utils.OSExecutor // Global OS executor instance

	// YouTube Services
	YTB_SERVICE        *youtube.Service
	YTB_SERVICE_ACTION services.IYoutubeService

	// m3u8 service
	M3U8_SERVICE services.IM3u8
)

// GetOSExecutor returns the global OS executor, initializing if needed
func GetOSExecutor() *utils.OSExecutor {
	if OS_EXECUTOR == nil {
		// Auto-initialize if not set
		osType := utils.GetCurrentOS()
		OS_SYSTEM = osType
		OS_EXECUTOR = utils.NewOSExecutor(osType)
	}
	return OS_EXECUTOR
}

// SetOSExecutor sets the global OS executor
func SetOSExecutor(executor *utils.OSExecutor) {
	OS_EXECUTOR = executor
}
