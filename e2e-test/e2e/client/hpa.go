package client

import (
	"context"

	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/pointer"
)

func NewHpa(name, namespace, metricName string, minReplicas, maxReplicas int32) *autoscalingv2beta2.HorizontalPodAutoscaler {
	avgValue := resource.MustParse("1000m")
	return &autoscalingv2beta2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: autoscalingv2beta2.HorizontalPodAutoscalerSpec{
			MaxReplicas: maxReplicas,
			MinReplicas: pointer.Int32Ptr(minReplicas),
			ScaleTargetRef: autoscalingv2beta2.CrossVersionObjectReference{
				Kind:       "Deployment",
				Name:       name,
				APIVersion: "apps/v1",
			},
			Metrics: []autoscalingv2beta2.MetricSpec{
				{
					Type: "Pods",
					Pods: &autoscalingv2beta2.PodsMetricSource{
						Metric: autoscalingv2beta2.MetricIdentifier{
							Name: metricName,
						},
						Target: autoscalingv2beta2.MetricTarget{
							Type:         "AverageValue",
							AverageValue: &avgValue,
						},
					},
				},
			},
		},
	}
}

func CreateHpa(cli kubernetes.Interface, name, namespace, metricName string, minReplicas, maxReplicas int32) (*autoscalingv2beta2.HorizontalPodAutoscaler, error) {
	hpa := NewHpa(name, namespace, metricName, minReplicas, maxReplicas)
	res, err := cli.AutoscalingV2beta2().HorizontalPodAutoscalers(namespace).Create(context.Background(), hpa, metav1.CreateOptions{})
	return res, err
}

func GetHpa(cli kubernetes.Interface, name, namespace string) (*autoscalingv2beta2.HorizontalPodAutoscaler, error) {
	res, err := cli.AutoscalingV2beta2().HorizontalPodAutoscalers(namespace).Get(context.Background(), name, metav1.GetOptions{})
	return res, err
}

func DeleteHpa(cli kubernetes.Interface, name, namespace string) error {
	err := cli.AutoscalingV2beta2().HorizontalPodAutoscalers(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	return err
}

func UpdateHpaStatus(cli kubernetes.Interface, name, namespace string, currentReplicas int32) (*autoscalingv2beta2.HorizontalPodAutoscaler, error) {
	hpa, err := cli.AutoscalingV2beta2().HorizontalPodAutoscalers(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	hpa.Status.CurrentReplicas = currentReplicas
	hpa, err = cli.AutoscalingV2beta2().HorizontalPodAutoscalers(namespace).UpdateStatus(context.Background(), hpa, metav1.UpdateOptions{})
	return hpa, err
}