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
	addPath          = "/v3/add"
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

func (p *Pocket) doRequestRaw(ctx context.Context, pocketPath string, reqData interface{}) ([]byte, error) {
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

	return data, nil
}

func (p *Pocket) doRequest(ctx context.Context, pocketPath string, reqData interface{}) (map[string]string, error) {
	res := make(map[string]string)
	data, err := p.doRequestRaw(ctx, pocketPath, reqData)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling http response body: %w", err)
	}

	return res, nil
}
