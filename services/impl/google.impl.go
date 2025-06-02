package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/youknow2509/crawl_vimeo/consts"
	"github.com/youknow2509/crawl_vimeo/services"
	"github.com/youknow2509/crawl_vimeo/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// #################################
// struct implement IGoogleService
type GoogleService struct{}

// GetConfigServer implements services.IGoogleService.
func (g *GoogleService) GetConfigServer(file string) (*oauth2.Config, error) {
    if file == "" {
        return nil, fmt.Errorf("file path cannot be empty")
    }
    
    // Read the client secret file
    b, err := os.ReadFile(file)
    if err != nil {
        return nil, fmt.Errorf("không thể đọc file client secret: %w", err)
    }
    
    // Create config from JSON
    config, err := google.ConfigFromJSON(b, g.GetScopesServer()...)
    if err != nil {
        return nil, fmt.Errorf("không thể parse client secret: %w", err)
    }
    
    return config, nil
}

// GetScopesServer implements services.IGoogleService.
func (g *GoogleService) GetScopesServer() []string {
    return []string{
        consts.YTB_SCOPE_MANAGET_YOUTUBE_ACCOUNT,
        consts.YTB_SCOPE_MANAGET_YOUTUBE_VIDEOS,
        consts.YTB_SCOPE_MANAGET_VIEW_YOUTUBE_ACCOUNT,
        consts.YTB_SCOPE_MANAGET_VIDEO,
		// v.v
    }
}

// #################################
// struct implement IGoogleUserService
type GoogleUserService struct{}

// GetUrlRedirectAuth implements services.IGoogleUserService.
func (g *GoogleUserService) GetUrlRedirectAuth(ctx context.Context, config *oauth2.Config) (string, error) {
    if config == nil {
        return "", fmt.Errorf("config cannot be nil")
    }
    
    options := []oauth2.AuthCodeOption{
        oauth2.SetAuthURLParam("access_type", "offline"),
        oauth2.SetAuthURLParam("approval_prompt", "force"),
		// v.v
    }

    return config.AuthCodeURL("state-token", options...), nil
}

// GetYoutubeService implements services.IGoogleUserService.
func (g *GoogleUserService) GetYoutubeService(ctx context.Context, opts ...option.ClientOption) (*youtube.Service, error) {
    return youtube.NewService(ctx, opts...)
}

// GetTokenFromFile implements services.IGoogleUserService.
func (g *GoogleUserService) GetTokenFromFile(file string) (*oauth2.Token, error) {
    if file == "" {
        return nil, fmt.Errorf("file path cannot be empty")
    }
    
    b, err := os.ReadFile(file)
    if err != nil {
        return nil, fmt.Errorf("không thể đọc file token: %w", err)
    }
    
    var token oauth2.Token
    err = json.Unmarshal(b, &token)
    if err != nil {
        return nil, fmt.Errorf("không thể parse token: %w", err)
    }
    
    return &token, nil
}

// GetTokenFromWeb implements services.IGoogleUserService.
func (g *GoogleUserService) GetTokenFromWeb(ctx context.Context, config *oauth2.Config, fileSave string) (*oauth2.Token, error) {
    if config == nil {
        return nil, fmt.Errorf("config cannot be nil")
    }
    if fileSave == "" {
        return nil, fmt.Errorf("file save path cannot be empty")
    }
    
    // Create context with timeout
    ctx, cancel := context.WithTimeout(ctx, consts.YTB_USER_TIME_EXPIRED_AUTH_REDIRECT*time.Minute)
    defer cancel()
    
    // open tcp listener
    listener, err := net.Listen("tcp", consts.YTB_REDIRECT_HOST)
    if err != nil {
        return nil, fmt.Errorf("unable to listen on %s: %w", consts.YTB_REDIRECT_HOST, err)
    }
    defer listener.Close()
    
    // get redirect URL
    urlRedirect, err := g.GetUrlRedirectAuth(ctx, config)
    if err != nil {
        return nil, fmt.Errorf("không thể lấy URL redirect: %w", err)
    }
    
    fmt.Printf("Hãy mở trình duyệt và truy cập:\n%s\n", urlRedirect)
    
    codeCh := make(chan string, 1)
    errCh := make(chan error, 1)
    
    // Server HTTP xử lý mã code trả về
    server := &http.Server{
        Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            code := r.URL.Query().Get("code")
            if code != "" {
                fmt.Fprintf(w, "<h1>Xác thực thành công! Bạn có thể đóng tab này.</h1>")
                select {
                case codeCh <- code:
                default:
                }
            } else {
                http.Error(w, "Không tìm thấy code xác thực", http.StatusBadRequest)
                select {
                case errCh <- fmt.Errorf("không tìm thấy code xác thực"):
                default:
                }
            }
        }),
    }
    
    go func() {
        if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
            select {
            case errCh <- err:
            default:
            }
        }
    }()
    
    var code string
    select {
    case code = <-codeCh:
    case err := <-errCh:
        return nil, err
    case <-ctx.Done():
        return nil, fmt.Errorf("timeout waiting for authorization code")
    }
    
    // Shutdown server
    server.Shutdown(context.Background())
    
    tok, err := config.Exchange(ctx, code)
    if err != nil {
        return nil, fmt.Errorf("không thể đổi code lấy token: %w", err)
    }
    
    if err := g.SaveToken(fileSave, tok); err != nil {
        return nil, fmt.Errorf("không thể lưu token: %w", err)
    }
    
    return tok, nil
}

