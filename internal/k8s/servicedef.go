package k8s

type ServiceDef struct {
	Addresses []string
	Name      string
	Namespace string
	Paths     []string
	Port      int
}
