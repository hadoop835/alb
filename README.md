# ALB -- Another Load Balancer

ALB (Another Load Balancer) is a Kubernetes Gateway powered by [OpenResty](https://github.com/openresty/) with years of production experience from Alauda.

> *Note*: We are in the process of preparing the necessary documentation and refactoring the code for open source. More information and detailed usage will be made available soon.

## Advantages

- **Isolation and Multi-Tenant**: With ALB operator, multiple ALB instances can be created and managed in one cluster. Each tenant can has a group of dedicated ALB instances.
- **Ingress and Gateway API Support**: Users can flexibly choose between Ingress and Gateway API according to their own preferences.
- **Flexible User Defined Traffic Rule**: ALB provides a traffic rule DSL that can support more complex traffic matching and distribution scenarios that beyond the capabilities of standard Ingress and Gateway API.
- **Multiple Protocol Support**: ALB can manage HTTP, HTTPS, TCP and UDP traffic.

## Architecture

![](docs/architecture.png)

## Quick Start

### Deploy the ALB Operator

1. Create a kind cluster
2. `helm repo add alb https://alauda.github.io/alb/;helm repo update;helm search repo|grep alb`
3. `helm install alb-operator alb/alauda-alb2` 

### Deploy an ALB Instance

```bash
cat <<EOF | kubectl apply -f -
apiVersion: crd.alauda.io/v2beta1
kind: ALB2
metadata:
    name: alb-demo
    namespace: kube-system
spec:
    address: "172.20.0.5"  # the ip address of node where alb been deployed
    type: "nginx" 
    config:
        networkMode: host
        loadbalancerName: alb-demo
        nodeSelector:
          alb-demo: "true"
        projects:
        - ALL_ALL
        replicas: 1
EOF
```

### Rua a Demo Application

```bash
cat <<EOF | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello-world
  labels:
    k8s-app: hello-world
spec:
  replicas: 1 
  selector:
    matchLabels:
      k8s-app: hello-world
  template:
    metadata:
      labels:
        k8s-app: hello-world
    spec:
      terminationGracePeriodSeconds: 60
      containers:
      - name: hello-world
        image: docker.io/crccheck/hello-world:latest 
        imagePullPolicy: IfNotPresent
---
apiVersion: v1
kind: Service
metadata:
  name: hello-world
  labels:
    k8s-app: hello-world
spec:
  ports:
  - name: http
    port: 80
    targetPort: 8000
  selector:
    k8s-app: hello-world 
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: hello-world
spec:
  rules:
  - http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: hello-world
            port:
              number: 80
EOF
```

Now you visit the app by `curl http://${ip}`

## Advanced Features

### `Frontend` and `Rule`

Complex traffic matching and distribution patterns can be configured by `Frontend` and `Rule`.

```yaml
apiVersion: crd.alauda.io/v1
kind: Frontend
metadata:
  labels:
    alb2.cpaas.io/name: alb-demo # required, indicate the ALB instance to which this Frontend belongs to
  name: alb-demo-00080
  namespace: kube-system
spec:
  backendProtocol: ""   # http|https 
  certificate_name: ""  # $secret_ns/$secret_name
  port: 80              
  protocol: http        # protocol of this Frontend itself
---
apiVersion: crd.alauda.io/v1
kind: Rule
metadata:
  labels:
    alb2.cpaas.io/frontend: alb-demo-00080  # required, indicate the Frontend to which this rule belongs
    alb2.cpaas.io/name: alb-demo            # required, indicate the ALB to which this rule belongs
  name: alb-demo-00080-topu
  namespace: kube-system
spec:
  backendProtocol: ""                       # as same as Frontend
  certificate_name: ""                      # as same as Frontend
  dslx:                                     # matching rule DSL
  - type: URL
    values:
    - - STARTS_WITH
      - /
  enableCORS: false
  priority: 5                              # the lower the number, the higher the priority
  serviceGroup:
    services:
    - name: hello-world
      namespace: default
      port: 80
      weight: 100
```

### Ingress Annotations

#### Rewrite Request

```yaml
alb.ingress.cpaas.io/rewrite-request: |
{"headers_remove":["h1"],"headers":{"a":"b"},"headers_add":{"aa": ["bb","cc"]}}
```

`headers_remove`: remove the header.
`headers_add`: append to the header instead of overwrite it.
`headers`: set the header.

#### Rewrite Response

```yaml
alb.ingress.cpaas.io/rewrite-response: |
{"headers_remove":["h1"],"headers":{"a":"b"},"headers_add":{"aa": ["bb","cc"]}}
```

`headers_remove`: remove the header.
`headers_add`: append to the header instead of overwrite it.
`headers`: set the header.

#### Annotations Compatible with ingress-nginx

```yaml
nginx.ingress.kubernetes.io/rewrite-target
nginx.ingress.kubernetes.io/enable-cors
nginx.ingress.kubernetes.io/cors-allow-headers
nginx.ingress.kubernetes.io/cors-allow-origin
nginx.ingress.kubernetes.io/backend-protocol
nginx.ingress.kubernetes.io/temporal-redirect
nginx.ingress.kubernetes.io/permanent-redirect
nginx.ingress.kubernetes.io/upstream-vhost
```

### Container Network

By default, ALB is deployed as a host network, which has the advantage of direct access via node ip, and the disadvantage that each ALB can only have exclusive access to the node, or you need to manually manage the ALB's ports.
But, ALB also supports container network mode deployment and provides external access through Loadbalancer type Service.

```yaml
apiVersion: crd.alauda.io/v2beta1
kind: ALB2
metadata:
    name: alb-demo
    namespace: kube-system
spec:
    type: "nginx" 
    config:
        networkMode: container           # use container networkmode
        vip:
            enableLbSvc: true            # automatically creates a Service of type LoadBalancer and treats the address assigned to the Service as the address of the ALB
        loadbalancerName: alb-demo
        nodeSelector:
          alb-demo: "true"
        projects:
        - ALL_ALL
        replicas: 1
```

### Gateway API

ALB supports Gateway out of box, just set the `gatewayClassName` to `exclusive-gateway` when creating gateways.

```yaml
apiVersion: gateway.networking.k8s.io/v1alpha2
kind: Gateway
metadata:
  name: g1 
  namespace: g1 
spec:
  gatewayClassName:  exclusive-gateway
  listeners:
  - name: http
    port: 80
    protocol: HTTP
    allowedRoutes:
      namespaces:
        from: All
```