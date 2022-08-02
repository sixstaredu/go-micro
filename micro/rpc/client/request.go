package client

import "net/http"

type rpcRequest struct {
	service     string
	method      string
	endpoint    string
	contentType string
	body        interface{}
	opts		requestOptiones

	header		http.Header
}

func newRequest(service, endpoint string, request interface{}, reqOpts ...RequestOption) Request {
	var opts requestOptiones

	for _, o := range reqOpts {
		o(&opts)
	}

	return &rpcRequest{
		service:     service,
		method:      endpoint,
		endpoint:    endpoint,
		body:        request,
		opts:        opts,

		header:http.Header{},
	}
}

func (r *rpcRequest) ContentType() string {
	return r.contentType
}

func (r *rpcRequest) Service() string {
	return r.service
}

func (r *rpcRequest) Method() string {
	return r.method
}

func (r *rpcRequest) Endpoint() string {
	return r.endpoint
}

func (r *rpcRequest) Body() interface{} {
	return r.body
}

func (r *rpcRequest) SetHeader(key string, value interface{}) {
	//r.header[key] = value
}

func (r *rpcRequest) Header() http.Header {
	return r.header
}
