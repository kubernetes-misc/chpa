# kubernetes-misc

## Cron-based Horizontal Pod Autoscaler

Trigger different Horizontal Pod Autoscaling strategies by cron 

## Roadmap
- Support for cron with seconds and without
- Official Docker images
- Docs
- HA: interim strategy involving redundant imperative approach
- HA: Integration with etcd to allow for HA deployment and locking when carrying out jobs
- Investigate approaching as an Operator
- Investigate pointing to the official HPA API client struct
- Investigate different concurrency model for increased reliability in various failure cases 
- More options around how to handle .kube/conf as file using in-cluster authentication
- Investigate using selectors instead - Kevin
