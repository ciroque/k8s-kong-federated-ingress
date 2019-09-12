/*

	Package eventing implements the interfaces required to subscribe and respond to Kubernetes API resources.

	k8s-kong-federated-ingress is interested only in the Ingress resource.

	- Controller manages the work queue that is populated by the Kubernetes Event Handlers registered in main().

	- ApiHandler manges translation of K8s resources to an intermediate representation, then onto a Kong=specific
		representation, then onto the Registrar for creating the resources in Kong.

	- Event is a representation of the K8s API event that is stored in the work queue for subsequent processing.

*/
package eventing
