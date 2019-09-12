package k8s

type ServiceDef struct {
	Addresses []string
	Name      string /// TODO Deprecate, prefer ServicesMap struct that maps from Name to ServiceDef
	Namespace string
	Paths     []string
	Port      int
}

type ServicesMap map[string]ServiceDef
