FROM toolsnexus.marchex.com:5000/alpine:3.4

RUN apk add --no-cache bash jq curl libc6-compat
WORKDIR /bin
COPY ./bin/k8s-kong-federated-ingress ./
ENTRYPOINT ["/bin/k8s-kong-federated-ingress"]
