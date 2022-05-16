package pocket

import (
	"context"
	"fmt"
	"net/url"
	"path"
)

const (
	requestTokenQueryParam = "request_token"
	redirectUriQueryParam  = "redirect_uri"
)

type codeRequest struct {
	ConsumerKey string `json:"consumer_key"`
	RedirectUri string `json:"redirect_uri"`
}

type AuthUserResponse struct {
	AccessToken string `json:"access_token"`
	Username    string `json:"username"`
}

type AuthAppResponse struct {
	Code string `json:"code"`
}

func (p *Pocket) GenerateRequestToken(ctx context.Context, redirectURI string) (*AuthAppResponse, error) {
	res := AuthAppResponse{}
	err := p.doRequest(ctx, requestTokenPath, &codeRequest{
		ConsumerKey: p.consumerKey,
		RedirectUri: redirectURI,
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (p *Pocket) GetRequestToken() string {
	return p.requestToken
}

func (p *Pocket) SetRequestToken(requestToken string) {
	p.requestToken = requestToken
}

func (p *Pocket) AuthApp(ctx context.Context, redirectURI string) error {
	var err error
	resp, err := p.GenerateRequestToken(ctx, redirectURI)
	if err != nil {
		return err
	}
	p.SetRequestToken(resp.Code)
	return nil
}

type accessTokenRequest struct {
	ConsumerKey string `json:"consumer_key"`
	Code        string `json:"code"`
}

func (p *Pocket) GenerateAccessToken(ctx context.Context) (*AuthUserResponse, error) {
	res := AuthUserResponse{}
	err := p.doRequest(ctx, accessTokenPath, &accessTokenRequest{
		ConsumerKey: p.consumerKey,
		Code:        p.requestToken,
	}, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (p *Pocket) GetAccessToken() string {
	return p.accessToken
}

func (p *Pocket) SetAccessToken(at string) {
	p.accessToken = at
}

func (p *Pocket) AuthUser(ctx context.Context) error {
	var err error
	resp, err := p.GenerateAccessToken(ctx)
	if err != nil {
		return err
	}
	p.SetAccessToken(resp.AccessToken)
	return nil
}

func (p *Pocket) MakeAuthUrl(redirectUri string) (string, error) {
	u, err := url.Parse(p.baseURL)
	if err != nil {
		return "", fmt.Errorf("error while building auth url: %w", err)
	}
	u.Path = path.Join(u.Path, authPath)
	q := u.Query()
	q.Add(requestTokenQueryParam, p.requestToken)
	q.Add(redirectUriQueryParam, redirectUri)
	u.RawQuery = q.Encode()

	return u.String(), nil
}
