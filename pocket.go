package pocket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"
)

const (
	baseURL          = "https://getpocket.com/"
	requestTokenPath = "/v3/oauth/request"
	accessTokenPath  = "/v3/oauth/authorize"
	authPath         = "/auth/authorize"
)

const (
	requestTokenKey        = "code"
	accessTokenKey         = "access_token"
	requestTokenQueryParam = "request_token"
	redirectUriQueryParam  = "redirect_uri"
)

type Pocket struct {
	consumerKey  string
	requestToken string
	accessToken  string
	baseURL      string
	httpClient   *http.Client
}

func New(consumerKey string) *Pocket {
	return &Pocket{
		consumerKey: consumerKey,
		httpClient: &http.Client{
			Timeout: time.Second * 5,
		},
		baseURL: baseURL,
	}
}

func (p *Pocket) WithBaseUrl(baseURL string) *Pocket {
	p.baseURL = baseURL
	return p
}

func (p *Pocket) WithHttpClient(client *http.Client) *Pocket {
	p.httpClient = client
	return p
}

func (p *Pocket) doRequest(ctx context.Context, pocketPath string, reqData interface{}) (map[string]string, error) {
	res := make(map[string]string)
	u, err := url.Parse(p.baseURL)
	if err != nil {
		return nil, fmt.Errorf("error while parsing base url: %w", err)
	}
	u.Path = path.Join(u.Path, pocketPath)

	body, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("error while marshalling request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error while building request: %w", err)
	}
	req.Header.Set("X-Accept", "application/json")
	req.Header.Set("Content-type", "application/json; charset=UTF8")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, NewErrorPocket(
			resp.Header.Get("X-Error"),
			resp.Header.Get("X-Error-Code"),
			resp.StatusCode,
		)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading data from response: %w", err)
	}

	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling http response body: %w", err)
	}

	return res, nil
}

type codeRequest struct {
	ConsumerKey string `json:"consumer_key"`
	RedirectUri string `json:"redirect_uri"`
}

func (p *Pocket) GenerateRequestToken(ctx context.Context, redirectURI string) (string, error) {
	resp, err := p.doRequest(ctx, requestTokenPath, &codeRequest{
		ConsumerKey: p.consumerKey,
		RedirectUri: redirectURI,
	})
	if err != nil {
		return "", err
	}
	return resp[requestTokenKey], nil
}

func (p *Pocket) GetRequestToken() string {
	return p.requestToken
}

func (p *Pocket) SetRequestToken(requestToken string) {
	p.requestToken = requestToken
}

func (p *Pocket) AuthApp(ctx context.Context, redirectURI string) error {
	var err error
	p.requestToken, err = p.GenerateRequestToken(ctx, redirectURI)
	if err != nil {
		return err
	}
	return nil
}

type accessTokenRequest struct {
	ConsumerKey string `json:"consumer_key"`
	Code        string `json:"code"`
}

func (p *Pocket) GenerateAccessToken(ctx context.Context) (string, error) {
	resp, err := p.doRequest(ctx, accessTokenPath, &accessTokenRequest{
		ConsumerKey: p.consumerKey,
		Code:        p.requestToken,
	})
	if err != nil {
		return "", err
	}
	return resp[accessTokenKey], nil
}

func (p *Pocket) GetAccessToken() string {
	return p.accessToken
}

func (p *Pocket) SetAccessToken(at string) {
	p.accessToken = at
}

func (p *Pocket) AuthUser(ctx context.Context) error {
	var err error
	p.accessToken, err = p.GenerateAccessToken(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (p *Pocket) BuildAuthUrl(redirectUri string) (string, error) {
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
