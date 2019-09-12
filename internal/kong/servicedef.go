package kong

type ServiceDef struct {
	Service  string
	Routes   []string
	Upstream string
	Targets  []string
}
