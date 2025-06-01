package initializes

import "github.com/youknow2509/crawl_vimeo/services"

func Initialize() {
	initializeOsSystem()

	url := "https://vod-adaptive-ak.vimeocdn.com/exp=1748799633~acl=%2F515c7aa7-e99b-4f16-8c56-85fe2f10bdc7%2F%2A~hmac=f47611463708c557cbf64372d81fb1626a5c93b13b55966f4e152ee5b48fdfb1/515c7aa7-e99b-4f16-8c56-85fe2f10bdc7/v2/playlist/av/primary/sub/19775540-c-en-x-autogen/prot/cXNyPTE/playlist.m3u8?ext-subs=1&omit=av1-hevc-opus&pathsig=8c953e4f~NRC5qRrF62xa4oycrge1bgoKRFryiPvHAIA4AAm8HFk&qsr=1&r=dXM%3D&rh=28Enmi&sf=ts"
	err := services.M3U8ToMP4(url, "~/data/1/output.mp4")
	if err != nil {
		panic(err)
	}
}