// This package implements a synchronized MRU (most recently used) map. Its a map that will hold at most cap key/value pairs
// and will evict entries if: 1) adding an element would exceed cap, in this case the oldest existing element is removed to
// make space for a new one. 2) an element's age exceeds maxage, the check is done when adding elements.
// Age is defined as the duration between now and the time the element was last accessed (added or looked up).
package mru

import (
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type holder struct {
	last int64
	k    interface{}
	v    interface{}
}

type holders []*holder

func (hs holders) Len() int {
	return len(hs)
}

func (hs holders) Less(i, j int) bool {
	return hs[j].last < hs[i].last
}

func (hs holders) Swap(i, j int) {
	t := hs[i]
	hs[i] = hs[j]
	hs[j] = t
}

type Map struct {
	m      map[interface{}]*holder
	idx    holders
	vidx   holders
	locker sync.RWMutex
	maxage time.Duration
}

func NewMap(cap int, maxage time.Duration) *Map {
	if !(cap > 0) {
		panic("capacity must be greater than zero")
	}

	idx := make(holders, cap)
	return &Map{
		m:      map[interface{}]*holder{},
		idx:    make(holders, cap),
		vidx:   idx[0:0],
		maxage: maxage,
	}
}

func (m *Map) Add(k interface{}, v interface{}) {
	m.locker.Lock()
	defer func() {
		m.locker.Unlock()
	}()

	sort.Sort(m.vidx)

	oldest := time.Now().Add(-m.maxage)
	var keep int
	if len(m.vidx) == cap(m.idx) {
		keep = len(m.vidx) - 1
		delete(m.m, m.vidx[keep].k)
		m.vidx[keep] = nil
		keep = len(m.vidx) - 2

	} else {
		keep = len(m.vidx) - 1
	}

	for ; keep >= 0; keep-- {
		if m.vidx[keep].last > oldest.UnixNano() {
			break
		}
		delete(m.m, m.vidx[keep].k)
		m.vidx[keep] = nil
	}
	m.vidx = m.idx[0 : keep+1]

	newHolder := &holder{
		last: time.Now().UnixNano(),
		k:    k,
		v:    v,
	}
	m.m[k] = newHolder
	m.vidx = m.idx[0 : len(m.vidx)+1]
	m.vidx[len(m.vidx)-1] = newHolder
}

func (m *Map) Lookup(k interface{}) (v interface{}, ok bool) {
	m.locker.RLock()
	defer func() {
		m.locker.RUnlock()
	}()

	h, ok := m.m[k]
	if ok {
		v = h.v
		_ = atomic.SwapInt64(&h.last, time.Now().UnixNano())
	}
	return v, ok
}
