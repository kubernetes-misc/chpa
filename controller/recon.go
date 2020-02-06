package controller

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"k8s.io/mouse/client"
	"k8s.io/mouse/model"
)

var ReconHub = NewReconHub()

func NewReconHub() *reconHub {
	r := &reconHub{in: make(chan model.CronScaleV1, 256)}
	go func() {
		for cs := range r.in {
			checkAndUpdate(cs)
		}
	}()
	return r
}

type reconHub struct {
	in chan model.CronScaleV1
}

func (r *reconHub) Add(cs model.CronScaleV1) {
	r.in <- cs
}

func checkAndUpdate(cs model.CronScaleV1) {
	logrus.Println(fmt.Sprintf("%s replicas: %v ==> %v @ CPU load of %v%% (cronscale/%s operating on %s/%s)", pad(cs.Spec.CronSpec, 12), cs.Spec.HorizontalPodAutoScaler.MinReplicas, cs.Spec.HorizontalPodAutoScaler.MaxReplicas, cs.Spec.HorizontalPodAutoScaler.TargetCPUUtilizationPercentage, cs.Metadata.Name, cs.Spec.ScaleTargetRef.Kind, cs.Spec.ScaleTargetRef.Name))

	checkAndUpdateDeployment(cs)
	checkAndUpdateHPA(cs)

}

func checkAndUpdateHPA(cs model.CronScaleV1) {
	//Check the hpa
	hpa, err := client.GetHPA(cs.Metadata.Namespace, cs.Spec.HorizontalPodAutoScaler.Name)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	//TODO: check if we need to change at all?
	hpa.Spec.MinReplicas = &cs.Spec.HorizontalPodAutoScaler.MinReplicas
	hpa.Spec.MaxReplicas = cs.Spec.HorizontalPodAutoScaler.MaxReplicas
	hpa.Spec.TargetCPUUtilizationPercentage = &cs.Spec.HorizontalPodAutoScaler.TargetCPUUtilizationPercentage
	client.UpdateHPA(cs.Metadata.Namespace, hpa)

}

func checkAndUpdateDeployment(cs model.CronScaleV1) {
	//Check the deployment
	dep, err := client.GetDeployment(cs.Metadata.Namespace, cs.Spec.ScaleTargetRef.Name)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	//TODO: check if we need to change at all?
	dep.Spec.Replicas = &cs.Spec.HorizontalPodAutoScaler.MinReplicas
	client.UpdateDeployment(dep)

}

func pad(in string, size int) string {
	result := in
	for len(result) < size {
		result += " "
	}
	return result
}
