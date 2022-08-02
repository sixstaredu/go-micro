package server

import (
	"encoding/json"
	"net/http"
)

type Message struct {
	Header http.Header

	Body json.RawMessage
}
