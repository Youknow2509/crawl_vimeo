package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/youknow2509/crawl_vimeo/consts"
	"golang.org/x/oauth2"
)

// get all scopes
func GetScopesYtb() []string {
	return []string{
		consts.YTB_SCOPE_MANAGET_YOUTUBE_ACCOUNT,
		consts.YTB_SCOPE_MANAGET_YOUTUBE_VIDEOS,
		consts.YTB_SCOPE_MANAGET_VIEW_YOUTUBE_ACCOUNT,
		consts.YTB_SCOPE_MANAGET_VIDEO,
	}
}

// Đọc token ytb user từ file
func tokenYtbUserFromFile(file string) (*oauth2.Token, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var token oauth2.Token
	err = json.Unmarshal(b, &token)
	return &token, err
}

// Lưu token user sau khi auth vào file
func saveTokenYtbUserAuth(path string, token *oauth2.Token) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}

// Lấy token user sau khi auth từ web
func getTokenUserAuthYtbFromWeb(config *oauth2.Config) *oauth2.Token {
	listener, err := net.Listen("tcp", consts.YTB_REDIRECT_HOST)
	if err != nil {
		log.Fatalf("Unable to listen on %s: %v", consts.YTB_REDIRECT_HOST, err)
	}
	defer listener.Close()
	config.RedirectURL = consts.YTB_REDIRECT_URI

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	fmt.Printf("Hãy mở trình duyệt và truy cập:\n%s\n", authURL)

	codeCh := make(chan string)
	// Server HTTP xử lý mã code trả về
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")
			if code != "" {
				fmt.Fprintf(w, "<h1>Xác thực thành công! Bạn có thể đóng tab này.</h1>")
				codeCh <- code
			} else {
				http.Error(w, "Không tìm thấy code xác thực", http.StatusBadRequest)
			}
		})
		http.Serve(listener, nil)
	}()
	code := <-codeCh
	tok, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Không thể đổi code lấy token: %v", err)
	}
	return tok
}

// valid token user
func GetValidTokenYtbUser(ctx context.Context, config *oauth2.Config, tokenPath string) (*oauth2.Token, error) {
	token, err := tokenYtbUserFromFile(tokenPath)
	if err != nil {
		return nil, fmt.Errorf("token file not found or error: %v", err)
	}
	// Nếu token hết hạn, thử refresh
	if !token.Valid() {
		ts := config.TokenSource(ctx, token)
		newToken, err := ts.Token()
		if err != nil || !newToken.Valid() {
			return nil, fmt.Errorf("token expired and refresh fail")
		}
		if newToken.AccessToken != token.AccessToken {
			_ = saveTokenYtbUserAuth(tokenPath, newToken)
		}
		return newToken, nil
	}
	return token, nil
}

// Lấy http.Client đã xác thực OAuth2, kiểm tra hạn token, refresh nếu hết hạn
func GetClientYtb(ctx context.Context, config *oauth2.Config, tokenPath string) (*http.Client, *oauth2.Token) {
	token, err := GetValidTokenYtbUser(ctx, config, tokenPath)
	if err != nil {
		// Nếu không có token hợp lệ thì bắt user xác thực lại
		token = getTokenUserAuthYtbFromWeb(config)
		if err := saveTokenYtbUserAuth(tokenPath, token); err != nil {
			log.Fatalf("Không thể lưu token: %v", err)
		}
	}
	return config.Client(ctx, token), token
}