package global

import (
	"github.com/youknow2509/crawl_vimeo/services"
	"google.golang.org/api/youtube/v3"
)

var (
	OS_SYSTEM          = 0 // Default OS system
	YTB_SERVICE        *youtube.Service
	YTB_SERVICE_ACTION services.IYoutubeService
)