// RefreshToken implements services.IGoogleUserService.
func (g *GoogleUserService) RefreshToken(ctx context.Context, config *oauth2.Config, token *oauth2.Token) (*oauth2.Token, error) {
    if config == nil {
        return nil, fmt.Errorf("config cannot be nil")
    }
    if token == nil {
        return nil, fmt.Errorf("token cannot be nil")
    }
    
    tokenSource := config.TokenSource(ctx, token)
    newToken, err := tokenSource.Token()
    if err != nil {
        return nil, fmt.Errorf("không thể refresh token: %w", err)
    }
    
    return newToken, nil
}

// SaveToken implements services.IGoogleUserService.
func (g *GoogleUserService) SaveToken(path string, token *oauth2.Token) error {
    if path == "" {
        return fmt.Errorf("path cannot be empty")
    }
    if token == nil {
        return fmt.Errorf("token cannot be nil")
    }
    
    f, err := os.Create(path)
    if err != nil {
        return fmt.Errorf("không thể tạo file: %w", err)
    }
    defer f.Close()
    
    return json.NewEncoder(f).Encode(token)
}

// ValidateToken implements services.IGoogleUserService.
func (g *GoogleUserService) ValidateToken(ctx context.Context, token *oauth2.Token) (int, bool) {
    // Check if token is nil
    if token == nil {
        return consts.YTB_USER_TOKEN_NOT_FOUND, false
    }
    
    // Check if token is expired
    if token.Expiry.Before(time.Now()) {
        return consts.YTB_USER_TOKEN_EXPIRED, false
    }
    
    // Check if token is valid (has access token and not expired)
    if !token.Valid() {
        return consts.YTB_USER_TOKEN_INVALIDATED, false
    }
    
    // Additional check: verify token with Google API (optional)
    if err := utils.ValidateTokenWithGoogle(ctx, token); err != nil {
        return consts.YTB_USER_TOKEN_NOT_AUTHORIZED, false
    }
    
    return consts.YTB_USER_TOKEN_VALIDATE, true
}

// #################################
// struct implement IYoutubeService
type YoutubeService struct {
    Client *youtube.Service
}

// GetChannel implements services.IYoutubeService.
func (y *YoutubeService) GetChannel(ctx context.Context, part []string, mine bool) (*youtube.ChannelListResponse, error) {
    call := y.Client.Channels.List(part).Mine(mine)
    return call.Do()
}

// ListVideos implements services.IYoutubeService.
func (y *YoutubeService) ListVideos(ctx context.Context, part []string, channelId string, maxResults int64) (*youtube.SearchListResponse, error) {
    call := y.Client.Search.List(part).
        ChannelId(channelId).
        MaxResults(maxResults).
        Type("video")
    return call.Do()
}

// UploadVideo implements services.IYoutubeService.
func (y *YoutubeService) UploadVideo(ctx context.Context, title, description, filename string) (*youtube.Video, error) {
    // Implementation for video upload
	// TODO
    return nil, fmt.Errorf("not implemented yet")
}

// UpdateVideo implements services.IYoutubeService.
func (y *YoutubeService) UpdateVideo(ctx context.Context, videoId, title, description string) (*youtube.Video, error) {
    video := &youtube.Video{
        Id: videoId,
        Snippet: &youtube.VideoSnippet{
            Title:       title,
            Description: description,
        },
    }
    call := y.Client.Videos.Update([]string{"snippet"}, video)
    return call.Do()
}

// DeleteVideo implements services.IYoutubeService.
func (y *YoutubeService) DeleteVideo(ctx context.Context, videoId string) error {
    call := y.Client.Videos.Delete(videoId)
    return call.Do()
}

// CreatePlaylist implements services.IYoutubeService.
func (y *YoutubeService) CreatePlaylist(ctx context.Context, title, description string) (*youtube.Playlist, error) {
    playlist := &youtube.Playlist{
        Snippet: &youtube.PlaylistSnippet{
            Title:       title,
            Description: description,
        },
        Status: &youtube.PlaylistStatus{
            PrivacyStatus: "private",
        },
    }
    call := y.Client.Playlists.Insert([]string{"snippet", "status"}, playlist)
    return call.Do()
}

// AddVideoToPlaylist implements services.IYoutubeService.
func (y *YoutubeService) AddVideoToPlaylist(ctx context.Context, playlistId, videoId string) (*youtube.PlaylistItem, error) {
    playlistItem := &youtube.PlaylistItem{
        Snippet: &youtube.PlaylistItemSnippet{
            PlaylistId: playlistId,
            ResourceId: &youtube.ResourceId{
                Kind:    "youtube#video",
                VideoId: videoId,
            },
        },
    }
    call := y.Client.PlaylistItems.Insert([]string{"snippet"}, playlistItem)
    return call.Do()
}

// #################################
// implement interface
func NewGoogleService() services.IGoogleService {
    return &GoogleService{}
}

func NewGoogleUserService() services.IGoogleUserService {
    return &GoogleUserService{}
}

func NewYoutubeService(client *youtube.Service) services.IYoutubeService {
    return &YoutubeService{
        Client: client,
    }
}