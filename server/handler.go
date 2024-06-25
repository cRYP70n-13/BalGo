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

// TODO: So now what we have to do is:
//   - Match the request with the correct backend server
//   - Check the heatlh of the backend server if it's healthy proceed if it's not
//   - then just drop this shit and return an error
//   - Add the correct headers to the request
//   - Periodically check the health of the servers and their replicas
//   - The request forwarding should be based on the strategy we have in the config file.
//     Hmm another idea is that I can do some checks (Where the request is coming from) on the load balancer
//     side then if it didn't work we gonna drop the connection with an
//     error.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	w.WriteHeader(res.StatusCode)
	_, _ = w.Write(body)
}

// Lets not use this atm.
func getDestinationUrl(r *http.Request) (*url.URL, error) {
	parsedUrl, err := url.Parse(r.URL.String())
	if err != nil {
		return nil, err
	}
	fmt.Println(parsedUrl.Path, parsedUrl.Host, parsedUrl)

	switch parsedUrl.Path {
	case "alpha":
		return &url.URL{
			// TODO: Get this from the config
			Scheme: "http",
			Host:   "127.0.0.1:8081",
			Path:   parsedUrl.Path,
		}, nil
	case "beta":
		return &url.URL{
			// TODO: Get this from the config
			Scheme: "http",
			Host:   "127.0.0.1:8082",
			Path:   parsedUrl.Path,
		}, nil
	case "gama":
		return &url.URL{
			// TODO: Get this from the config
			Scheme: "http",
			Host:   "127.0.0.1:8083",
			Path:   parsedUrl.Path,
		}, nil
	case "sigma":
		return &url.URL{
			// TODO: Get this from the config
			Scheme: "http",
			Host:   "127.0.0.1:8084",
			Path:   parsedUrl.Path,
		}, nil
	}

	return parsedUrl, err
}
