package rest

import (
	"io"
	"net"
	"net/http"
	"time"
)

type req struct {
	Headers map[string]string
	Timeout int
}

func NewClient(to int) *http.Client {
	timeout := time.Duration(to) * time.Second
	return &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				deadline := time.Now().Add(timeout)
				c, err := net.DialTimeout(network, addr, timeout)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
			DisableKeepAlives:     true,
			ResponseHeaderTimeout: timeout,
			DisableCompression:    false,
		},
	}
}

func NewRequest(timeout int, headers map[string]string) *req {
	return &req{
		Headers: headers,
		Timeout: timeout,
	}
}

func (r *req) Get(url string) (*http.Response, error) {
	c := NewClient(r.Timeout)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}
	return c.Do(req)
}

func (r *req) Post(url string, body io.Reader) (*http.Response, error) {
	c := NewClient(r.Timeout)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accpet", "application/json")
	req.Header.Set("Content-Type", "application/json")
	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}
	return c.Do(req)
}
