package client

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/pointer"
)

func NewDeployment(name, namespace, image string, replicas int32) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app": name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32Ptr(replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image:           image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Name:            name,
						},
					},
				},
			},
		},
	}
}

func CreateDeployment(cli kubernetes.Interface, name, namespace, image string, replicas int32) (*appsv1.Deployment, error) {
	deploy := NewDeployment(name, namespace, image, replicas)
	result, err := cli.AppsV1().Deployments(namespace).Create(context.Background(), deploy, metav1.CreateOptions{})
	return result, err
}

func GetDeployment(cli kubernetes.Interface, name, namespace string) (*appsv1.Deployment, error) {
	deploy, err := cli.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	return deploy, err
}

func PatchDeployReplicas(cli kubernetes.Interface, name, namespace string, newReplicas int32) (*appsv1.Deployment, error) {
	var patch = fmt.Sprintf(`{"spec": {"replicas": %d}}`, newReplicas)
	res, err := cli.AppsV1().Deployments(namespace).
		Patch(context.Background(), name, types.StrategicMergePatchType, []byte(patch), metav1.PatchOptions{})
	return res, err
}

func DeleteDeployment(cli kubernetes.Interface, name, namespace string) error {
	err := cli.AppsV1().Deployments(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	return err
}

func WaitDeployReplicas(cli kubernetes.Interface, name, namespace string, replicas int32) error {
	return wait.Poll(2*time.Second, 2*time.Minute, func() (done bool, err error) {
		d, err := cli.AppsV1().Deployments(namespace).Get(
			context.Background(),
			name,
			metav1.GetOptions{},
		)
		if err != nil {
			return false, err
		}
		if d.Status.ReadyReplicas != replicas {
			return false, nil
		}
		if d.Status.AvailableReplicas != replicas {
			return false, nil
		}
		if d.Status.UpdatedReplicas != replicas {
			return false, nil
		}
		return true, nil
	})
}
