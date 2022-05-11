package pocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	redirectURL  = "google.com"
	consumerKey  = "consumer-key"
	requestToken = "request-token"
	accessToken  = "access-token"
)

// messages
const (
	msgMissConsumerKey    = "Missing consumer key"
	msgMissRedirectUrl    = "Missing redirect url."
	msgInvalidConsumerKey = "Invalid consumer key."
	msgPocketServerIssue  = "Pocket server issue."
	msgInvalidRedirectUri = "Invalid redirect uri."
	msgMissingCode        = "Missing code."
	msgCodeNotFound       = "Code not found."
	msgRejectedCode       = "User rejected code."
	msgCodeAlreadyUsed    = "Already used code."
)

// x error codes
const (
	xMissConsumerKey    = "138"
	xMissRedirectUrl    = "140"
	xInvalidConsumerKey = "152"
	xPocketServerIssue  = "199"
	xInvalidRedirectUri = "181"
	xMissingCode        = "182"
	xCodeNotFound       = "185"
	xRejectedCode       = "158"
	xCodeAlreadyUsed    = "159"
)

func TestPocket_AuthApp(t *testing.T) {
	tests := []struct {
		name        string
		consumerKey string
		redirectUrl string
		expErr      *ErrorPocket
		expToken    string
		handler     func(t *testing.T) http.HandlerFunc
	}{
		{
			name:        "Missing consumer key",
			redirectUrl: redirectURL,
			expErr: &ErrorPocket{
				Message:  msgMissConsumerKey,
				Xcode:    xMissConsumerKey,
				HttpCode: http.StatusBadRequest,
			},
			handler: func(t *testing.T) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					require.Equal(t, requestTokenPath, r.URL.Path)

					w.Header().Add("X-Error-Code", xMissConsumerKey)
					w.Header().Add("X-Error", msgMissConsumerKey)
					w.WriteHeader(http.StatusBadRequest)
				}
			},
		},
		{
			name:        "Missing redirect url",
			consumerKey: consumerKey,
			redirectUrl: "",
			expErr: &ErrorPocket{
				Message:  msgMissRedirectUrl,
				Xcode:    xMissRedirectUrl,
				HttpCode: http.StatusBadRequest,
			},
			handler: func(t *testing.T) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					require.Equal(t, requestTokenPath, r.URL.Path)

					w.Header().Add("X-Error-Code", xMissRedirectUrl)
					w.Header().Add("X-Error", msgMissRedirectUrl)
					w.WriteHeader(http.StatusBadRequest)
				}
			},
		},
		{
			name:        "Invalid consumer key",
			consumerKey: consumerKey + "invalid",
			redirectUrl: redirectURL,
			expErr: &ErrorPocket{
				Message:  msgInvalidConsumerKey,
				Xcode:    xInvalidConsumerKey,
				HttpCode: http.StatusForbidden,
			},
			handler: func(t *testing.T) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					require.Equal(t, requestTokenPath, r.URL.Path)

					w.Header().Add("X-Error-Code", xInvalidConsumerKey)
					w.Header().Add("X-Error", msgInvalidConsumerKey)
					w.WriteHeader(http.StatusForbidden)
				}
			},
		},
		{
			name:        "Pocket server issue",
			consumerKey: consumerKey,
			redirectUrl: redirectURL,
			expErr: &ErrorPocket{
				Message:  msgPocketServerIssue,
				Xcode:    xPocketServerIssue,
				HttpCode: http.StatusInternalServerError,
			},
			handler: func(t *testing.T) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					require.Equal(t, requestTokenPath, r.URL.Path)

					w.Header().Add("X-Error-Code", xPocketServerIssue)
					w.Header().Add("X-Error", msgPocketServerIssue)
					w.WriteHeader(http.StatusInternalServerError)
				}
			},
		},
		{
			name:        "Success",
			consumerKey: consumerKey,
			redirectUrl: redirectURL,
			expToken:    requestToken,
			handler: func(t *testing.T) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					data, err := io.ReadAll(r.Body)
					require.NoError(t, err)

					req := codeRequest{}
					require.NoError(t, json.Unmarshal(data, &req))

					require.Equal(t, consumerKey, req.ConsumerKey)
					require.Equal(t, redirectURL, req.RedirectUri)
					require.Equal(t, requestTokenPath, r.URL.Path)

					resp := map[string]string{
						requestTokenKey: requestToken,
					}
					data, err = json.Marshal(resp)
					require.NoError(t, err)
					w.Header().Add("Content-type", "application/json")
					w.Write(data)
					w.WriteHeader(http.StatusOK)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(tc.handler(t))
			defer srv.Close()

			p := New(tc.consumerKey).WithBaseUrl(srv.URL)

			err := p.AuthApp(context.Background(), tc.redirectUrl)

			if tc.expErr != nil {
				var perr *ErrorPocket
				if errors.As(err, &perr) {
					require.Equal(t, tc.expErr.Message, perr.Message)
					require.Equal(t, tc.expErr.Xcode, perr.Xcode)
					require.Equal(t, tc.expErr.HttpCode, perr.HttpCode)
				} else {
					require.Fail(t, "unknown error", err)
				}
			} else {
				require.Equal(t, tc.expToken, p.GetRequestToken())
			}
		})
	}
}

