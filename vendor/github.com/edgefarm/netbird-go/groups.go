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

type Groups []Group

type Group struct {
	ID         string       `json:"id,omitempty"`
	Name       string       `json:"name,omitempty"`
	PeersCount int          `json:"peers_count,omitempty"`
	Issued     string       `json:"issued,omitempty"`
	Peers      []GroupPeers `json:"peers,omitempty"`
}
type GroupPeers struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func (c *Client) ListGroups() ([]Group, error) {
	groups := []Group{}
	body, err := c.doCall(http.MethodGet, "api/groups", nil)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &groups)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (c *Client) GetGroup(id string) (*Group, error) {
	group := &Group{}
	body, err := c.doCall(http.MethodGet, fmt.Sprintf("api/groups/%s", id), nil)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, group)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (c *Client) GetGroupByName(name string) (*Group, error) {
	groups, err := c.ListGroups()
	if err != nil {
		return nil, err
	}
	for _, g := range groups {
		if g.Name == name {
			return &g, nil
		}
	}
	return nil, fmt.Errorf("group not found")
}

func (c *Client) DeleteGroup(id string) error {
	_, err := c.doCall(http.MethodDelete, fmt.Sprintf("api/groups/%s", id), nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateGroup(p *Group) (*Group, error) {
	j, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	data := strings.NewReader(string(j))

	body, err := c.doCall(http.MethodPut, fmt.Sprintf("api/groups/%s", p.ID), data)
	if err != nil {
		return nil, err
	}
	resp := Group{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) CreateGroup(p *Group) (*Group, error) {
	j, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	data := strings.NewReader(string(j))

	body, err := c.doCall(http.MethodPost, "api/groups", data)
	if err != nil {
		return nil, err
	}
	resp := Group{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
