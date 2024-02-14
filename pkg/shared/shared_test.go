package shared

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnly(t *testing.T) {
	assert := assert.New(t)
	onlyFlags := OnlyFlags{
		Network: true,
		Flannel: true,
	}

	skipFlags := ConvertOnlyToSkip(onlyFlags)
	assert.Nil(nil)
	fmt.Printf("%#v\n", skipFlags)
}
