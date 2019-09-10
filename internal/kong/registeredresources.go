package kong

import "github.com/hbagdi/go-kong/kong"

type RegisteredResources struct {
	Service  *kong.Service
	Targets  []*kong.Target
	Route    []*kong.Route
	Upstream *kong.Upstream
}
