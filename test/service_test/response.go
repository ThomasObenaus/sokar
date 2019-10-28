package main

import "net/http"

type Response struct {
	statusCode int
	data       []byte
	header     http.Header
}

type ResponseOption func(r *Response)

func StatusCode(code int) ResponseOption {
	return func(r *Response) {
		r.statusCode = code
	}
}

func AddHeader(key, value string) ResponseOption {
	return func(r *Response) {
		r.header.Add(key, value)
	}
}

func SetHeader(key, value string) ResponseOption {
	return func(r *Response) {
		r.header.Set(key, value)
	}
}

func DelHeader(key, value string) ResponseOption {
	return func(r *Response) {
		r.header.Del(key)
	}
}

func NewStringResponse(data string, options ...ResponseOption) Response {
	resp := Response{
		statusCode: http.StatusOK,
		data:       []byte(data),
		header:     make(http.Header, 0),
	}

	// apply the options
	for _, opt := range options {
		opt(&resp)
	}
	return resp
}

func NewBinResponse(data []byte, options ...ResponseOption) Response {
	resp := Response{
		statusCode: http.StatusOK,
		data:       data,
		header:     make(http.Header, 0),
	}

	// apply the options
	for _, opt := range options {
		opt(&resp)
	}
	return resp
}
