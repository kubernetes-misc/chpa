package model

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var CronHPAV1CRDSchema = schema.GroupVersionResource{
	Group:    "kubernetes-misc.xyz",
	Version:  "v1",
	Resource: "cronhpas",
}

type CronHPAV1 struct {
	Metadata MetadataV1 `json:"metadata"`
	Spec     SpecV1     `json:"spec"`
}

func (cs CronHPAV1) PrettyString() string {
	return fmt.Sprintf("%s replicas: %v ==> %v @ CPU load of %v%% (cronscale/%s operating on %s/%s)", pad(cs.Spec.CronSpec, 12), cs.Spec.HorizontalPodAutoScaler.MinReplicas, cs.Spec.HorizontalPodAutoScaler.MaxReplicas, cs.Spec.HorizontalPodAutoScaler.TargetCPUUtilizationPercentage, cs.Metadata.Name, cs.Spec.ScaleTargetRef.Kind, cs.Spec.ScaleTargetRef.Name)

}

func (cs CronHPAV1) GetID() string {
	return "chpaV1." + cs.Metadata.Namespace + "." + cs.Metadata.Name
}

type MetadataV1 struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type SpecV1 struct {
	CronSpec                string                  `json:"cronSpec"`
	ScaleTargetRef          ScaleTargetRefV1        `json:"scaleTargetRef"`
	HorizontalPodAutoScaler HorizontalPodAutoScaler `json:"horizontalPodAutoScaler"`
}

type ScaleTargetRefV1 struct {
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Name       string `json:"name"`
}

type HorizontalPodAutoScaler struct {
	Name                           string `json:"name"`
	MaxReplicas                    int32  `json:"maxReplicas"`
	MinReplicas                    int32  `json:"minReplicas"`
	TargetCPUUtilizationPercentage int32  `json:"targetCPUUtilizationPercentage"`
}

func pad(in string, size int) string {
	result := in
	for len(result) < size {
		result += " "
	}
	return result
}
