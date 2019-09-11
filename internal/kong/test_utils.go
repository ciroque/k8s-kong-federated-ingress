package kong

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

func ServicesMatch(l ServiceDef, r ServiceDef) bool {
	return l.Name == r.Name &&
		CompareStringArrays(l.Paths, r.Paths) &&
		l.Port == r.Port &&
		CompareStringArrays(l.Addresses, r.Addresses) &&
		ResourceNamesMatch(l.Names, r.Names)
}

func ResourceNamesMatch(l ResourceNames, r ResourceNames) bool {
	return l.ServiceName == r.ServiceName && l.UpstreamName == r.UpstreamName
}
