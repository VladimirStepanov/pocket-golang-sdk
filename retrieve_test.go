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

func TestPocket_Retrieve(t *testing.T) {
	tests := []struct {
		name        string
		consumerKey string
		expErr      *ErrorPocket
		ad          *RetrieveInput
		handler     func(t *testing.T) http.HandlerFunc
	}{
		{
			name: "Unauthorized",
			expErr: &ErrorPocket{
				Message:  msgUnauthorized,
				Xcode:    xUnauthorized,
				HttpCode: http.StatusUnauthorized,
			},
			ad: &RetrieveInput{
				State: Unread,
			},
			handler: func(t *testing.T) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					require.Equal(t, retrievePath, r.URL.Path)

					w.Header().Add("X-Error-Code", xUnauthorized)
					w.Header().Add("X-Error", msgUnauthorized)
					w.WriteHeader(http.StatusUnauthorized)
				}
			},
		},
		{
			name: "Success",
			ad: &RetrieveInput{
				State: Unread,
			},
			handler: func(t *testing.T) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					data, err := io.ReadAll(r.Body)
					require.NoError(t, err)

					req := retrieveRequest{}
					require.NoError(t, json.Unmarshal(data, &req))

					require.Equal(t, consumerKey, req.ConsumerKey)
					require.Equal(t, accessToken, req.AccessToken)
					require.Equal(t, Unread, req.State)

					resp := &RetrieveResponse{
						Status: successStatus,
					}
					data, err = json.Marshal(resp)
					require.NoError(t, err)
					w.Header().Add("Content-type", "application/json")
					_, err = w.Write(data)
					require.NoError(t, err)
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

			res, err := p.Retrieve(context.Background(), tc.ad)

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
