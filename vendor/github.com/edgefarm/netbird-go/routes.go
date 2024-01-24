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
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Routes []Route

type Route struct {
	ID          string   `json:"id,omitempty"`
	NetworkType string   `json:"network_type,omitempty"`
	Description string   `json:"description,omitempty"`
	NetworkID   string   `json:"network_id,omitempty"`
	Enabled     bool     `json:"enabled,omitempty"`
	Peer        string   `json:"peer,omitempty"`
	PeerGroups  []string `json:"peer_groups,omitempty"`
	Network     string   `json:"network,omitempty"`
	Metric      int      `json:"metric,omitempty"`
	Masquerade  bool     `json:"masquerade,omitempty"`
	Groups      []string `json:"groups,omitempty"`
}

func (c *Client) ListRoutes() (Routes, error) {
	routes := Routes{}
	body, err := c.doCall(http.MethodGet, "api/routes", nil)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &routes)
	if err != nil {
		return nil, err
	}

	return routes, nil
}

func (c *Client) GetRoute(id string) (Route, error) {
	route := Route{}
	body, err := c.doCall(http.MethodGet, fmt.Sprintf("api/routes/%s", id), nil)
	if err != nil {
		return Route{}, err
	}
	err = json.Unmarshal(body, &route)
	if err != nil {
		return Route{}, err
	}
	return route, nil
}

func (c *Client) DeleteRoute(id string) error {
	_, err := c.doCall(http.MethodDelete, fmt.Sprintf("api/routes/%s", id), nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateRoute(p *Route) (*Route, error) {
	j, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	data := strings.NewReader(string(j))

	body, err := c.doCall(http.MethodPut, fmt.Sprintf("api/routes/%s", p.ID), data)
	if err != nil {
		return nil, err
	}
	resp := Route{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) CreateRoute(p *Route) (*Route, error) {
	j, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	data := strings.NewReader(string(j))

	body, err := c.doCall(http.MethodPost, "api/routes", data)
	if err != nil {
		return nil, err
	}
	resp := Route{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
