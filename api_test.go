package webapi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type testApiServer struct {
	mux    *http.ServeMux
	server *httptest.Server
}

func newTestServer() (*testApiServer, *Client) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	client, err := DefaultClient("test-token")

	if err != nil {
		log.Fatalln(err)
	}

	url, _ := url.Parse(server.URL)
	client.BaseUrl = url

	ts := &testApiServer{
		mux:    mux,
		server: server,
	}

	return ts, client
}

func (ts *testApiServer) Close() {
	ts.server.Close()
}

func (ts *testApiServer) HandleFunc(t *testing.T, method, pattern, filename string, outRes interface{}) {
	ts.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Method; got != method {
			t.Errorf("Request method: %v, want %v", got, method)
		}
		if len(filename) > 0 {
			stub, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Fatalln(err)
			}
			if err := json.Unmarshal([]byte(stub), outRes); err != nil {
				log.Fatalln(err)
			}
			if stub != nil {
				w.Write([]byte(stub))
			}
		}
	})
}
