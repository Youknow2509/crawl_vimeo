package consts

// const v.v
const (
	YTB_USER_TIME_EXPIRED_AUTH_REDIRECT = 10 // 10 minutes
)

// scope use
const (
	YTB_SCOPE_MANAGET_YOUTUBE_ACCOUNT      = "https://www.googleapis.com/auth/youtube"
	YTB_SCOPE_MANAGET_YOUTUBE_VIDEOS       = "https://www.googleapis.com/auth/youtube.force-ssl"
	YTB_SCOPE_MANAGET_VIEW_YOUTUBE_ACCOUNT = "https://www.googleapis.com/auth/youtube.readonly"
	YTB_SCOPE_MANAGET_VIDEO                = "https://www.googleapis.com/auth/youtube.upload"
)

// token
const (
	YTB_USER_TOKEN_VALIDATE                 = 0
	YTB_USER_TOKEN_EXPIRED                  = 1
	YTB_USER_TOKEN_INVALIDATED              = 2
	YTB_USER_TOKEN_REFRESHED                = 3
	YTB_USER_TOKEN_NOT_FOUND                = 4
	YTB_USER_TOKEN_NOT_AUTHORIZED           = 5
	YTB_USER_TOKEN_NOT_FOUND_OR_EXPIRED     = 6
	YTB_USER_TOKEN_NOT_FOUND_OR_INVALIDATED = 7
)

//
const (
	YTB_CLIENT_SECRET_PATH = "secrets/client_secret.json"
	YTB_USER_AUTH_FILE     = "secrets/user_auth.json"
	YTB_REDIRECT_URI       = "http://localhost:8080"
	YTB_REDIRECT_HOST      = "localhost:8080"
)
