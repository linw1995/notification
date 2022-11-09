package notification

import "context"

type SendOptions struct {
	Title       string
	Description string
}

type SendOption interface {
	Apply(*SendOptions)
}

type Title string

func (t Title) Apply(o *SendOptions) {
	o.Title = string(t)
}

type Description string

func (t Description) Apply(o *SendOptions) {
	o.Description = string(t)
}

func GenerateOptions(opts ...SendOption) (rv SendOptions) {
	for _, opt := range opts {
		opt.Apply(&rv)
	}
	return
}

type Channel interface {
	Send(ctx context.Context, opts ...SendOption) (err error)
}
