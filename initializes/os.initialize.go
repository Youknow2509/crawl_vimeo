package initializes

import (
	"fmt"

	"github.com/youknow2509/crawl_vimeo/consts"
	"github.com/youknow2509/crawl_vimeo/global"
	"github.com/youknow2509/crawl_vimeo/utils"
)

func initializeOsSystem() {
	global.OS_SYSTEM = utils.GetOSSystem()
	
	switch global.OS_SYSTEM {
	case consts.OS_LINUX:
		fmt.Println("Operating System: Linux")
	case consts.OS_MACOS:
		fmt.Println("Operating System: macOS")
	case consts.OS_WINDOWS:
		fmt.Println("Operating System: Windows")	
	case consts.OS_UBUNTU:
		fmt.Println("Operating System: Ubuntu")
	default:
		fmt.Println("Operating System: Unknown")
	}
}