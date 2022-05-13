package pocket

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	testUrl = "https://vk.com"

	urlKey    = "url"
	actionKey = "action"
)

func TestPocket_Modify(t *testing.T) {
	tests := []struct {
		name        string
		consumerKey string
		expErr      *ErrorPocket
		actions     Actions
		handler     func(t *testing.T) http.HandlerFunc
	}{
		{
			name: "Unauthorized",
			expErr: &ErrorPocket{
				Message:  msgUnauthorized,
				Xcode:    xUnauthorized,
				HttpCode: http.StatusUnauthorized,
			},
			actions: Actions{
				&ActionAdd{
					Action: ActionAddType,
					Url:    testUrl,
				},
			},
			handler: func(t *testing.T) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					require.Equal(t, modifyPath, r.URL.Path)

					w.Header().Add("X-Error-Code", xUnauthorized)
					w.Header().Add("X-Error", msgUnauthorized)
					w.WriteHeader(http.StatusUnauthorized)
				}
			},
		},
		{
			name: "Success[action=add]",
			actions: Actions{
				&ActionAdd{
					Action: ActionAddType,
					Url:    testUrl,
				},
			},
			handler: func(t *testing.T) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					data, err := io.ReadAll(r.Body)
					require.NoError(t, err)

					req := modifyRequest{}
					require.NoError(t, json.Unmarshal(data, &req))

					require.Equal(t, consumerKey, req.ConsumerKey)
					require.Equal(t, accessToken, req.AccessToken)
					addReq, ok := req.Actions[0].(map[string]interface{})

					if ok {
						require.Equal(t, testUrl, addReq[urlKey])
						require.Equal(t, ActionAddType, addReq[actionKey])
					}

					resp := &ModifyResponse{
						Status: successStatus,
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

			p := &Pocket{
				consumerKey:  consumerKey,
				requestToken: requestToken,
				accessToken:  accessToken,
				baseURL:      srv.URL,
				httpClient:   &http.Client{Timeout: 5 * time.Second},
			}

			res, err := p.Modify(context.Background(), tc.actions)

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
				require.Equal(t, successStatus, res.Status)
			}
		})
	}
}
