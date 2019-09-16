package kong

type ServiceDef struct {
	ServiceName  string
	RoutesMap    map[string]string
	UpstreamName string
	Targets      []string
}
