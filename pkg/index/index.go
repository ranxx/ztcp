package index

import "sync/atomic"

// Index64 ...
type Index64 struct {
	i int64
}

// NewIndexI64 ...
func NewIndexI64() *Index64 {
	return &Index64{i: -1}
}

// NewIndex ...
func (i *Index64) NewIndex() int64 {
	return atomic.AddInt64(&i.i, 1)
}
