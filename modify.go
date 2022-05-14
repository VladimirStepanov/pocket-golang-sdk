package pocket

import (
	"context"
	"encoding/json"
	"fmt"
)

type ActionType string

const (
	ActionAddType        ActionType = "add"
	ActionArchiveType    ActionType = "archive"
	ActionReaddType      ActionType = "readd"
	ActionFavoriteType   ActionType = "favorite"
	ActionUnfavoriteType ActionType = "unfavorite"
	ActionDeleteType     ActionType = "delete"
)

type Actions []interface{}

type Action struct {
	Action ActionType `json:"action"`
	ItemID int64      `json:"item_id"`
	Time   int64      `json:"time,omitempty"`
}

type ActionAdd struct {
	Action ActionType `json:"action"`
	RefID  int64      `json:"ref_id,omitempty"`
	Tags   string     `json:"tags,omitempty"`
	Time   int64      `json:"time,omitempty"`
	Title  string     `json:"title,omitempty"`
	Url    string     `json:"url,omitempty"` // MUST BE ENCODED
}

type modifyRequest struct {
	Actions     Actions `json:"actions"`
	ConsumerKey string  `json:"consumer_key"`
	AccessToken string  `json:"access_token"`
}

type ModifyResponse struct {
	ActionResult []interface{} `json:"action_results"`
	Status       int           `json:"status"`
}

func (p *Pocket) Modify(ctx context.Context, actions Actions) (*ModifyResponse, error) {
	req := modifyRequest{
		Actions:     actions,
		ConsumerKey: p.consumerKey,
		AccessToken: p.accessToken,
	}

	data, err := p.doRequestRaw(ctx, modifyPath, req)
	if err != nil {
		return nil, err
	}

	res := ModifyResponse{}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling response data: %w", err)
	}
	return &res, nil
}
