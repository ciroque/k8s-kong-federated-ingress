package kong

type KongServiceDef struct {
	ServiceName  string
	Routes       []string
	UpstreamName string
	Targets      []string
}
