package services

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/youknow2509/crawl_vimeo/consts"
	"github.com/youknow2509/crawl_vimeo/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// Create youtobe client
func CreateYtbClient(ctx context.Context) (*youtube.Service, *oauth2.Token, error) {
	b, err := ioutil.ReadFile(consts.YTB_CLIENT_SECRET_PATH)
	if err != nil {
		return nil, nil, fmt.Errorf("Không thể đọc file client secret: %v", err)
	}
	config, err := google.ConfigFromJSON(b, utils.GetScopesYtb()...)
	if err != nil {
		return nil, nil, fmt.Errorf("Không thể parse client secret: %v", err)
	}
	client, token := utils.GetClientYtb(ctx, config, consts.YTB_USER_AUTH_FILE)
	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	return service, token, err
}