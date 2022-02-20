package rest

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type HttpStatusCode int

type Requester interface {
	Get(url string, headers map[string]string) ([]byte, HttpStatusCode, error)
	Post(url string, headers map[string]string, body io.Reader) ([]byte, HttpStatusCode, error)
}

type RestClient struct{}

func (h *RestClient) Get(url string, headers map[string]string) ([]byte, HttpStatusCode, error) {
	return nil, 000, nil
}

func (h *RestClient) Post(url string, headers map[string]string, body io.Reader) ([]byte, HttpStatusCode, error) {
	return nil, 000, nil
}

//Makes an http request and returns the response as a slice of bytes.
//Returns response, status code and error (if any)
func makeHttpRequest(url, method string, body io.Reader, headers map[string]string) ([]byte, HttpStatusCode, error) {
	c := http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, 000, errors.Wrap(err, "couldn't initialize the request")
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, 000, errors.Wrap(err, "couldn't execute the request")
	}
	defer res.Body.Close()
	var statusCode HttpStatusCode = HttpStatusCode(res.StatusCode)

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, statusCode, errors.Wrap(err, "couldn't read the response")
	}

	return bytes, statusCode, nil
}
