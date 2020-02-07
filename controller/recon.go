package controller

import (
	"fmt"
	"github.com/kubernetes-misc/chpa/client"
	"github.com/kubernetes-misc/chpa/model"
	"github.com/sirupsen/logrus"
)

var ReconHub = NewReconHub()

func NewReconHub() *reconHub {
	r := &reconHub{in: make(chan model.CronHPAV1, 256)}
	go func() {
		for cs := range r.in {
			logrus.Debugln("recon hub has received", cs.GetID(), "event")
			checkAndUpdate(cs)
		}
	}()
	return r
}

type reconHub struct {
	in chan model.CronHPAV1
}

func (r *reconHub) Add(cs model.CronHPAV1) {
	r.in <- cs
}

func checkAndUpdate(cs model.CronHPAV1) {

	checkAndUpdateDeployment(cs)
	checkAndUpdateHPA(cs)

}

func checkAndUpdateHPA(cs model.CronHPAV1) {
	//Check the hpa
	hpa, err := client.GetHPA(cs.Metadata.Namespace, cs.Spec.HorizontalPodAutoScaler.Name)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	hpa.Spec.MinReplicas = &cs.Spec.HorizontalPodAutoScaler.MinReplicas
	hpa.Spec.MaxReplicas = cs.Spec.HorizontalPodAutoScaler.MaxReplicas
	hpa.Spec.TargetCPUUtilizationPercentage = &cs.Spec.HorizontalPodAutoScaler.TargetCPUUtilizationPercentage
	client.UpdateHPA(cs.Metadata.Namespace, hpa)
	logrus.Infoln(fmt.Sprintf(">> Updating hpa/%s from %v to %v @ CPU load of %v%%", cs.Spec.HorizontalPodAutoScaler.Name, cs.Spec.HorizontalPodAutoScaler.MinReplicas, cs.Spec.HorizontalPodAutoScaler.MaxReplicas, cs.Spec.HorizontalPodAutoScaler.TargetCPUUtilizationPercentage))

}

func checkAndUpdateDeployment(cs model.CronHPAV1) {
	//Check the deployment
	dep, err := client.GetDeployment(cs.Metadata.Namespace, cs.Spec.ScaleTargetRef.Name)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	dep.Spec.Replicas = &cs.Spec.HorizontalPodAutoScaler.MinReplicas
	client.UpdateDeployment(dep)
	logrus.Infoln(fmt.Sprintf(">> Updating deployment/%s to min replicas %v", cs.Spec.HorizontalPodAutoScaler.Name, cs.Spec.HorizontalPodAutoScaler.MinReplicas))

}
