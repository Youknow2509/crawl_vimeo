package cmd

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net"
// 	"net/http"
// 	"os"
// 	"time"

// 	"golang.org/x/oauth2"
// 	"golang.org/x/oauth2/google"
// 	"google.golang.org/api/option"
// 	"google.golang.org/api/youtube/v3"
// )

// // scope use
// const (
// 	MANAGET_YOUTUBE_ACCOUNT      = "https://www.googleapis.com/auth/youtube"
// 	MANAGET_YOUTUBE_VIDEOS       = "https://www.googleapis.com/auth/youtube.force-ssl"
// 	MANAGET_VIEW_YOUTUBE_ACCOUNT = "https://www.googleapis.com/auth/youtube.readonly"
// 	MANAGET_VIDEO                = "https://www.googleapis.com/auth/youtube.upload"
// )

// const (
// 	CLIENT_SECRET_PATH = "secrets/client_secret.json"
// 	USER_AUTH_FILE     = "secrets/user_auth.json"
// 	REDIRECT_URI       = "http://localhost:8080"
// 	REDIRECT_HOST      = "localhost:8080"
// )

// // get all scopes
// func GetScopes() []string {
// 	return []string{
// 		MANAGET_YOUTUBE_ACCOUNT,
// 		MANAGET_YOUTUBE_VIDEOS,
// 		MANAGET_VIEW_YOUTUBE_ACCOUNT,
// 		MANAGET_VIDEO,
// 	}
// }

// // Đọc token từ file
// func tokenFromFile(file string) (*oauth2.Token, error) {
// 	b, err := ioutil.ReadFile(file)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var token oauth2.Token
// 	err = json.Unmarshal(b, &token)
// 	return &token, err
// }

// // Lưu token vào file
// func saveToken(path string, token *oauth2.Token) error {
// 	f, err := os.Create(path)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()
// 	return json.NewEncoder(f).Encode(token)
// }

// // Lấy token qua server local (browser login)
// func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
// 	listener, err := net.Listen("tcp", REDIRECT_HOST)
// 	if err != nil {
// 		log.Fatalf("Unable to listen on %s: %v", REDIRECT_HOST, err)
// 	}
// 	defer listener.Close()
// 	config.RedirectURL = REDIRECT_URI

// 	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
// 	fmt.Printf("Hãy mở trình duyệt và truy cập:\n%s\n", authURL)

// 	codeCh := make(chan string)
// 	// Server HTTP xử lý mã code trả về
// 	go func() {
// 		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 			code := r.URL.Query().Get("code")
// 			if code != "" {
// 				fmt.Fprintf(w, "<h1>Xác thực thành công! Bạn có thể đóng tab này.</h1>")
// 				codeCh <- code
// 			} else {
// 				http.Error(w, "Không tìm thấy code xác thực", http.StatusBadRequest)
// 			}
// 		})
// 		http.Serve(listener, nil)
// 	}()
// 	code := <-codeCh
// 	tok, err := config.Exchange(context.Background(), code)
// 	if err != nil {
// 		log.Fatalf("Không thể đổi code lấy token: %v", err)
// 	}
// 	return tok
// }

// // Lấy http.Client đã xác thực OAuth2, kiểm tra hạn token, refresh nếu hết hạn
// func getClient(ctx context.Context, config *oauth2.Config, tokenPath string) (*http.Client, *oauth2.Token) {
// 	token, err := getValidToken(ctx, config, tokenPath)
// 	if err != nil {
// 		// Nếu không có token hợp lệ thì bắt user xác thực lại
// 		token = getTokenFromWeb(config)
// 		if err := saveToken(tokenPath, token); err != nil {
// 			log.Fatalf("Không thể lưu token: %v", err)
// 		}
// 	}
// 	return config.Client(ctx, token), token
// }

// // valid token
// func getValidToken(ctx context.Context, config *oauth2.Config, tokenPath string) (*oauth2.Token, error) {
// 	token, err := tokenFromFile(tokenPath)
// 	if err != nil {
// 		return nil, fmt.Errorf("token file not found or error: %v", err)
// 	}
// 	// Nếu token hết hạn, thử refresh
// 	if !token.Valid() {
// 		ts := config.TokenSource(ctx, token)
// 		newToken, err := ts.Token()
// 		if err != nil || !newToken.Valid() {
// 			return nil, fmt.Errorf("token expired and refresh fail")
// 		}
// 		if newToken.AccessToken != token.AccessToken {
// 			_ = saveToken(tokenPath, newToken)
// 		}
// 		return newToken, nil
// 	}
// 	return token, nil
// }

// // Tạo Youtube client với xác thực OAuth2
// func createClient(ctx context.Context) (*youtube.Service, *oauth2.Token, error) {
// 	b, err := ioutil.ReadFile(CLIENT_SECRET_PATH)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("Không thể đọc file client secret: %v", err)
// 	}
// 	config, err := google.ConfigFromJSON(b, GetScopes()...)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("Không thể parse client secret: %v", err)
// 	}
// 	client, token := getClient(ctx, config, USER_AUTH_FILE)
// 	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
// 	return service, token, err
// }

// // Main function
// func CmdYtb() {
// 	ctx := context.Background()
// 	ytService, token, err := createClient(ctx)
// 	if err != nil {
// 		log.Fatalf("Failed to create YouTube client: %v", err)
// 	}
// 	fmt.Println("AccessToken:", token.AccessToken)
// 	fmt.Println("Expiry:", token.Expiry.Format(time.RFC3339))
// 	// Ping YouTube API
// 	resp, err := ytService.Channels.List([]string{"snippet"}).Mine(true).Do()
// 	if err != nil {
// 		log.Fatalf("Failed to ping YouTube API: %v", err)
// 	}
// 	for _, item := range resp.Items {
// 		fmt.Printf("Kênh của bạn: %s (ID: %s)\n", item.Snippet.Title, item.Id)
// 	}
// }
