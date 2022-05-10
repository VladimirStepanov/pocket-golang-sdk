package pocket

import (
	"fmt"
)

type ErrorPocket struct {
	Message  string
	Xcode    string // see X-Code-Error here https://getpocket.com/developer/docs/authentication
	HttpCode int
}

func (pe *ErrorPocket) Error() string {
	return fmt.Sprintf("%s-%s-%d", pe.Message, pe.Xcode, pe.HttpCode)
}

func NewErrorPocket(message string, xCode string, httpCode int) *ErrorPocket {
	return &ErrorPocket{
		Message:  message,
		Xcode:    xCode,
		HttpCode: httpCode,
	}
}
