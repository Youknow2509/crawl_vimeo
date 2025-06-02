package global

import (
	"github.com/youknow2509/crawl_vimeo/services"
	"google.golang.org/api/youtube/v3"
)

var (
	// OS System
	OS_SYSTEM   = 0                  // Default OS system type
	OS_EXECUTOR services.IOSExecutor // Global OS executor instance

	// YouTube Services
	YTB_SERVICE        *youtube.Service
	YTB_SERVICE_ACTION services.IYoutubeService

	// m3u8 service
	M3U8_SERVICE services.IM3u8
)
