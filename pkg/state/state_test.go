package state

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestState(t *testing.T) {
	assert := assert.New(t)
	dir, err := os.MkdirTemp("", "state")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
	state, err := GetState(WithStoragePath(dir))
	assert.Nil(err)
	assert.NotNil(state)

	state.SetNetbirdSetupKeyID("foo")
	assert.Equal("foo", state.Netbird.NetbirdSetupKeyID)
	state.SetNetbirdSetupKey("bar")
	assert.Equal("bar", state.Netbird.NetbirdSetupKey)

	loadedState, err := GetState(WithStoragePath(dir))
	assert.Nil(err)
	assert.NotNil(loadedState)
	assert.Equal("foo", loadedState.Netbird.NetbirdSetupKeyID)
	assert.Equal("bar", loadedState.Netbird.NetbirdSetupKey)
}

func TestStateDefaultLocation(t *testing.T) {
	assert := assert.New(t)
	state, err := GetState()
	assert.Nil(err)
	assert.NotNil(state)

	state.SetNetbirdSetupKeyID("foo")
	assert.Equal("foo", state.Netbird.NetbirdSetupKeyID)
	state.SetNetbirdSetupKey("bar")
	assert.Equal("bar", state.Netbird.NetbirdSetupKey)

	loadedState, err := GetState()
	assert.Nil(err)
	assert.NotNil(loadedState)
	assert.Equal("foo", loadedState.Netbird.NetbirdSetupKeyID)
	assert.Equal("bar", loadedState.Netbird.NetbirdSetupKey)
}
