package webapi

import (
	"testing"
)

func TestDefaultClient_tokenIsEmpty(t *testing.T) {
	c, err := DefaultClient("")

	if c != nil {
		t.Errorf("Expected nil to be returned")
	}
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestNewClient_baseUrlIsEmpty(t *testing.T) {
	c, err := NewClient("test token", "", defaultHttpClient, defaultLogger)

	if c != nil {
		t.Errorf("Expected nil to be returned")
	}
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestNewClient_httpClientIsNil(t *testing.T) {
	c, err := NewClient("test token", defaultBaseUrl, nil, defaultLogger)

	if c != nil {
		t.Errorf("Expected nil to be returned")
	}
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestNewClient_loggerIsNil(t *testing.T) {
	c, err := NewClient("test token", defaultBaseUrl, defaultHttpClient, nil)

	if c != nil {
		t.Errorf("Expected nil to be returned")
	}
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}
