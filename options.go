/**
 *
 * @author liangjf
 * @create on 2020/5/9
 * @version 1.0
 */
package gglog

import (
	"context"

	"github.com/liangjfblue/gglog/vlog"
)

type Options struct {
	Name    string
	Log     vlog.Log
	Before  []func()
	After   []func()
	Context context.Context
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Log:     DefaultLog,
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

func Log(log gglog.Log) Option {
	return func(o *Options) {
		o.Log = log
	}
}

func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

func Before(fn func()) Option {
	return func(o *Options) {
		o.Before = append(o.Before, fn)
	}
}

func After(fn func()) Option {
	return func(o *Options) {
		o.After = append(o.After, fn)
	}
}
