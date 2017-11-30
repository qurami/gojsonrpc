package gojsonrpc

import (
	"testing"
	"time"
)

func TestThatNewClientReturnsTheExpectedClient(t *testing.T) {
	mockURL := "http://mock.url"

	sut := NewClient(mockURL)

	if sut.timeout != defaultTimeout ||
		sut.proxyURL != "" ||
		sut.httpClient.Timeout != time.Duration(defaultTimeout)*time.Second {
		t.Fatal("expected Client was not received.")
	}
}

func TestThatSetTimeoutSucceeds(t *testing.T) {
	mockURL := "http://mock.url"
	mockTimeout := 123

	sut := NewClient(mockURL)
	sut.SetTimeout(mockTimeout)

	if sut.timeout != 123 ||
		sut.proxyURL != "" ||
		sut.httpClient.Timeout != time.Duration(123)*time.Second {
		t.Fatal("expected Client was not received.")
	}
}

func TestThatSetProxySucceeds(t *testing.T) {
	mockURL := "http://mock.url"
	mockProxyURL := "http://proxy.url:1234"

	sut := NewClient(mockURL)
	sut.SetHTTPProxy(mockProxyURL)

	if sut.timeout != defaultTimeout ||
		sut.proxyURL != mockProxyURL {
		t.Fatal("expected Client was not received.")
	}
}
