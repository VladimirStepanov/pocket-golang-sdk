package pocket

import "errors"

var (
	ErrMissingConsumerKey = errors.New("missing consumer key")
	ErrInvalidConsumerKey = errors.New("invalid consumer key")
	ErrMissingRedirectUri = errors.New("missing redirect uri")
	ErrInvalidRedirectUri = errors.New("invalid redirect uri")
	ErrPocketServerIssue  = errors.New("pocket server issue")
	ErrMissingCode        = errors.New("missing code")
	ErrNoCode             = errors.New("code not found")
	ErrRejectedCode       = errors.New("user rejected code")
	ErrCodeAlreadyUsed    = errors.New("already used code")
)
