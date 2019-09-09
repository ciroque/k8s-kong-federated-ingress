package k8s

type Service struct {
	Addresses []*string
	Name      string
	Path      string
	Port      int
}
