# kubernetes-cronscale
Trigger different Horizontal Pod Autoscaling strategies by cron 

## To do
- Pick a name for the project
- Pick a home for the project?
- Change the crd spec.group to new project home?
- Integration with Etcd to allow for HA deployment and locking when carrying out jobs
- Consider approaching as an Operator
- Consider pointing to the official HPA API client struct
- Consider different concurrency model for increased reliability in various failure cases 
