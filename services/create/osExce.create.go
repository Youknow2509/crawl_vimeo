package create

import (
	"github.com/youknow2509/crawl_vimeo/consts"
	"github.com/youknow2509/crawl_vimeo/services"
	"github.com/youknow2509/crawl_vimeo/services/impl"
)

// factory create osExce
func FactoryCreateOsExce(osType int) services.IOSExecutor {
	switch osType {
	case consts.OS_WINDOWS:
		return impl.NewWindowsOsExecutorService()
	case consts.OS_LINUX:
		return impl.NewLinuxOsExecutorService()
	case consts.OS_MACOS:
		return impl.NewMacOsOsExecutorService()
	default:
		return nil
	}
}
