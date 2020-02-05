package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"k8s.io/mouse/client"
	"k8s.io/mouse/model"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

func main() {
	err := client.BuildClient()
	if err != nil {
		panic(err)
	}

	c, err := client.GetAllCRD("default", model.CronScaleV1CRDSchema)
	if err != nil {
		logrus.Errorln(err)
		return
	}

	for cs := range c {
		fmt.Println(fmt.Sprintf("%s replicas: %v ==> %v @ CPU load of %v%% (cronscale/%s operating on %s/%s)", pad(cs.Spec.CronSpec, 12), cs.Spec.HorizontalPodAutoScaler.MinReplicas, cs.Spec.HorizontalPodAutoScaler.MaxReplicas, cs.Spec.HorizontalPodAutoScaler.TargetCPUUtilizationPercentage, cs.Metadata.Name, cs.Spec.ScaleTargetRef.Kind, cs.Spec.ScaleTargetRef.Name))
		dep, err := client.GetDeployment(cs.Metadata.Namespace, cs.Spec.ScaleTargetRef.Name)
		if err != nil {
			logrus.Errorln(err)
			continue
		}
		dep.Spec.Replicas = &cs.Spec.HorizontalPodAutoScaler.MinReplicas
		client.UpdateDeployment(dep)

		hpa, err := client.GetHPA(cs.Metadata.Namespace, cs.Spec.HorizontalPodAutoScaler.Name)
		if err != nil {
			logrus.Errorln(err)
			continue
		}
		hpa.Spec.MinReplicas = &cs.Spec.HorizontalPodAutoScaler.MinReplicas
		hpa.Spec.MaxReplicas = cs.Spec.HorizontalPodAutoScaler.MaxReplicas
		hpa.Spec.TargetCPUUtilizationPercentage = &cs.Spec.HorizontalPodAutoScaler.TargetCPUUtilizationPercentage
		client.UpdateHPA(cs.Metadata.Namespace, hpa)

		if true {
			return
		}

	}

}

func pad(in string, size int) string {
	result := in
	for len(result) < size {
		result += " "
	}
	return result
}
