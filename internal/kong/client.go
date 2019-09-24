package kong

import (
	"context"
	gokong "github.com/hbagdi/go-kong/kong"
)

type RoutesInterface interface {
	Create(ctx context.Context, route *gokong.Route) (*gokong.Route, error)
	Delete(context context.Context, routeNameOrId *string) error
}

type ServicesInterface interface {
	Create(context context.Context, service *gokong.Service) (*gokong.Service, error)
	Delete(context context.Context, serviceNameOrId *string) error
}

type TargetsInterface interface {
	Create(context context.Context, target *gokong.Target) (*gokong.Target, error)
	Delete(context context.Context, upstreamNameOrId *string, targetOrId *string) error
}

type UpstreamsInterface interface {
	Create(context context.Context, upstream *gokong.Upstream) (*gokong.Upstream, error)
	Delete(context context.Context, upstreamNameOrId *string) error
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
	if r, err := routes.Kong.Routes.Get(context, route.Name); err == nil {
		route.ID = r.ID // update if exists
	}
	return routes.Kong.Routes.Create(context, route)
}

func (routes Routes) Delete(context context.Context, routeNameOrId *string) error {
	return routes.Kong.Routes.Delete(context, routeNameOrId)
}

type Services service

func (services Services) Create(context context.Context, service *gokong.Service) (*gokong.Service, error) {
	if s, err := services.Kong.Services.Get(context, service.Name); err == nil {
		service.ID = s.ID // update if exists
	}
	return services.Kong.Services.Create(context, service)
}

func (services Services) Delete(context context.Context, serviceNameOrId *string) error {
	return services.Kong.Services.Delete(context, serviceNameOrId)
}

type Targets service

func (targets Targets) Create(context context.Context, target *gokong.Target) (*gokong.Target, error) {
	return targets.Kong.Targets.Create(context, target.Upstream.Name, target)
}

func (targets Targets) Delete(context context.Context, upstreamNameOrId *string, targetOrId *string) error {
	return targets.Kong.Targets.Delete(context, upstreamNameOrId, targetOrId)
}

type Upstreams service

func (upstreams Upstreams) Create(context context.Context, upstream *gokong.Upstream) (*gokong.Upstream, error) {
	if u, err := upstreams.Kong.Upstreams.Get(context, upstream.Name); err == nil {
		upstream.ID = u.ID // update if exists
	}
	return upstreams.Kong.Upstreams.Create(context, upstream)
}

func (upstreams Upstreams) Delete(context context.Context, upstreamNameOrId *string) error {
	return upstreams.Kong.Upstreams.Delete(context, upstreamNameOrId)
}
