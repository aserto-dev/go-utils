package mru

import (
	"testing"
	"time"

	"github.com/alecthomas/assert"
)

func TestAddLookup(t *testing.T) {
	m := NewMap(3, time.Duration(time.Second*3600))
	m.Add("one", 1)
	m.Add("two", 2)
	m.Add("three", 3)
	_, ok := m.Lookup("one", false, false)
	assert.True(t, ok)
	v, ok := m.Lookup("two", false, false)
	assert.True(t, ok)
	assert.Equal(t, v, 2)
	_, ok = m.Lookup("three", false, false)
	assert.True(t, ok)

	_, ok = m.Lookup("one", false, true)
	assert.True(t, ok)
	m.Add("four", 4)
	_, ok = m.Lookup("one", false, false)
	assert.True(t, ok)
	_, ok = m.Lookup("two", false, false)
	assert.False(t, ok)
	_, ok = m.Lookup("three", false, false)
	assert.True(t, ok)
	_, ok = m.Lookup("four", false, false)
	assert.True(t, ok)

	m.maxage = time.Duration(time.Second * 0)
	m.Add("five", 5)
	m.Add("five", "five")
	_, ok = m.Lookup("one", false, false)
	assert.False(t, ok)
	_, ok = m.Lookup("two", false, false)
	assert.False(t, ok)
	_, ok = m.Lookup("three", false, false)
	assert.False(t, ok)
	_, ok = m.Lookup("four", false, false)
	assert.False(t, ok)
	v, ok = m.Lookup("five", false, false)
	assert.True(t, ok)
	assert.Equal(t, v, "five")
	_, ok = m.Lookup("five", true, false)
	assert.False(t, ok)
}
