package session

import (
	"net/http"
)

type Writer interface {
	WriteSession(rw http.ResponseWriter, accessToken *string)
	ClearSession(rw http.ResponseWriter)
}