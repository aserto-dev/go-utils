package mru

import (
	"testing"
	"time"

	"github.com/alecthomas/assert"
)

func TestAddLookup(t *testing.T) {
	m := NewMRUMap(3, time.Duration(time.Second*3600))
	m.Add("one", 1)
	m.Add("two", 2)
	m.Add("three", 3)
	_, ok := m.Lookup("one")
	assert.True(t, ok)
	v, ok := m.Lookup("two")
	assert.True(t, ok)
	assert.Equal(t, v, 2)
	_, ok = m.Lookup("three")
	assert.True(t, ok)

	m.Add("four", 4)
	_, ok = m.Lookup("one")
	assert.False(t, ok)
	_, ok = m.Lookup("two")
	assert.True(t, ok)
	_, ok = m.Lookup("three")
	assert.True(t, ok)
	_, ok = m.Lookup("four")
	assert.True(t, ok)

	m.maxage = time.Duration(time.Second * 0)
	m.Add("five", 5)
	_, ok = m.Lookup("one")
	assert.False(t, ok)
	_, ok = m.Lookup("two")
	assert.False(t, ok)
	_, ok = m.Lookup("three")
	assert.False(t, ok)
	_, ok = m.Lookup("four")
	assert.False(t, ok)
	_, ok = m.Lookup("five")
	assert.True(t, ok)
}
