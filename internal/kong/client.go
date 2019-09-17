package kong

import (
	"context"
	gokong "github.com/hbagdi/go-kong/kong"
)

type RoutesInterface interface {
	Create(ctx context.Context, route *gokong.Route) (*gokong.Route, error)
}

type ServicesInterface interface {
	Create(context context.Context, service *gokong.Service) (*gokong.Service, error)
}

type TargetsInterface interface {
	Create(context context.Context, target *gokong.Target) (*gokong.Target, error)
}

type UpstreamsInterface interface {
	Create(context context.Context, upstream *gokong.Upstream) (*gokong.Upstream, error)
}

type Client struct {
	Routes    RoutesInterface
	Services  ServicesInterface
	Targets   TargetsInterface
	Upstreams UpstreamsInterface
}

type service struct {
	Kong gokong.Client
}

type Routes service

func (routes Routes) Create(context context.Context, route *gokong.Route) (*gokong.Route, error) {
	return routes.Kong.Routes.Create(context, route)
}

type Services service

func (services Services) Create(context context.Context, service *gokong.Service) (*gokong.Service, error) {
	return services.Kong.Services.Create(context, service)
}

type Targets service

func (targets Targets) Create(context context.Context, target *gokong.Target) (*gokong.Target, error) {
	return targets.Kong.Targets.Create(context, target.Upstream.Name, target)
}

type Upstreams service

func (upstreams Upstreams) Create(context context.Context, upstream *gokong.Upstream) (*gokong.Upstream, error) {
	return upstreams.Kong.Upstreams.Create(context, upstream)
}
