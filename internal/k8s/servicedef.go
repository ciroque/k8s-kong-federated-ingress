package k8s

type ServiceDef struct {
	Addresses []string
	Namespace string
	Paths     []string
	Port      int
}

type ServicesMap map[string]ServiceDef
