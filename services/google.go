package services

import (
    "context"
    "errors"
    "golang.org/x/oauth2"
    "google.golang.org/api/option"
    "google.golang.org/api/youtube/v3"
)

// interface for google services
type (
    // For google
    IGoogleService interface {
        GetScopesServer() []string
        GetConfigServer(file string) (*oauth2.Config, error)
    }

    // For google user
    IGoogleUserService interface {
        ValidateToken(ctx context.Context, token *oauth2.Token) (int, bool)
        RefreshToken(ctx context.Context, config *oauth2.Config, token *oauth2.Token) (*oauth2.Token, error)
        GetUrlRedirectAuth(ctx context.Context, config *oauth2.Config) (string, error)
        GetTokenFromWeb(ctx context.Context, config *oauth2.Config, fileSave string) (*oauth2.Token, error)
        GetTokenFromFile(file string) (*oauth2.Token, error)
        SaveToken(path string, token *oauth2.Token) error
        GetYoutubeService(ctx context.Context, opts ...option.ClientOption) (*youtube.Service, error)
    }

    // For youtube
    IYoutubeService interface {
        // Channel operations
        GetChannel(ctx context.Context, part []string, mine bool) (*youtube.ChannelListResponse, error)
        // Video operations
        ListVideos(ctx context.Context, part []string, channelId string, maxResults int64) (*youtube.SearchListResponse, error)
        UploadVideoBase(ctx context.Context, title, description, privacyStatus, videoPath string) (*youtube.Video, error)
        DeleteVideo(ctx context.Context, videoId string) error
        // Playlist operations
        CreatePlaylist(ctx context.Context, title, description string) (*youtube.Playlist, error)
        AddVideoToPlaylist(ctx context.Context, playlistId, videoId string) (*youtube.PlaylistItem, error)
    }
)

// #####################################################

// Error definitions
var (
    ErrGoogleServiceNotInitialized     = errors.New("google service is not initialized")
    ErrGoogleUserServiceNotInitialized = errors.New("google user service is not initialized")
    ErrYoutubeServiceNotInitialized    = errors.New("youtube service is not initialized")
)

// local var interface
var (
    vIGoogleService     IGoogleService
    vIGoogleUserService IGoogleUserService
    vIYoutubeService    IYoutubeService
)

// #####################################################

// initialize google services
func InitGoogleService(iGoogleService IGoogleService) {
    vIGoogleService = iGoogleService
}

// get google service instance
func GetGoogleService() (IGoogleService, error) {
    if vIGoogleService == nil {
        return nil, ErrGoogleServiceNotInitialized
    }
    return vIGoogleService, nil
}

// #####################################################

// initialize google user services
func InitGoogleUserService(iGoogleUserService IGoogleUserService) {
    vIGoogleUserService = iGoogleUserService
}

// get google user service instance
func GetGoogleUserService() (IGoogleUserService, error) {
    if vIGoogleUserService == nil {
        return nil, ErrGoogleUserServiceNotInitialized
    }
    return vIGoogleUserService, nil
}

// #####################################################

// initialize youtube services
func InitYoutubeService(iYoutubeService IYoutubeService) {
    vIYoutubeService = iYoutubeService
}

// get youtube service instance
func GetYoutubeService() (IYoutubeService, error) {
    if vIYoutubeService == nil {
        return nil, ErrYoutubeServiceNotInitialized
    }
    return vIYoutubeService, nil
}

// #####################################################
