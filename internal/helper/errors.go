package helper

import (
	"fmt"
	"strings"
)

const (
	NOTFOUND = "NOT FOUND"
	INTERNALSERVER  = "INTERNAL SERVER"
	CONFLICT        = "CONFLICT ERROR"
)

func HandleErrors(err error, model string) (status int, message string) {
	if err == nil {
		return 200, ""
	}
	errMessage := err.Error()
	message = fmt.Sprintf("err: %s, model: %s", errMessage, model)
	if strings.Contains(errMessage, NOTFOUND) {
		return 404, message
	} else if strings.Contains(errMessage, CONFLICT) {
		return 409, message
	} else {
		return 500, message
	}
}