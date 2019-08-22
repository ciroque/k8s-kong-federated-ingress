# k8s-kong-federated-ingress

A custom Kubernetes Controller that watches Ingress resources and updates an external Kong. 

## References

This project is based on the article and sample code from:

[Extending Kubernetes: Create Controllers for Core and Custom Resources](https://medium.com/@trstringer/create-kubernetes-controllers-for-core-and-custom-resources-62fc35ad64a3)

[Base sample for a custom controller in Kubernetes working with core resources](https://github.com/trstringer/k8s-controller-core-resource)

## Questions

- Retry semantics for failed creation of King Kong objects, are there samples in k8s land?
- Unit testing?

