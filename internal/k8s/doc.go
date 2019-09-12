/*

	Package k8s contains the various types and functions to map K8s-specific resources into an intermediate format
	that allows easier manipulation within the k9s-kong-federated-ingress domain.

	- ServiceDef is used to extract the pertinent data from the K8s resources.

	- Translation contains the logic to translate the K8s resources to the k8s-kong-federated-ingress domain objects.

*/
package k8s
