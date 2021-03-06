package pocket

import (
	"context"
)

type ActionType string
type Actions []interface{}

const (
	ActionAddType         ActionType = "add"
	ActionArchiveType     ActionType = "archive"
	ActionReaddType       ActionType = "readd"
	ActionFavoriteType    ActionType = "favorite"
	ActionUnfavoriteType  ActionType = "unfavorite"
	ActionDeleteType      ActionType = "delete"
	ActionTagsAddType     ActionType = "tags_add"
	ActionTagsRemoveType  ActionType = "tags_remove"
	ActionTagsReplaceType ActionType = "tags_replace"
	ActionTagsClearType   ActionType = "tags_clear"
	ActionTagRenameType   ActionType = "tag_rename"
	ActionTagDeleteType   ActionType = "tag_delete"
)

type (
	ActionArchive     action
	ActionReadd       action
	ActionFavorite    action
	ActionUnfavorite  action
	ActionDelete      action
	ActionTagsAdd     tagsAction
	ActionTagsRemove  tagsAction
	ActionTagsReplace tagsAction
	ActionTagsClear   action
)

type (
	action struct {
		Action ActionType `json:"action"`
		ItemID int64      `json:"item_id"`
		Time   int64      `json:"time,omitempty"`
	}

	ActionAdd struct {
		Action ActionType `json:"action"`
		RefID  int64      `json:"ref_id,omitempty"`
		Tags   string     `json:"tags,omitempty"`
		Time   int64      `json:"time,omitempty"`
		Title  string     `json:"title,omitempty"`
		Url    string     `json:"url"` // MUST BE ENCODED
	}

	tagsAction struct {
		Action ActionType `json:"action"`
		ItemID int64      `json:"item_id"`
		Tags   string     `json:"tags"`
		Time   int64      `json:"time,omitempty"`
	}

	ActionTagRename struct {
		Action ActionType `json:"action"`
		OldTag string     `json:"old_tag"`
		NewTag string     `json:"new_tag"`
		Time   int64      `json:"time,omitempty"`
	}

	ActionTagDelete struct {
		Action ActionType `json:"action"`
		Tag    string     `json:"tag"`
		Time   int64      `json:"time,omitempty"`
	}

	modifyRequest struct {
		Actions     Actions `json:"actions"`
		ConsumerKey string  `json:"consumer_key"`
		AccessToken string  `json:"access_token"`
	}

	ModifyResponse struct {
		ActionResult []interface{} `json:"action_results"`
		Status       int           `json:"status"`
	}
)

func (p *Pocket) Modify(ctx context.Context, actions Actions) (*ModifyResponse, error) {
	req := modifyRequest{
		Actions:     actions,
		ConsumerKey: p.consumerKey,
		AccessToken: p.accessToken,
	}

	res := ModifyResponse{}
	err := p.doRequest(ctx, modifyPath, req, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
