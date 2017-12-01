package gojsonrpc

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestThatSetHTTPProxyURLSetsAnHTTPProxyInClient(t *testing.T) {
	proxyCalled := false

	mockProxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxyCalled = true
	}))

	mockRPCServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	sut := NewClient(mockRPCServer.URL)

	mockHTTPProxyURL, _ := url.Parse(mockProxy.URL)

	sut.SetHTTPProxyURL(mockHTTPProxyURL)
	_ = sut.Notify("SometMethod", nil)

	if !proxyCalled {
		t.Fatal("expected proxy was not called")
	}
}
