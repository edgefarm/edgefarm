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

type PeerGroup struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	PeersCount int    `json:"peers_count,omitempty"`
	Issued     string `json:"issued,omitempty"`
}

type Peers []Peer

type Peer struct {
	ID                     string            `json:"id,omitempty"`
	Name                   string            `json:"name,omitempty"`
	IP                     string            `json:"ip,omitempty"`
	Connected              bool              `json:"connected,omitempty"`
	LastSeen               string            `json:"last_seen,omitempty"`
	Os                     string            `json:"os,omitempty"`
	Version                string            `json:"version,omitempty"`
	Groups                 []PeerGroup       `json:"groups,omitempty"`
	SSHEnabled             bool              `json:"ssh_enabled,omitempty"`
	UserID                 string            `json:"user_id,omitempty"`
	Hostname               string            `json:"hostname,omitempty"`
	UIVersion              string            `json:"ui_version,omitempty"`
	DNSLabel               string            `json:"dns_label,omitempty"`
	LoginExpirationEnabled bool              `json:"login_expiration_enabled,omitempty"`
	LoginExpired           bool              `json:"login_expired,omitempty"`
	LastLogin              string            `json:"last_login,omitempty"`
	ApprovalRequired       bool              `json:"approval_required,omitempty"`
	AccessiblePeers        []AccessiblePeers `json:"accessible_peers,omitempty"`
	AcessiblePeersCount    int               `json:"accessible_peers_count,omitempty"`
}
type AccessiblePeers struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	IP       string `json:"ip,omitempty"`
	DNSLabel string `json:"dns_label,omitempty"`
	UserID   string `json:"user_id,omitempty"`
}

func (c *Client) ListPeers() (Peers, error) {
	peers := Peers{}
	body, err := c.doCall(http.MethodGet, "api/peers", nil)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &peers)
	if err != nil {
		return nil, err
	}

	return peers, nil
}

func (c *Client) GetPeer(id string) (Peer, error) {
	peer := Peer{}
	body, err := c.doCall(http.MethodGet, fmt.Sprintf("api/peers/%s", id), nil)
	if err != nil {
		return Peer{}, err
	}
	err = json.Unmarshal(body, &peer)
	if err != nil {
		return Peer{}, err
	}
	return peer, nil
}

func (c *Client) GetPeerIdByHostname(hostname string) (string, error) {
	peers, err := c.ListPeers()
	if err != nil {
		return "", err
	}
	for _, peer := range peers {
		if peer.Hostname == hostname {
			return peer.ID, nil
		}
	}
	return "", fmt.Errorf("peer not found")
}

func (c *Client) DeletePeer(id string) error {
	_, err := c.doCall(http.MethodDelete, fmt.Sprintf("api/peers/%s", id), nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdatePeer(p *Peer) (*Peer, error) {
	j, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	data := strings.NewReader(string(j))

	body, err := c.doCall(http.MethodPut, fmt.Sprintf("api/peers/%s", p.ID), data)
	if err != nil {
		return nil, err
	}
	resp := Peer{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
