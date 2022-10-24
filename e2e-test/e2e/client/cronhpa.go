package client

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	autoscalingv1beta1 "github.com/AliyunContainerService/kubernetes-cronhpa-controller/pkg/apis/autoscaling/v1beta1"
	"github.com/AliyunContainerService/kubernetes-cronhpa-controller/pkg/client/clientset/versioned"
)

func NewCronhpaSchedule(name, namespace string, targetReplicas int32) *autoscalingv1beta1.CronHorizontalPodAutoscaler {
	return &autoscalingv1beta1.CronHorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: autoscalingv1beta1.CronHorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv1beta1.ScaleTargetRef{
				ApiVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       name,
			},
			Jobs: []autoscalingv1beta1.Job{
				{
					Name:       "scale",
					Schedule:   "*/10 * * * * *",
					TargetSize: targetReplicas,
				},
			},
		},
	}
}

func CreateCronhpaSchedule(cli *versioned.Clientset, name, namespace string, targetReplicas int32) (*autoscalingv1beta1.CronHorizontalPodAutoscaler, error) {
	cronhpa := NewCronhpaSchedule(name, namespace, targetReplicas)
	res, err := cli.AutoscalingV1beta1().CronHorizontalPodAutoscalers(namespace).Create(context.Background(), cronhpa, metav1.CreateOptions{})
	return res, err
}

func NewCronhpaDate(name, namespace string, targetReplicas int32) *autoscalingv1beta1.CronHorizontalPodAutoscaler {
	later := time.Now().Add(10 * time.Second)
	schedule := "@date " + later.Format("2006-01-02 15:04:05")
	return &autoscalingv1beta1.CronHorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: autoscalingv1beta1.CronHorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv1beta1.ScaleTargetRef{
				ApiVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       name,
			},
			Jobs: []autoscalingv1beta1.Job{
				{
					Name:       name,
					Schedule:   schedule,
					TargetSize: targetReplicas,
				},
			},
		},
	}
}

func CreateCronhpaDate(cli *versioned.Clientset, name, namespace string, targetReplicas int32) (*autoscalingv1beta1.CronHorizontalPodAutoscaler, error) {
	cronhpa := NewCronhpaDate(name, namespace, targetReplicas)
	res, err := cli.AutoscalingV1beta1().CronHorizontalPodAutoscalers(namespace).Create(context.Background(), cronhpa, metav1.CreateOptions{})
	return res, err
}

func NewCronhpaExcludeToday(name, namespace string, targetReplicas int32) *autoscalingv1beta1.CronHorizontalPodAutoscaler {
	now := time.Now()
	month, day := now.Month(), now.Day()
	return &autoscalingv1beta1.CronHorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: autoscalingv1beta1.CronHorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv1beta1.ScaleTargetRef{
				ApiVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       name,
			},
			ExcludeDates: []string{
				fmt.Sprintf("* * * %d %d *", day, month),
			},
			Jobs: []autoscalingv1beta1.Job{
				{
					Name:       name,
					Schedule:   "*/10 * * * * *",
					TargetSize: targetReplicas,
					RunOnce:    true,
				},
			},
		},
	}
}

func CreateCronhpaExcludeToday(cli *versioned.Clientset, name, namespace string, targetReplicas int32) (*autoscalingv1beta1.CronHorizontalPodAutoscaler, error) {
	cronhpa := NewCronhpaExcludeToday(name, namespace, targetReplicas)
	res, err := cli.AutoscalingV1beta1().CronHorizontalPodAutoscalers(namespace).Create(context.Background(), cronhpa, metav1.CreateOptions{})
	return res, err
}

func NewCronhpaRunonce(name, namespace string, targetReplicas int32) *autoscalingv1beta1.CronHorizontalPodAutoscaler {
	return &autoscalingv1beta1.CronHorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: autoscalingv1beta1.CronHorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv1beta1.ScaleTargetRef{
				ApiVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       name,
			},
			Jobs: []autoscalingv1beta1.Job{
				{
					Name:       name,
					Schedule:   "*/10 * * * * *",
					TargetSize: targetReplicas,
					RunOnce:    true,
				},
			},
		},
	}
}

func CreateCronhpaRunonce(cli *versioned.Clientset, name, namespace string, targetReplicas int32) (*autoscalingv1beta1.CronHorizontalPodAutoscaler, error) {
	cronhpa := NewCronhpaRunonce(name, namespace, targetReplicas)
	res, err := cli.AutoscalingV1beta1().CronHorizontalPodAutoscalers(namespace).Create(context.Background(), cronhpa, metav1.CreateOptions{})
	return res, err
}

func NewCronhpaSts(name, namespace string, targetReplicas int32) *autoscalingv1beta1.CronHorizontalPodAutoscaler {
	return &autoscalingv1beta1.CronHorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: autoscalingv1beta1.CronHorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv1beta1.ScaleTargetRef{
				ApiVersion: "apps/v1",
				Kind:       "StatefulSet",
				Name:       name,
			},
			Jobs: []autoscalingv1beta1.Job{
				{
					Name:       name,
					Schedule:   "*/10 * * * * *",
					TargetSize: targetReplicas,
				},
			},
		},
	}
}

func CreateCronhpaSts(cli *versioned.Clientset, name, namespace string, targetReplicas int32) (*autoscalingv1beta1.CronHorizontalPodAutoscaler, error) {
	cronhpa := NewCronhpaSts(name, namespace, targetReplicas)
	res, err := cli.AutoscalingV1beta1().CronHorizontalPodAutoscalers(namespace).Create(context.Background(), cronhpa, metav1.CreateOptions{})
	return res, err
}

func NewCronhpaHpa(name, namespace string, targetReplicas int32) *autoscalingv1beta1.CronHorizontalPodAutoscaler {
	return &autoscalingv1beta1.CronHorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: autoscalingv1beta1.CronHorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv1beta1.ScaleTargetRef{
				ApiVersion: "autoscaling/v2beta2",
				Kind:       "HorizontalPodAutoscaler",
				Name:       name,
			},
			Jobs: []autoscalingv1beta1.Job{
				{
					Name:       name,
					Schedule:   "*/10 * * * * *",
					TargetSize: targetReplicas,
				},
			},
		},
	}
}

func CreateCronhpaHpa(cli *versioned.Clientset, name, namespace string, targetReplicas int32) (*autoscalingv1beta1.CronHorizontalPodAutoscaler, error) {
	cronhpa := NewCronhpaHpa(name, namespace, targetReplicas)
	res, err := cli.AutoscalingV1beta1().CronHorizontalPodAutoscalers(namespace).Create(context.Background(), cronhpa, metav1.CreateOptions{})
	return res, err
}

func GetCronhpa(cli *versioned.Clientset, name, namespace string) (*autoscalingv1beta1.CronHorizontalPodAutoscaler, error) {
	res, err := cli.AutoscalingV1beta1().CronHorizontalPodAutoscalers(namespace).Get(context.Background(), name, metav1.GetOptions{})
	return res, err
}

func DeleteCronhpa(cli *versioned.Clientset, name, namespace string) error {
	err := cli.AutoscalingV1beta1().CronHorizontalPodAutoscalers(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	return err
}
