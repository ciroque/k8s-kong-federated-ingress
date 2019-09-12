package kong

import (
	"context"
	gokong "github.com/hbagdi/go-kong/kong"
)

type RoutesInterface interface {
	Create(ctx context.Context, route gokong.Route) (gokong.Route, error)
}

type ServicesInterface interface {
	Create(context context.Context, upstream gokong.Service) (gokong.Service, error)
}

type TargetsInterface interface {
	Create(context context.Context, target gokong.Target) (gokong.Target, error)
}

type UpstreamsInterface interface {
	Create(context context.Context, upstream gokong.Upstream) (gokong.Upstream, error)
}

type ClientInterface struct {
	Routes    RoutesInterface
	Services  ServicesInterface
	Targets   TargetsInterface
	Upstreams UpstreamsInterface
}
