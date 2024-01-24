package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	stateDirectory = ".edgefarm-local-up"
	stateFile      = "state.json"
)

type CurrentState struct {
	Netbird  NetbirdState `json:"netbird"`
	filePath string       `json:"-"`
}

type NetbirdState struct {
	FullyConfigured   bool   `json:"fully_configured"`
	NetbirdToken      string `json:"netbird_token"`
	NetbirdSetupKey   string `json:"netbird_setup_key"`
	NetbirdSetupKeyID string `json:"netbird_setup_key_id"`
	NetbirdGroupID    string `json:"netbird_group_id"`
	NetbirdRouteID    string `json:"netbird_route_id"`
}

var currentState *CurrentState

type CurrentStateOption func(*CurrentState)

func WithStoragePath(path string) CurrentStateOption {
	return func(s *CurrentState) {
		s.filePath = path
	}
}

func GetState(opts ...CurrentStateOption) (*CurrentState, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	tmp := &CurrentState{
		filePath: fmt.Sprintf("%s/%s", homeDir, stateDirectory),
	}
	for _, opt := range opts {
		opt(tmp)
	}

	currentState, err = importFromFile(tmp.filePath)
	if err != nil {
		currentState = tmp
	}

	currentState.filePath = tmp.filePath

	if _, err := os.Stat(filepath.Base(currentState.filePath)); os.IsNotExist(err) {
		err := os.MkdirAll(fmt.Sprintf("%s/%s", homeDir, stateDirectory), 0755)
		if err != nil {
			return nil, err
		}
	}
	return currentState, nil
}

func (s *CurrentState) IsFullyConfigured() bool {
	return s.Netbird.FullyConfigured
}

func (s *CurrentState) SetFullyConfigured() {
	s.Netbird.FullyConfigured = true
	s.Export()
}

func (s *CurrentState) SetNetbirdSetupKey(key string) {
	s.Netbird.NetbirdSetupKey = key
	s.Export()
}

func (s *CurrentState) GetNetbirdSetupKey() string {
	return s.Netbird.NetbirdSetupKey
}

func (s *CurrentState) SetNetbirdToken(token string) {
	s.Netbird.NetbirdToken = token
	s.Export()
}

func (s *CurrentState) GetNetbirdToken() string {
	return s.Netbird.NetbirdToken
}

func (s *CurrentState) SetNetbirdSetupKeyID(id string) {
	s.Netbird.NetbirdSetupKeyID = id
	s.Export()
}

func (s *CurrentState) GetNetbirdSetupKeyID() string {
	return s.Netbird.NetbirdSetupKeyID
}

func (s *CurrentState) SetNetbirdGroupID(id string) {
	s.Netbird.NetbirdGroupID = id
	s.Export()
}

func (s *CurrentState) GetNetbirdGroupID() string {
	return s.Netbird.NetbirdGroupID
}

func (s *CurrentState) SetNetbirdRouteID(id string) {
	s.Netbird.NetbirdRouteID = id
	s.Export()
}

func (s *CurrentState) GetNetbirdRouteID() string {
	return s.Netbird.NetbirdRouteID
}

func (s *CurrentState) Export() error {
	j, err := json.Marshal(s)
	if err != nil {
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", s.filePath, stateFile), j, 0644)
	if err != nil {
		return err
	}
	return nil
}

func importFromFile(path string) (*CurrentState, error) {
	content, err := os.ReadFile(fmt.Sprintf("%s/%s", path, stateFile))
	if err != nil {
		return nil, err
	}
	c := &CurrentState{}
	err = json.Unmarshal([]byte(content), c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CurrentState) Clear() {
	s.Netbird = NetbirdState{}
	s.Export()
}
