package server

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Handler struct{}

func NewHandler() *Handler {
    return &Handler{}
}

// TODO: Now I need to forward the request that I got here to the backend server
// based on what I have in the config
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Hmm another idea is that I can do some checks on the load balancer
	// side then if it didn't work we gonna drop the connection with an
	// error
	if r.RequestURI == "/hello" {
		if _, err := w.Write([]byte("Hello World\n")); err != nil {
			fmt.Println("error when writing response for /hello request")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	url, err := url.Parse("http://localhost:8081/healthz")
	if err != nil {
		panic(err)
	}

	client := http.Client{}
	request := &http.Request{
		Method: r.Method,
		URL:    url,
	}
	res, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

    fmt.Println("******** RESPONSE *********", string(body))

	w.WriteHeader(http.StatusNotFound)
}
