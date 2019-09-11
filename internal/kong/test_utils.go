package kong

import "github.com/ciroque/k8s-kong-federated-ingress/internal/k8s"

func CompareStringArrays(l []string, r []string) bool {
	if len(l) != len(r) {
		return false
	}
	for i, e := range l {
		if r[i] != e {
			return false
		}
	}
	return true
}

func ServicesMatch(l k8s.Service, r k8s.Service) bool {
	return l.Name == r.Name &&
		CompareStringArrays(l.Paths, r.Paths) &&
		l.Port == r.Port &&
		CompareStringArrays(l.Addresses, r.Addresses)
}
