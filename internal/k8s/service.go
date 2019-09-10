package k8s

type Service struct {
	Addresses []string
	Name      string
	Paths     []string
	Port      int
}
