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

type SetupKeys []SetupKey

type SetupKey struct {
	ID         string   `json:"id,omitempty"`
	Key        string   `json:"key,omitempty"`
	Name       string   `json:"name,omitempty"`
	ExpiresIn  int      `json:"expires_in,omitempty"`
	Expires    string   `json:"expires,omitempty"`
	Type       string   `json:"type,omitempty"`
	Valid      bool     `json:"valid,omitempty"`
	Revoked    bool     `json:"revoked,omitempty"`
	UsedTimes  int      `json:"used_times,omitempty"`
	LastUsed   string   `json:"last_used,omitempty"`
	State      string   `json:"state,omitempty"`
	AutoGroups []string `json:"auto_groups,omitempty"`
	UpdatedAt  string   `json:"updated_at,omitempty"`
	UsageLimit int      `json:"usage_limit,omitempty"`
	Ephemeral  bool     `json:"ephemeral,omitempty"`
}

func (c *Client) ListSetupKeys() (SetupKeys, error) {
	setupKeys := SetupKeys{}
	body, err := c.doCall(http.MethodGet, "api/setup-keys", nil)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &setupKeys)
	if err != nil {
		return nil, err
	}

	return setupKeys, nil
}

func (c *Client) GetSetupKeyByName(name string) ([]SetupKey, error) {
	setupKeys := []SetupKey{}
	keys, err := c.ListSetupKeys()
	if err != nil {
		return []SetupKey{}, err
	}
	for _, k := range keys {
		if k.Name == name {
			setupKeys = append(setupKeys, k)
		}
	}
	if len(setupKeys) == 0 {
		return []SetupKey{}, fmt.Errorf("no setup key with name %s found", name)
	}
	return setupKeys, nil
}

func (c *Client) GetSetupKey(id string) (SetupKey, error) {
	setupKey := SetupKey{}
	body, err := c.doCall(http.MethodGet, fmt.Sprintf("api/setup-keys/%s", id), nil)
	if err != nil {
		return SetupKey{}, err
	}
	err = json.Unmarshal(body, &setupKey)
	if err != nil {
		return SetupKey{}, err
	}
	return setupKey, nil
}

func (c *Client) CreateSetupKey(k *SetupKey) (*SetupKey, error) {
	j, err := json.Marshal(k)
	if err != nil {
		return nil, err
	}
	data := strings.NewReader(string(j))
	body, err := c.doCall(http.MethodPost, "api/setup-keys", data)
	if err != nil {
		return nil, err
	}
	resp := SetupKey{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) DeleteSetupKey(id string) error {
	_, err := c.doCall(http.MethodDelete, fmt.Sprintf("api/setup-keys/%s", id), nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateSetupKey(k *SetupKey) (SetupKey, error) {
	j, err := json.Marshal(k)
	if err != nil {
		return SetupKey{}, err
	}
	data := strings.NewReader(string(j))

	body, err := c.doCall(http.MethodPut, fmt.Sprintf("api/setup-keys/%s", k.ID), data)
	if err != nil {
		return SetupKey{}, err
	}
	err = json.Unmarshal(body, &k)
	if err != nil {
		return SetupKey{}, err
	}
	return *k, nil
}
