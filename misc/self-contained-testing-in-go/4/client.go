package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

func main() {
	// Nothing to see here
}

// START CLIENT OMIT
type Client struct {
	u *url.URL
}

func NewClient(addr string) *Client {
	u, _ := url.Parse(addr)
	return &Client{u: u}
}

func (c *Client) Status() (string, error) {
	u, _ := url.Parse("/status")
	res, _ := http.Get(c.u.ResolveReference(u).String())
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	return string(body), err
}

// END CLIENT OMIT
