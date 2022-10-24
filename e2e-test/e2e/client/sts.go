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

func NewStatefulset(name, namespace, image string, replicas int32) *appsv1.StatefulSet {
	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: name,
			Replicas:    pointer.Int32Ptr(replicas),
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
							Name:            name,
							ImagePullPolicy: corev1.PullIfNotPresent,
						},
					},
				},
			},
		},
	}
}

func CreateStatefulset(cli kubernetes.Interface, name, namespace, image string, replicas int32) (*appsv1.StatefulSet, error) {
	sts := NewStatefulset(name, namespace, image, replicas)
	res, err := cli.AppsV1().StatefulSets(namespace).Create(context.Background(), sts, metav1.CreateOptions{})
	return res, err
}

func WaitStatefulsetReplicas(cli kubernetes.Interface, name, namespace string, replicas int32) error {
	return wait.Poll(2*time.Second, 2*time.Minute, func() (done bool, err error) {
		sts, err := cli.AppsV1().StatefulSets(namespace).Get(
			context.Background(),
			name,
			metav1.GetOptions{},
		)
		if err != nil {
			return false, nil
		}
		if sts.Status.ReadyReplicas != *sts.Spec.Replicas {
			return false, nil
		}
		if sts.Status.UpdatedReplicas != *sts.Spec.Replicas {
			return false, nil
		}
		return true, nil
	})
}

func DeleteStatefulset(cli kubernetes.Interface, name, namespace string) error {
	err := cli.AppsV1().StatefulSets(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	return err
}

func PatchStsReplicas(cli kubernetes.Interface, name, namespace string, replicas int32) (*appsv1.StatefulSet, error) {
	var patch = fmt.Sprintf(`{"spec": {"replicas": %d}}`, replicas)
	res, err := cli.AppsV1().StatefulSets(namespace).
		Patch(context.Background(), name, types.StrategicMergePatchType, []byte(patch), metav1.PatchOptions{})
	return res, err
}

func GetStatefulset(cli kubernetes.Interface, name, namespace string) (*appsv1.StatefulSet, error) {
	res, err := cli.AppsV1().StatefulSets(namespace).Get(context.Background(), name, metav1.GetOptions{})
	return res, err
}
