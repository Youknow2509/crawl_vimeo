package services

// interfaces
type (
	IM3u8 interface {
		M3U8ToMP4(inputUrl string, outPath string) error
		GetPathM3U8(videoID string) (string, error)
	}
)

// #######################################################

// local variable
var (
	vIM3u8 IM3u8
)

// InitializeM3u8 initializes the M3U8 service
func InitializeM3u8(m IM3u8) {
	vIM3u8 = m
}

// Get instance of M3U8 service
func GetM3u8Service() IM3u8 {
	if vIM3u8 == nil {
		panic("M3U8 service is not initialized")
	}
	return vIM3u8
}

