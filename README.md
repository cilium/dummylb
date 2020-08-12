# dummylb

A dummy k8s loadbalancer which assigns a user requested IP addr to a LoadBalancer
Service.

Used only for testing.

```
make compile
docker build -t cilium/dummylb:$VSN .
```
