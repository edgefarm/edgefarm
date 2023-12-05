package context

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetContext(t *testing.T) {
	assert := assert.New(t)
	foo := Context("foo")
	assert.Equal("foo", foo.Name)
	assert.NotNil(foo.Data)
	assert.Equal(0, len(foo.Data))
	foo.Set("key", "value")
	value, ok := foo.Get("key")
	assert.True(ok)
	assert.Equal("value", value)
	assert.Equal(1, len(foo.Data))

	value, ok = foo.Get("key2")
	assert.False(ok)
	assert.Nil(value)

	barData := map[string]interface{}{
		"foo": "baz",
		"test": map[string]int{
			"one": 1,
			"two": 2,
		},
	}

	bar := Context("bar", WithData(barData))
	assert.Equal("bar", bar.Name)
	assert.NotNil(bar.Data)
	assert.Equal(2, len(bar.Data))
	first, ok := bar.Get("foo")
	assert.True(ok)
	assert.Equal("baz", first)
	second, ok := bar.Get("test")
	assert.True(ok)
	assert.NotNil(second)
	assert.Equal(2, len(second.(map[string]int)))
	assert.Equal(1, second.(map[string]int)["one"])
	assert.Equal(2, second.(map[string]int)["two"])

	value, ok = foo.Get("key")
	assert.True(ok)
	assert.Equal("value", value)

	assert.Equal(1, len(foo.Data))

	secondFoo := Context("foo")
	assert.Equal("foo", secondFoo.Name)
	assert.NotNil(secondFoo.Data)
	assert.Equal(1, len(secondFoo.Data))
	value, ok = secondFoo.Get("key")
	assert.True(ok)
	assert.Equal("value", value)
}

func TestExists(t *testing.T) {
	assert := assert.New(t)
	assert.False(Exists("asdf"))
	asdf := Context("asdf")
	assert.True(Exists("asdf"))
	assert.Equal("asdf", asdf.Name)
	assert.NotNil(asdf.Data)
	assert.Equal(0, len(asdf.Data))
}
