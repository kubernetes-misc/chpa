# CHPA

Declarative Cron-based Horizontal Pod Autoscaler for Kubernetes 

## Features

### Trigger different HPAs with cron jobs
- CRD for Kubernetes called `cronhpa` or `chpa`. See yaml/ directory.
- Specify a conventional cron or a cron with seconds.
- Specify different CPU load, min replicas and max replicas to be applied to a named horizontal pod autoscaler.
- Deployment `replicas` is updated to min replicas in the CRD.  
- All namespaces are searched for entries matching the CRD allowing multiple teams to make use of them.


## Roadmap

### Version 1
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
- Exclude namespaces
