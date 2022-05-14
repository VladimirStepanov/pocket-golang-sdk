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
	actionKey = "action"
)

func checkActionHandler(t *testing.T, expectedAction ActionType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		req := modifyRequest{}
		require.NoError(t, json.Unmarshal(data, &req))

		require.Equal(t, consumerKey, req.ConsumerKey)
		require.Equal(t, accessToken, req.AccessToken)
		addReq, ok := req.Actions[0].(map[string]interface{})

		require.True(t, ok)
		reqAct, ok := addReq[actionKey].(string)
		require.True(t, ok)
		require.Equal(t, expectedAction, ActionType(reqAct))

		resp := &ModifyResponse{
			Status: successStatus,
		}
		data, err = json.Marshal(resp)
		require.NoError(t, err)
		w.Header().Add("Content-type", "application/json")
		w.Write(data)
		w.WriteHeader(http.StatusOK)
	}
}

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
			name: "Success [action=add]",
			actions: Actions{
				&ActionAdd{
					Action: ActionAddType,
				},
			},
			handler: func(t *testing.T) http.HandlerFunc {
				return checkActionHandler(t, ActionAddType)
			},
		},
		{
			name: "Success [action=archive]",
			actions: Actions{
				&ActionArchive{
					Action: ActionArchiveType,
				},
			},
			handler: func(t *testing.T) http.HandlerFunc {
				return checkActionHandler(t, ActionArchiveType)
			},
		},
		{
			name: "Success [action=readd]",
			actions: Actions{
				&ActionReadd{
					Action: ActionReaddType,
				},
			},
			handler: func(t *testing.T) http.HandlerFunc {
				return checkActionHandler(t, ActionReaddType)
			},
		},
		{
			name: "Success [action=favorite]",
			actions: Actions{
				&ActionFavorite{
					Action: ActionFavoriteType,
				},
			},
			handler: func(t *testing.T) http.HandlerFunc {
				return checkActionHandler(t, ActionFavoriteType)
			},
		},
		{
			name: "Success [action=unfavorite]",
			actions: Actions{
				&ActionUnfavorite{
					Action: ActionUnfavoriteType,
				},
			},
			handler: func(t *testing.T) http.HandlerFunc {
				return checkActionHandler(t, ActionUnfavoriteType)
			},
		},
		{
			name: "Success [action=delete]",
			actions: Actions{
				&ActionDelete{
					Action: ActionDeleteType,
				},
			},
			handler: func(t *testing.T) http.HandlerFunc {
				return checkActionHandler(t, ActionDeleteType)
			},
		},
		{
			name: "Success [action=tags_add]",
			actions: Actions{
				&ActionTagsAdd{
					Action: ActionTagsAddType,
				},
			},
			handler: func(t *testing.T) http.HandlerFunc {
				return checkActionHandler(t, ActionTagsAddType)
			},
		},
		{
			name: "Success [action=tags_remove]",
			actions: Actions{
				&ActionTagsRemove{
					Action: ActionTagsRemoveType,
				},
			},
			handler: func(t *testing.T) http.HandlerFunc {
				return checkActionHandler(t, ActionTagsRemoveType)
			},
		},
		{
			name: "Success [action=tags_replace]",
			actions: Actions{
				&ActionTagsReplace{
					Action: ActionTagsReplaceType,
				},
			},
			handler: func(t *testing.T) http.HandlerFunc {
				return checkActionHandler(t, ActionTagsReplaceType)
			},
		},
		{
			name: "Success [action=tags_clear]",
			actions: Actions{
				&ActionTagsClear{
					Action: ActionTagsClearType,
				},
			},
			handler: func(t *testing.T) http.HandlerFunc {
				return checkActionHandler(t, ActionTagsClearType)
			},
		},
		{
			name: "Success [action=tag_rename]",
			actions: Actions{
				&ActionTagRename{
					Action: ActionTagRenameType,
				},
			},
			handler: func(t *testing.T) http.HandlerFunc {
				return checkActionHandler(t, ActionTagRenameType)
			},
		},
		{
			name: "Success [action=tag_delete]",
			actions: Actions{
				&ActionTagDelete{
					Action: ActionTagDeleteType,
				},
			},
			handler: func(t *testing.T) http.HandlerFunc {
				return checkActionHandler(t, ActionTagDeleteType)
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
