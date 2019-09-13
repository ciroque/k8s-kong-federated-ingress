package k8s

type ServiceDef struct {
	Addresses []string
	Namespace string
	Paths     []string
}

type ServicesMap map[string]ServiceDef
