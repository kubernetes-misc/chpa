package model

import "k8s.io/apimachinery/pkg/runtime/schema"

var CronScaleV1CRDSchema = schema.GroupVersionResource{
	Group:    "captainjustin.space",
	Version:  "v1",
	Resource: "cronscales",
}

type CronScaleV1 struct {
	Metadata MetadataV1 `json:"metadata"`
	Spec     SpecV1     `json:"spec"`
}

type MetadataV1 struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type SpecV1 struct {
	CronSpec                       string           `json:"cronSpec"`
	MaxReplicas                    int32            `json:"maxReplicas"`
	MinReplicas                    int32            `json:"minReplicas"`
	TargetCPUUtilizationPercentage int              `json:"targetCPUUtilizationPercentage"`
	ScaleTargetRef                 ScaleTargetRefV1 `json:"scaleTargetRef"`
}

type ScaleTargetRefV1 struct {
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Name       string `json:"name"`
}
