package mock

import (
	"bufio"
	"net"
	"net/http"
)

type MockResponseWriter struct{}

func (m *MockResponseWriter) Header() http.Header {
	return map[string][]string{}
}

func (m *MockResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (m *MockResponseWriter) WriteHeader(statusCode int) {

}

func (m *MockResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, &bufio.ReadWriter{}, nil
}
