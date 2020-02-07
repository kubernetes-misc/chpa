# kubernetes-misc

## Cron-based Horizontal Pod Autoscaler

Trigger different Horizontal Pod Autoscaling strategies by cron 

## Roadmap

### Version 1
- Support for cron with seconds and without
- Official Docker images
- Docs
- Yaml for Kubernetes deployment
- HA: interim strategy involving redundant imperative approach
- HA: Integration with etcd to allow for HA deployment and locking when carrying out jobs

### Version 2
- More options around how to kube conf - using in-cluster authentication, RBAC instead
- Investigate using selectors to pick HPAs - Kevin
- Investigate pointing to the official HPA API client struct
- Investigate different concurrency model for increased reliability in various failure cases 
- Investigate approaching as an Operator
