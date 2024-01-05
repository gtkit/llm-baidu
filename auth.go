package baidu

import (
	"context"
	"errors"
	"net/http"
	"time"
)

const authAPI = "https://aip.baidubce.com/oauth/2.0/token"

type AuthRequest struct {
	GrantType    string `json:"grant_type"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type AuthResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int64  `json:"expires_in"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type AuthToken struct {
	token     string
	expiresAt int64 // 到期时间
	expiresIn int64 // 有效期
}

func (t *AuthToken) IsValid() bool {
	if t.token == "" {
		return false
	}
	currentTime := time.Now().Unix()
	return currentTime >= t.expiresAt-5
}

// CreateAccessToken — API call to Create a completion for the chat message.
func (c *Client) CreateAccessToken(ctx context.Context) (response AuthResponse, err error) {
	return c.RequestAccessToken(ctx, AuthRequest{
		GrantType:    c.config.GrantType,
		ClientId:     c.config.ClientId,
		ClientSecret: c.config.ClientSecret,
	})
}

// RequestAccessToken — API call to Create a completion for the chat message.
func (c *Client) RequestAccessToken(ctx context.Context, request AuthRequest) (response AuthResponse, err error) {
	api := authAPI
	if c.config.AuthAPI != "" {
		api = c.config.AuthAPI
	}

	query := withQuery(map[string]string{
		"client_id":     request.ClientId,
		"client_secret": request.ClientSecret,
		"grant_type":    request.GrantType,
	})
	ct := withContentType("application/json")

	req, err := c.newRequest(ctx, http.MethodPost, api, query, ct)
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}

// AutoHandleAccessToken 获取access_token.
func (c *Client) AutoHandleAccessToken(ctx context.Context) error {
	if c.authToken.IsValid() {
		return nil
	}

	// 生成新的token.
	resp, err := c.CreateAccessToken(ctx)
	if err != nil {
		return errors.New("create access token error: " + err.Error())
	}

	c.authToken = AuthToken{
		token:     resp.AccessToken,
		expiresAt: time.Now().Unix() + resp.ExpiresIn,
		expiresIn: resp.ExpiresIn,
	}

	return nil
}
