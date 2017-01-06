package webapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const (
	libraryVersion = "1"

	defaultBaseUrl   = "https://api.example.com"
	defaultUserAgent = "go-webapi-client-template/" + libraryVersion
)

var (
	defaultLogger     = log.New(ioutil.Discard, "", log.LstdFlags)
	defaultHttpClient = http.DefaultClient
)

type Client struct {
	HttpClient *http.Client
	Logger     *log.Logger
	BaseUrl    *url.URL
	UserAgent  string

	Token string
}

func DefaultClient(token string) (*Client, error) {
	return NewClient(token, defaultBaseUrl, defaultHttpClient, defaultLogger)
}

func NewClient(token, urlStr string, httpClient *http.Client, logger *log.Logger) (*Client, error) {
	if len(token) == 0 {
		return nil, errors.New("You must set a token.")
	}

	baseUrl, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to parse url %s: ", urlStr)
	}

	if httpClient == nil {
		return nil, errors.New("You must set a *http.Client.")
	}

	if logger == nil {
		return nil, errors.New("You must set a *log.Logger.")
	}

	return &Client{
		HttpClient: httpClient,
		Logger:     logger,
		UserAgent:  defaultUserAgent,
		BaseUrl:    baseUrl,
		Token:      token,
	}, nil
}

type Response struct {
	*http.Response
}

type ErrorResponse struct {
	Response *http.Response
	Type     string `json:"error"`
	Message  string `json:"message"`
}

func (r *ErrorResponse) Error() string {
	// NOTE: Shoud sanitize r.Response.Request.URL if the URL query has some secret params.
	return fmt.Sprintf("%v %v: %d %v %+v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message, r.Type)
}

func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseUrl.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	response := &Response{Response: resp}

	err = checkResponse(resp)
	if err != nil {
		return response, err
	}

	defer resp.Body.Close()

	if v != nil {
		if err := responseUnmarshal(resp.Body, v); err != nil {
			return nil, err
		}
	}

	return response, err
}

func checkResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}

func responseUnmarshal(body io.ReadCloser, v interface{}) error {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}
