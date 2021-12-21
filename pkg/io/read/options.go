package read

import (
	"github.com/ranxx/ztcp/pkg/message"
	"github.com/ranxx/ztcp/pkg/pack"
)

// Options ...
type Options struct {
	packer pack.Packer
}

// Option ...
type Option func(*Options)

// DefaultOptions ...
func DefaultOptions() *Options {
	return &Options{
		packer: pack.DefaultPacker(message.DefaultMessager),
	}
}

// WithPacker ...
func WithPacker(packer pack.Packer) Option {
	return func(o *Options) {
		o.packer = packer
	}
}
