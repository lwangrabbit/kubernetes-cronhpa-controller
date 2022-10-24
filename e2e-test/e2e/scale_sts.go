package e2e

import (
	"time"

	. "github.com/onsi/ginkgo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/kubernetes/test/e2e/framework"

	"github.com/AliyunContainerService/kubernetes-cronhpa-controller/e2e-test/e2e/client"
	e2econfig "github.com/AliyunContainerService/kubernetes-cronhpa-controller/e2e-test/e2e/config"
	"github.com/AliyunContainerService/kubernetes-cronhpa-controller/pkg/client/clientset/versioned"
)

var _ = Describe("Statefulset with schedule", func() {

	name := "test-sts"
	namespace := metav1.NamespaceDefault
	image := e2econfig.TestConfig.Image

	var cli kubernetes.Interface
	var cronhpaCli *versioned.Clientset

	var err error
	cli, err = framework.LoadClientset()
	framework.ExpectNoError(err, "failed to load client")
	var config *rest.Config
	config, err = framework.LoadConfig()
	framework.ExpectNoError(err, "faild to load config")
	cronhpaCli, err = versioned.NewForConfig(config)
	framework.ExpectNoError(err, "failed to create cronhpa client")

	Context("schedule to scale up", func() {
		var initReplicas int32 = 2
		var targetSize int32 = 3

		BeforeEach(func() {
			err := PrepareScaleSts(cli, name, namespace, image, initReplicas)
			framework.ExpectNoError(err, "failed to prepare sts")
			_, err = client.CreateCronhpaSts(cronhpaCli, name, namespace, targetSize)
			framework.ExpectNoError(err, "failed to create cronhpa")
		})
		AfterEach(func() {
			err := client.DeleteCronhpa(cronhpaCli, name, namespace)
			framework.ExpectNoError(err, "faild to delete cronhpa")
			err = CleanScaleSts(cli, name, namespace)
			framework.ExpectNoError(err, "failed to delete sts")
		})

		It("check replicas ok", func() {
			err := wait.Poll(2*time.Second, 1*time.Minute, func() (done bool, err error) {
				sts, err := client.GetStatefulset(cli, name, namespace)
				if err != nil {
					return false, err
				}
				if sts.Status.ReadyReplicas != targetSize {
					return false, nil
				}
				return true, nil
			})
			framework.ExpectNoError(err, "failed to scale up deploy")
		})
	})

	Context("schedule to scale down", func() {
		var initReplicas int32 = 2
		var targetSize int32 = 1

		BeforeEach(func() {
			err := PrepareScaleSts(cli, name, namespace, image, initReplicas)
			framework.ExpectNoError(err, "failed to prepare sts")
			_, err = client.CreateCronhpaSts(cronhpaCli, name, namespace, targetSize)
			framework.ExpectNoError(err, "failed to create cronhpa")
		})
		AfterEach(func() {
			err := client.DeleteCronhpa(cronhpaCli, name, namespace)
			framework.ExpectNoError(err, "faild to delete cronhpa")
			err = CleanScaleSts(cli, name, namespace)
			framework.ExpectNoError(err, "faild to delete cronhpa")
		})

		It("check replicas ok", func() {
			err := wait.Poll(2*time.Second, 1*time.Minute, func() (done bool, err error) {
				sts, err := client.GetStatefulset(cli, name, namespace)
				if err != nil {
					return false, err
				}
				if sts.Status.ReadyReplicas != targetSize {
					return false, nil
				}
				return true, nil
			})
			framework.ExpectNoError(err, "failed to scale up sts")
		})
	})
})

func PrepareScaleSts(cli kubernetes.Interface, name, namespace, image string, replicas int32) error {
	_, err := client.CreateStatefulset(cli, name, namespace, image, replicas)
	if err != nil {
		return err
	}
	err = client.WaitStatefulsetReplicas(cli, name, namespace, replicas)
	return err
}

func CleanScaleSts(cli kubernetes.Interface, name, namespace string) error {
	err := client.DeleteStatefulset(cli, name, namespace)
	return err
}
