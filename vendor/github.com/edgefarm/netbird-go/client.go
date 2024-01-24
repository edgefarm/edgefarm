/*
Copyright © 2024 EdgeFarm Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package netbird

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	token   string
	http    *http.Client
	address string
	timeout time.Duration
}

type ClientOption func(*Client)

func WithManagementAddress(address string) ClientOption {
	return func(c *Client) {
		c.address = address
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.timeout = timeout
	}
}

func NewClient(token string, opts ...ClientOption) *Client {
	c := &Client{
		token: token,
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
		address: "https://api.netbird.io",
	}

	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) doCall(method string, endpoint string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", c.address, endpoint), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Token "+c.token)

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resp, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
