// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

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
