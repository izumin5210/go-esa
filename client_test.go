package webapi

import (
	"reflect"
	"testing"
)

type TestResponse struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}
type TestCase struct {
	method, pattern, filename string
	want, got                 TestResponse
}

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

func TestNewResponseAndDo(t *testing.T) {
	testCase := TestCase{
		method:   "GET",
		pattern:  "/sample",
		filename: "./tests/stub/sample.json",
	}
	ts, c := newTestServer()
	ts.HandleFunc(t, testCase.method, testCase.pattern, testCase.filename, &testCase.want)

	req, err := c.NewRequest(testCase.method, testCase.pattern, nil)
	if err != nil {
		t.Errorf("Error create new request: %v", err)
	}
	resp, err := c.Do(req, &testCase.got)
	if err != nil {
		t.Errorf("Error do a request: %v", err)
	}
	if resp == nil {
		t.Errorf("Expected response to be returned: %v", err)
	}

	if !reflect.DeepEqual(testCase.got, testCase.want) {
		t.Errorf("Response: %v, want %v", testCase.got, testCase.want)
	}
}

func TestNewResponseAndDo_responseIsEmpty(t *testing.T) {
	testCase := TestCase{
		method:  "DELETE",
		pattern: "/sample",
	}
	ts, c := newTestServer()
	ts.HandleFunc(t, testCase.method, testCase.pattern, testCase.filename, &testCase.want)

	req, err := c.NewRequest(testCase.method, testCase.pattern, nil)
	if err != nil {
		t.Errorf("Error create new request: %v", err)
	}
	resp, err := c.Do(req, nil)
	if err != nil {
		t.Errorf("Error do a request: %v", err)
	}
	if resp == nil {
		t.Errorf("Expected response to be returned: %v", err)
	}

	if !reflect.DeepEqual(testCase.got, testCase.want) {
		t.Errorf("Response: %v, want %v", testCase.got, testCase.want)
	}
}