func TestPocket_AuthUser(t *testing.T) {
	tests := []struct {
		name        string
		consumerKey string
		expErr      *ErrorPocket
		expToken    string
		handler     func(t *testing.T) http.HandlerFunc
	}{
		{
			name: "Missing consumer key",
			expErr: &ErrorPocket{
				Message:  msgMissConsumerKey,
				Xcode:    xMissConsumerKey,
				HttpCode: http.StatusBadRequest,
			},
			handler: func(t *testing.T) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					require.Equal(t, accessTokenPath, r.URL.Path)

					w.Header().Add("X-Error-Code", xMissConsumerKey)
					w.Header().Add("X-Error", msgMissConsumerKey)
					w.WriteHeader(http.StatusBadRequest)
				}
			},
		},
		{
			name: "Invalid consumer key",
			expErr: &ErrorPocket{
				Message:  msgInvalidConsumerKey,
				Xcode:    xInvalidConsumerKey,
				HttpCode: http.StatusForbidden,
			},
			handler: func(t *testing.T) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					require.Equal(t, accessTokenPath, r.URL.Path)

					w.Header().Add("X-Error-Code", xInvalidConsumerKey)
					w.Header().Add("X-Error", msgInvalidConsumerKey)
					w.WriteHeader(http.StatusForbidden)
				}
			},
		},
		{
			name: "Invalid redirect uri",
			expErr: &ErrorPocket{
				Message:  msgInvalidRedirectUri,
				Xcode:    xInvalidRedirectUri,
				HttpCode: http.StatusBadRequest,
			},
			handler: func(t *testing.T) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					require.Equal(t, accessTokenPath, r.URL.Path)

					w.Header().Add("X-Error-Code", xInvalidRedirectUri)
					w.Header().Add("X-Error", msgInvalidRedirectUri)
					w.WriteHeader(http.StatusBadRequest)
				}
			},
		},
		{
			name: "Missing code",
			expErr: &ErrorPocket{
				Message:  msgMissingCode,
				Xcode:    xMissingCode,
				HttpCode: http.StatusBadRequest,
			},
			handler: func(t *testing.T) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					require.Equal(t, accessTokenPath, r.URL.Path)

					w.Header().Add("X-Error-Code", xMissingCode)
					w.Header().Add("X-Error", msgMissingCode)
					w.WriteHeader(http.StatusBadRequest)
				}
			},
		},
		{
			name: "Code not found",
			expErr: &ErrorPocket{
				Message:  msgCodeNotFound,
				Xcode:    xCodeNotFound,
				HttpCode: http.StatusBadRequest,
			},
			handler: func(t *testing.T) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					require.Equal(t, accessTokenPath, r.URL.Path)

					w.Header().Add("X-Error-Code", xCodeNotFound)
					w.Header().Add("X-Error", msgCodeNotFound)
					w.WriteHeader(http.StatusBadRequest)
				}
			},
		},
		{
			name: "User rejected code",
			expErr: &ErrorPocket{
				Message:  msgRejectedCode,
				Xcode:    xRejectedCode,
				HttpCode: http.StatusForbidden,
			},
			handler: func(t *testing.T) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					require.Equal(t, accessTokenPath, r.URL.Path)

					w.Header().Add("X-Error-Code", xRejectedCode)
					w.Header().Add("X-Error", msgRejectedCode)
					w.WriteHeader(http.StatusForbidden)
				}
			},
		},
		{
			name: "Already used code",
			expErr: &ErrorPocket{
				Message:  msgCodeAlreadyUsed,
				Xcode:    xCodeAlreadyUsed,
				HttpCode: http.StatusForbidden,
			},
			handler: func(t *testing.T) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					require.Equal(t, accessTokenPath, r.URL.Path)

					w.Header().Add("X-Error-Code", xCodeAlreadyUsed)
					w.Header().Add("X-Error", msgCodeAlreadyUsed)
					w.WriteHeader(http.StatusForbidden)
				}
			},
		},
		{
			name:        "Success",
			consumerKey: consumerKey,
			expToken:    accessToken,
			handler: func(t *testing.T) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					data, err := io.ReadAll(r.Body)
					require.NoError(t, err)

					req := accessTokenRequest{}
					require.NoError(t, json.Unmarshal(data, &req))

					require.Equal(t, consumerKey, req.ConsumerKey)
					require.Equal(t, requestToken, req.Code)
					require.Equal(t, accessTokenPath, r.URL.Path)

					resp := map[string]string{
						accessTokenKey: accessToken,
					}
					data, err = json.Marshal(resp)
					require.NoError(t, err)
					w.Header().Add("Content-type", "application/json")
					w.Write(data)
					w.WriteHeader(http.StatusOK)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(tc.handler(t))
			defer srv.Close()

			p := New(tc.consumerKey).WithBaseUrl(srv.URL)

			p.SetRequestToken(requestToken)
			err := p.AuthUser(context.Background())

			if tc.expErr != nil {
				var perr *ErrorPocket
				if errors.As(err, &perr) {
					require.Equal(t, tc.expErr.Message, perr.Message)
					require.Equal(t, tc.expErr.Xcode, perr.Xcode)
					require.Equal(t, tc.expErr.HttpCode, perr.HttpCode)
				} else {
					require.Fail(t, "unknown error", err)
				}
			} else {
				require.Equal(t, tc.expToken, p.GetAccessToken())
			}
		})
	}
}

func TestPocket_BuildAuthUrl(t *testing.T) {
	p := &Pocket{
		consumerKey:  consumerKey,
		requestToken: requestToken,
		baseURL:      baseURL,
	}

	u, err := url.Parse(p.baseURL)
	require.NoError(t, err)
	u.Path = path.Join(u.Path, authPath)
	u.RawQuery = fmt.Sprintf(
		"%s=%s&%s=%s",
		redirectUriQueryParam, redirectURL,
		requestTokenQueryParam, requestToken)

	bu, err := p.BuildAuthUrl(redirectURL)
	require.NoError(t, err)
	require.Equal(t, u.String(), bu)
}
