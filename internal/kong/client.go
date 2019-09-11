package kong

import (
	"context"
	go_kong "github.com/hbagdi/go-kong/kong"
)

type ServicesInterface interface {
	Create(context context.Context, upstream go_kong.Service) (go_kong.Service, error)
}

type UpstreamsInterface interface {
	Create(context context.Context, upstream go_kong.Upstream) (go_kong.Upstream, error)
}

type ClientInterface struct {
	Services  ServicesInterface
	Upstreams UpstreamsInterface
}
