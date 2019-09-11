package kong

import (
	"context"
	go_kong "github.com/hbagdi/go-kong/kong"
)

type UpstreamsInterface interface {
	Create(context context.Context, upstream go_kong.Upstream) (go_kong.Upstream, error)
}

type ClientInterface struct {
	Upstreams UpstreamsInterface
}
