## Overview

This is simple SDK implementation for the [Pocket](https://getpocket.com/developer/?src=footer_v2).
Before using it, you should create a [Pocket Application](https://getpocket.com/developer/apps/new)

The SDK implements all developer API features:

- [Authentication](https://getpocket.com/developer/docs/authentication)
- [Modify](https://getpocket.com/developer/docs/v3/modify)
- [Retrieve](https://getpocket.com/developer/docs/v3/retrieve)
- [Add](https://getpocket.com/developer/docs/v3/add)

## Content

- [Installation](#installation)
- [Create a pocket object](#create-a-pocket-object)
- [Authentication](#authentication)
  - [Generate a request token](#generate-a-request-token)
  - [Generate an authorization link](#generate-an-authorization-link)
  - [Generate an access token](#generate-an-access-token)
- [Add](#add)
- [Retrieve](#retrieve)
- [Modification](#modification)
  - [Actions](#actions)
  - [Tags](#tags)
  - [Usage](#usage)
- [Errors](#errors)

## Installation

To install the package run:

```bash
go get -u github.com/VladimirStepanov/pocket-golang-sdk
```

## Create a pocket object

For creating a Pocket object, you need to use the consumer key, which you've gotten after app registration.

```go
p := pocket.New("consumer-key")
```

## Authentication

Authentication performs in 3 steps:
- get a request token
- user authorization
- convert a request token into a Pocket access token

### Generate a request token

Method 1 
```go
err := p.AuthApp(context.Background(), "redirect-url")

if err != nil {
    log.Fatal(err)
}

fmt.Println("Request token", p.GetRequestToken())
```

Method 2
```go
res, err := p.GenerateRequestToken(context.Background(), "redirect-url")

if err != nil {
    log.Fatal(err)
}

p.SetRequestToken(res.Code)
fmt.Println("Request token", p.GetRequestToken())
```

### Generate an authorization link

Once you've had a request token, you need to redirect the user to Pocket to authorize your
application's request token. For getting authorize link, use this method:

```go
link := p.MakeAuthUrl("redirect-url")

// result example: https://getpocket.com/auth/authorize?redirect_uri=https%3A%2F%2Fgoogle.com&request_token=ffffcc4e-ffff-ffff-ffff-f7f68f 
```

### Generate an access token

After successfully user authorization, you can get an access token

Method 1
```go
err := p.AuthUser(context.Background())

if err != nil {
    log.Fatal(err)
}

fmt.Println("Access token", p.GetAccessToken())
```

Method 2
```go
at, err := p.GenerateAccessToken(context.Background())
if err != nil {
    log.Fatal(err)
}
p.SetAccessToken(at.AccessToken)
fmt.Println("Access token", p.GetAccessToken())
```

## Add

**NOTE**: You can add multiple items at the same time. See [Modification](#modification).

Input model:
```go
type AddInput struct {
    Url     string `json:"url"` // The URL of the item you want to save. MUST BE ENCODED
    Title   string `json:"title,omitempty"`
    Tags    string `json:"tags,omitempty"` // A comma-separated list of tags to apply to the item
    TweetID string `json:"tweet_id,omitempty"` // If you are adding Pocket support to a Twitter client, please send along a reference to the tweet status
}
```

Response model:
```go
type AddResponse struct {
    Item   Item `json:"item"` // very big struct, see file add.go
    Status int  `json:"status"`
}
```

Usage:
```go
_, err = p.Add(context.Background(), &pocket.AddInput{
    Url:  "https://www.youtube.com/watch?v=fJHNhL1FUEs&ab_channel=GolangCafe",
    Tags: "codding",
})

if err != nil {
    log.Fatal(err)
}
```

## Retrieve

Retrieve model:
```go
type RetrieveInput struct {
	State       State       `json:"state,omitempty"` // [unread | archive | all]. see  State consts
	Favorite    *Favorite   `json:"favorite"` //  [0 - un-favorited | 1 = favorited]. see Favorite consts
	Tag         string      `json:"tag,omitempty"` // [*tag_name* | _untagged_]
	ContentType ContentType `json:"contentType,omitempty"` // [article | video | image]. see ContentType consts
	Sort        Sort        `json:"sort,omitempty"` // [newest | oldest | title | site]. see Sort consts
	DetailType  DetailType  `json:"detailType,omitempty"` // [simple | complete]. see DetailType consts
	Search      string      `json:"search,omitempty"` // Only return items whose title or url contain the search string
	Domain      string      `json:"domain,omitempty"` // Only return items from a particular domain
	Since       *int64      `json:"since,omitempty"` // Only return items modified since the given since unix timestamp
	Count       int64       `json:"count,omitempty"` // Only return count number of items
	Offset      int64       `json:"offset,omitempty"` // Used only with count; start returning from offset position of results
}
```

Retrieve response:
```go
type RetrieveResponse struct {
	Status     int                         `json:"status"`
	Complete   int                         `json:"complete"`
	SearchMeta SearchMeta                  `json:"search_meta"`
	Since      int                         `json:"since"`
	List       map[string]RetrieveListItem `json:"list"` // very big struct, see file retrieve.go
}
```

Example:

```go
retrRes, err := p.Retrieve(context.Background(), &pocket.RetrieveInput{
    State:       pocket.Unread,
    Tag:         pocket.Untagged,
    ContentType: pocket.ArticleType,
    Sort:        pocket.Title,
    DetailType:  pocket.Simple,
})

if err != nil {
    log.Fatal(err)
}

fmt.Println(retrRes)
```

## Modification

[Modify](https://getpocket.com/developer/docs/v3/modify) method accept different actions in one array. For
 the actions exist special type:

```go
type Actions []interface{}
```

### Actions

Every action is a structure, with special Action field. The field has own value for every action.

```go
const (
    ActionAddType         ActionType = "add"
    ActionArchiveType     ActionType = "archive"
    ActionReaddType       ActionType = "readd"
    ActionFavoriteType    ActionType = "favorite"
    ActionUnfavoriteType  ActionType = "unfavorite"
    ActionDeleteType      ActionType = "delete"
)

type (
    action struct {
        Action ActionType `json:"action"`
        ItemID int64      `json:"item_id"`
        Time   int64      `json:"time,omitempty"`
    }
    
    ActionAdd struct {
        Action ActionType `json:"action"`
        RefID  int64      `json:"ref_id,omitempty"` // A Twitter status id; this is used to show tweet attribution
        Tags   string     `json:"tags,omitempty"` // A comma-delimited list of one or more tags
        Time   int64      `json:"time,omitempty"` // The time the action occurred
        Title  string     `json:"title,omitempty"` // 	The title of the item
        Url    string     `json:"url"` // The url of the item; provide this only if you do not have. MUST BE ENCODED
    }
)

type (
    ActionArchive     action
    ActionReadd       action
    ActionFavorite    action
    ActionUnfavorite  action
    ActionDelete      action
)
```

### Tags

```go
const (
    ActionTagsAddType     ActionType = "tags_add"
    ActionTagsRemoveType  ActionType = "tags_remove"
    ActionTagsReplaceType ActionType = "tags_replace"
    ActionTagsClearType   ActionType = "tags_clear"
    ActionTagRenameType   ActionType = "tag_rename"
    ActionTagDeleteType   ActionType = "tag_delete"
)

type (
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
)

type (
	ActionTagsAdd     tagsAction
	ActionTagsRemove  tagsAction
	ActionTagsReplace tagsAction
	ActionTagsClear   action
)
```

### Usage

```go
modRes, err := p.Modify(context.Background(), pocket.Actions{
    &pocket.ActionAdd{
        Action: pocket.ActionAddType,
        Tags:   "codding",
        Url:    "https://www.youtube.com/watch?v=fJHNhL1FUEs&ab_channel=GolangCafe",
    },
    &pocket.ActionDelete{
        Action: pocket.ActionDeleteType,
        ItemID: 777,
    },
    &pocket.ActionTagsAdd{
        Action: pocket.ActionTagsAddType,
        ItemID: 777,
        Tags:   "tag1,tag2",
    },
})

if err != nil {
    log.Fatal(err)
}

fmt.Println(modRes)
```

## Errors

For Pocket errors exist this structure:

```go
type ErrorPocket struct {
	Message  string
	Xcode    string // see X-Code-Error here https://getpocket.com/developer/docs/authentication
	HttpCode int
}
```