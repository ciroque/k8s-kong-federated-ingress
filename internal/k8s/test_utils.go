package k8s

import "reflect"

func CompareStringArrays(expected []string, actual []string) bool {
	if len(expected) != len(actual) {
		return false
	}
	for i, e := range expected {
		if actual[i] != e {
			return false
		}
	}
	return true
}

func ServicesMatch(expected ServiceDef, actual ServiceDef) bool {
	return expected.Name == actual.Name &&
		CompareStringArrays(expected.Paths, actual.Paths) &&
		expected.Port == actual.Port &&
		CompareStringArrays(expected.Addresses, actual.Addresses)
}

func ServicesMapMatch(expected map[string]ServiceDef, actual map[string]ServiceDef) bool {
	//for serviceName, leftServiceDef := range expected {
	//	rightServiceDef, found := actual[serviceName]
	//	if !found || !ServicesMatch(leftServiceDef, rightServiceDef) {
	//		return false
	//	}
	//}
	//
	//actual.

	return reflect.DeepEqual(expected, actual)
}
