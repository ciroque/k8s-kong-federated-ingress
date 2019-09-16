package kong

type ServiceDef struct {
	ServiceName  string
	Routes       []string
	UpstreamName string
	Targets      []string
}
