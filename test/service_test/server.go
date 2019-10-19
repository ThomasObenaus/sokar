package main

import (
	"io"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/thomasobenaus/sokar/api"
)

type server struct {
	receiver *api.API
	mock     *MockHTTP
	mockCtrl *gomock.Controller
}

func New(port int, t *testing.T) *server {

	return nil

	//mockCtrl := gomock.NewController(t)
	//
	//srv := server{
	//	receiver: api.New(port),
	//	mock:     NewMockHTTP(mockCtrl),
	//	mockCtrl: mockCtrl,
	//}
	//
	//srv.receiver.Run()
	//
	//return &srv
}

func (s *server) Close() {
	s.receiver.Stop()
	s.mockCtrl.Finish()
}

func (s *server) POST(expect string, returnCode int, returnData string) {
	//receiver.Router.HandlerFunc("GET", "/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//
	//	code, _ := mock.POST("HELLO")
	//	w.Header().Set("Content-Type", "application/json")
	//	w.WriteHeader(code)
	//	io.WriteString(w, "data")
	//
	//	//receiver.Stop()
	//}))
}

func (s *server) GET(path string, returnCode int, returnData string) {
	s.mock.EXPECT().GET(path).Return(returnCode, returnData)

	s.receiver.Router.HandlerFunc("GET", path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code, data := s.mock.GET(path)
		w.WriteHeader(code)
		io.WriteString(w, data)
	}))
}
