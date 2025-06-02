package initializes

import (
	"log"

	"github.com/youknow2509/crawl_vimeo/global"
	"github.com/youknow2509/crawl_vimeo/services"
	"github.com/youknow2509/crawl_vimeo/services/impl"
)

// ininitialize m3u8 service
func initializeM3u8() {
	// create M3U8 service instance
	services.InitializeM3u8(impl.NewM3u8Service(global.OS_EXECUTOR))
	// set M3U8 service to global variable
	global.M3U8_SERVICE = services.GetM3u8Service()
	// log initialization
	log.Println("M3U8 service initialized successfully")
}
