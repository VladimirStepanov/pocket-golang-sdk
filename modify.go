package pocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	ActionAddType = "add"
)

type Actions []interface{}

func (actions *Actions) AppendActionAdd(ad *ActionAdd) (*ActionAdd, error) {
	u, err := url.Parse(ad.Url)
	if err != nil {
		return nil, fmt.Errorf("error while parsing url field from ActionAdd struct: %w", err)
	}
	ad.Url = u.String()

	ad.Action = ActionAddType

	*actions = append(*actions, ad)

	return ad, nil
}

type ActionAdd struct {
	Action string `json:"action"`
	RefID  int64  `json:"ref_id,omitempty"`
	Tags   string `json:"tags,omitempty"`
	Time   int64  `json:"time,omitempty"`
	Title  string `json:"title,omitempty"`
	Url    string `json:"url,omitempty"`
}

type modifyRequest struct {
	Actions     []interface{} `json:"actions"`
	ConsumerKey string        `json:"consumer_key"`
	AccessToken string        `json:"access_token"`
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
