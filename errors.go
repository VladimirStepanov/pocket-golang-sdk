package pocket

import (
	"fmt"
)

type ErrorPocket struct {
	Message string
	// see X-Code-Error here https://getpocket.com/developer/docs/authentication
	Xcode    int
	HttpCode int
}

func (pe *ErrorPocket) Error() string {
	return fmt.Sprintf("%s-%d-%d", pe.Message, pe.Xcode, pe.HttpCode)
}

func NewErrorPocket(message string, xCode int, httpCode int) *ErrorPocket {
	return &ErrorPocket{
		Message:  message,
		Xcode:    xCode,
		HttpCode: httpCode,
	}
}
