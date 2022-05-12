package pocket

import (
	"context"
	"encoding/json"
	"fmt"
)

// state

type State string

const (
	Unread  State = "unread"
	Archive State = "archive"
	All     State = "all"
)

// favorite

type Favorite int

const (
	Unfavorited Favorite = 0
	Favorited   Favorite = 1
)

//tag

const (
	Untagged = "_untagged_"
)

// contentType

type ContentType string

const (
	ArticleType ContentType = "article"
	VideoType   ContentType = "video"
	ImageType   ContentType = "image"
)

// sort

type Sort string

const (
	Newest Sort = "newest"
	Oldest Sort = "oldest"
	Title  Sort = "title"
	Site   Sort = "site"
)

// detailType

type DetailType string

const (
	Simple   DetailType = "simple"
	Complete DetailType = "complete"
)

type RetrieveInput struct {
	State       State       `json:"state,omitempty"`
	Favorite    *Favorite   `json:"favorite"`
	Tag         string      `json:"tag,omitempty"`
	ContentType ContentType `json:"contentType,omitempty"`
	Sort        Sort        `json:"sort,omitempty"`
	DetailType  DetailType  `json:"detailType,omitempty"`
	Search      string      `json:"search,omitempty"`
	Domain      string      `json:"domain,omitempty"`
	Since       *int64      `json:"since,omitempty"`
	Count       int64       `json:"count,omitempty"`
	Offset      int64       `json:"offset,omitempty"`
}

type SearchMeta struct {
	SearchType string `json:"search_type"`
}

type Author struct {
	ItemID   string `json:"item_id"`
	AuthorID string `json:"author_id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
}

type Image struct {
	ItemID string `json:"item_id"`
	Src    string `json:"src"`
	Width  string `json:"width"`
	Height string `json:"height"`
}

type ImagesItem struct {
	ItemID  string `json:"item_id"`
	ImageID string `json:"image_id"`
	Src     string `json:"src"`
	Width   string `json:"width"`
	Height  string `json:"height"`
	Credit  string `json:"credit"`
	Caption string `json:"caption"`
}

type Video struct {
	ItemID  string `json:"item_id"`
	VideoID string `json:"video_id"`
	Src     string `json:"src"`
	Width   string `json:"width"`
	Height  string `json:"height"`
	Type    string `json:"type"`
	Vid     string `json:"vid"`
	Length  string `json:"length"`
}

type RetrieveListItem struct {
	ItemID                 string                `json:"item_id"`
	ResolvedID             string                `json:"resolved_id"`
	GivenURL               string                `json:"given_url"`
	GivenTitle             string                `json:"given_title"`
	Favorite               string                `json:"favorite"`
	Status                 string                `json:"status"`
	TimeAdded              string                `json:"time_added"`
	TimeUpdated            string                `json:"time_updated"`
	TimeRead               string                `json:"time_read"`
	TimeFavorited          string                `json:"time_favorited"`
	SortID                 int                   `json:"sort_id"`
	ResolvedTitle          string                `json:"resolved_title"`
	ResolvedURL            string                `json:"resolved_url"`
	Excerpt                string                `json:"excerpt"`
	IsArticle              string                `json:"is_article"`
	IsIndex                string                `json:"is_index"`
	HasVideo               string                `json:"has_video"`
	HasImage               string                `json:"has_image"`
	WordCount              string                `json:"word_count"`
	Lang                   string                `json:"lang"`
	Authors                map[string]Author     `json:"authors"`
	Image                  Image                 `json:"image"`
	Images                 map[string]ImagesItem `json:"images"`
	Videos                 map[string]Video      `json:"videos"`
	DomainMetadata         DomainMetadata        `json:"domain_metadata"`
	ListenDurationEstimate int                   `json:"listen_duration_estimate"`
}

type RetrieveResponse struct {
	Status     int                         `json:"status"`
	Complete   int                         `json:"complete"`
	SearchMeta SearchMeta                  `json:"search_meta"`
	Since      int                         `json:"since"`
	List       map[string]RetrieveListItem `json:"list"`
}

type retrieveRequest struct {
	*RetrieveInput
	ConsumerKey string `json:"consumer_key"`
	AccessToken string `json:"access_token"`
}

func (p *Pocket) Retrieve(ctx context.Context, rd *RetrieveInput) (*RetrieveResponse, error) {
	req := retrieveRequest{
		RetrieveInput: rd,
		ConsumerKey:   p.consumerKey,
		AccessToken:   p.accessToken,
	}

	data, err := p.doRequestRaw(ctx, retrievePath, req)
	if err != nil {
		return nil, err
	}

	res := RetrieveResponse{}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling request data: %w", err)
	}

	return &res, nil
}
