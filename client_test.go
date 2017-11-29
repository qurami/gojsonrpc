package gojsonrpc

import (
	"reflect"
	"testing"
)

func TestThatNewClientReturnsTheExpectedClient(t *testing.T) {
	mockURL := "http://mock.url"

	expectedClient := &Client{
		URL:     mockURL,
		Timeout: defaultTimeout,
	}

	sut := NewClient(mockURL)

	if !reflect.DeepEqual(sut, expectedClient) {
		t.Fatal("expected Client was not received.")
	}
}

func TestThatSetTimeoutSucceeds(t *testing.T) {
	mockURL := "http://mock.url"
	mockTimeout := 123

	expectedClient := &Client{
		URL:     mockURL,
		Timeout: 123,
	}

	sut := NewClient(mockURL)
	sut.SetTimeout(mockTimeout)

	if !reflect.DeepEqual(sut, expectedClient) {
		t.Fatal("expected Client was not received.")
	}
}

func TestThatSetProxySucceeds(t *testing.T) {
	mockURL := "http://mock.url"
	mockProxyURL := "http://proxy.url:1234"

	expectedClient := &Client{
		URL:          mockURL,
		Timeout:      defaultTimeout,
		proxyAddress: mockProxyURL,
	}

	sut := NewClient(mockURL)
	sut.SetHTTPProxy(mockProxyURL)

	if !reflect.DeepEqual(sut, expectedClient) {
		t.Fatal("expected Client was not received.")
	}
}
