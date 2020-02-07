# CHPA
<img src="https://img.shields.io/badge/Version-alpha-f5bc42">&nbsp;<a href="https://goreportcard.com/report/github.com/kubernetes-misc/chpa"><img src="https://goreportcard.com/badge/github.com/kubernetes-misc/chpa"></a>&nbsp;<a href="https://codebeat.co/projects/github-com-kubernetes-misc-chpa-master"><img alt="codebeat badge" src="https://codebeat.co/badges/5c5ccd5a-a48b-400b-8400-f8cedfd93c63" /></a>&nbsp;<a href="https://codeclimate.com/github/kubernetes-misc/chpa/maintainability"><img src="https://api.codeclimate.com/v1/badges/0d70e5e60e9cdc89c9ff/maintainability" /></a>


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
- Docs (running, deployment and production)
- Yaml for Kubernetes deployment
- HA: interim strategy involving redundant imperative approach
- HA: Integration with etcd to allow for HA deployment and locking when carrying out jobs
- Investigate different concurrency model for increased reliability in various failure cases 

### Version 2
- More options around how to kube conf - using in-cluster authentication, RBAC instead
- Investigate using selectors to pick HPAs - Kevin
- Investigate approaching as an Operator
- Investigate pointing to the official HPA API client struct
- Exclude namespaces
