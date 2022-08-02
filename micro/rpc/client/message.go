package client

import "net/http"

type Message struct {
	Header http.Header
	Body interface{}
}