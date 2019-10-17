# k8s-kong-federated-ingress

A custom Kubernetes Controller that watches Ingress resources and updates an external Kong. 

## References

This project is based on the article and sample code from:

[Extending Kubernetes: Create Controllers for Core and Custom Resources](https://medium.com/@trstringer/create-kubernetes-controllers-for-core-and-custom-resources-62fc35ad64a3)

[Base sample for a custom controller in Kubernetes working with core resources](https://github.com/trstringer/k8s-controller-core-resource)

## Questions


## Running Locally

1. Clone the repository
1. Define the required Environment Variables:
    1. KUBERNETES_SERVICE_HOST: The host for the k8s API
        (for local dev environments this is obtained using `minikube ip`, or finding the cluster's ip in the Rancher console)
    1. KUBERNETES_SERVICE_PORT: The port for the k8s API
        (this is most often ***6443***)
    1. KONG_HOST: The Kong host that is the target for Ingress synchronization
        (For local dev environments this is most often ***http://localhost:8001***)
1. Ensure the Service Account's token and ca.crt are available at the following locations:
    1. /var/run/secrets/kubernetes.io/serviceaccount/token
    1. /var/run/secrets/kubernetes.io/serviceaccount/ca.crt

### Creating the RBAC entities and export the token and ca.crt to the proper locations

1. `kubectl apply -f helm/k8s-kong-federated-ingress/templates/rbac-config.yaml`
1. `sudo hack/create-rbac-entities`

You should now be able to run the project.

### With IntelliJ
1. Ensure you have the Go plugin installed
1. Create a new Run/Debug Configuration using the *Go Build* base configuration
1. Set the Files to _/site/k8s-kong-federated-ingress/cmd/k8s-kong-federated-ingress/main.go_
1. Ensure *Run after build* is checked
1. The Environment Variables mentioned above can be defined in this dialog if you prefer
1. Click the green Run button, or use Shift + F10

### From terminal
- `go build -o bin/k8s-kong-federated-ingress cmd/k8s-kong-federated-ingress/main.go && bin/k8s-kong-federated-ingress`

_ -- OR -- _

- `go run cmd/k8s-kong-federated-ingress/main.go`

## Running tests

`go test ./...`
