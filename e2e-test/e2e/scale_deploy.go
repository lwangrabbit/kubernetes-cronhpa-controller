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

var _ = Describe("Deploy with schedule", func() {

	name := "test-deploy"
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
			err := PrepareScaleDeploy(cli, name, namespace, image, initReplicas)
			framework.ExpectNoError(err, "failed to patch deploy")
			_, err = client.CreateCronhpaSchedule(cronhpaCli, name, namespace, targetSize)
			framework.ExpectNoError(err, "faild to create cronhpa")
		})
		AfterEach(func() {
			err = CleanScaleDeploy(cli, name, namespace)
			framework.ExpectNoError(err, "failed to clean deploy")
			err := client.DeleteCronhpa(cronhpaCli, name, namespace)
			framework.ExpectNoError(err, "failed to delete cronhpa")
		})

		It("check replicas ok", func() {
			err := wait.Poll(2*time.Second, 1*time.Minute, func() (done bool, err error) {
				deploy, err := client.GetDeployment(cli, name, namespace)
				if err != nil {
					return false, err
				}
				if deploy.Status.ReadyReplicas != targetSize {
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
			err := PrepareScaleDeploy(cli, name, namespace, image, initReplicas)
			framework.ExpectNoError(err, "failed to patch deploy")
			_, err = client.CreateCronhpaSchedule(cronhpaCli, name, namespace, targetSize)
			framework.ExpectNoError(err, "faild to create cronhpa")
		})
		AfterEach(func() {
			err = CleanScaleDeploy(cli, name, namespace)
			framework.ExpectNoError(err, "failed to clean deploy")
			err := client.DeleteCronhpa(cronhpaCli, name, namespace)
			framework.ExpectNoError(err, "failed to delete cronhpa")
		})
		It("check replicas ok", func() {
			err := wait.Poll(2*time.Second, 1*time.Minute, func() (done bool, err error) {
				deploy, err := client.GetDeployment(cli, name, namespace)
				if err != nil {
					return false, err
				}
				if deploy.Status.ReadyReplicas != targetSize {
					return false, nil
				}
				return true, nil
			})
			framework.ExpectNoError(err, "failed to scale up deploy")
		})
	})

	Context("date to scale up", func() {
		var initReplicas int32 = 2
		var targetSize int32 = 3

		BeforeEach(func() {
			err := PrepareScaleDeploy(cli, name, namespace, image, initReplicas)
			framework.ExpectNoError(err, "failed to patch deploy")
			_, err = client.CreateCronhpaDate(cronhpaCli, name, namespace, targetSize)
			framework.ExpectNoError(err, "faild to create cronhpa")
		})
		AfterEach(func() {
			err = CleanScaleDeploy(cli, name, namespace)
			framework.ExpectNoError(err, "failed to clean deploy")
			err := client.DeleteCronhpa(cronhpaCli, name, namespace)
			framework.ExpectNoError(err, "failed to delete cronhpa")
		})

		It("check replicas ok", func() {
			// scale ok
			err := wait.Poll(2*time.Second, 1*time.Minute, func() (done bool, err error) {
				deploy, err := client.GetDeployment(cli, name, namespace)
				if err != nil {
					return false, err
				}
				if deploy.Status.ReadyReplicas != targetSize {
					return false, nil
				}
				return true, nil
			})
			framework.ExpectNoError(err, "failed to scale up deploy")
			// runonce: 1. scale deploy to 1
			var newReplicas int32 = 1
			_, err = client.PatchDeployReplicas(cli, name, namespace, newReplicas)
			framework.ExpectNoError(err, "failed to patch deploy")
			err = client.WaitDeployReplicas(cli, name, namespace, newReplicas)
			framework.ExpectNoError(err, "failed to wait deploy replicas")
			// runonce: 2. check not scale anymore
			err = wait.Poll(30*time.Second, 40*time.Second, func() (done bool, err error) {
				deploy, err := client.GetDeployment(cli, name, namespace)
				if err != nil {
					return false, err
				}
				if deploy.Status.ReadyReplicas != newReplicas {
					return false, nil
				}
				return true, nil
			})
		})
	})

	Context("execludeDate to scale up", func() {
		var initReplicas int32 = 2
		var targetSize int32 = 3

		BeforeEach(func() {
			err := PrepareScaleDeploy(cli, name, namespace, image, initReplicas)
			framework.ExpectNoError(err, "failed to patch deploy")
			_, err = client.CreateCronhpaExcludeToday(cronhpaCli, name, namespace, targetSize)
			framework.ExpectNoError(err, "faild to create cronhpa")
		})
		AfterEach(func() {
			err = CleanScaleDeploy(cli, name, namespace)
			framework.ExpectNoError(err, "failed to clean deploy")
			err := client.DeleteCronhpa(cronhpaCli, name, namespace)
			framework.ExpectNoError(err, "failed to delete cronhpa")
		})

		It("check replicas ok", func() {
			// wait deploy replicas unchanged
			err := wait.Poll(30*time.Second, 40*time.Second, func() (done bool, err error) {
				deploy, err := client.GetDeployment(cli, name, namespace)
				if err != nil {
					return false, err
				}
				if deploy.Status.ReadyReplicas != initReplicas {
					return false, nil
				}
				return true, nil
			})
			framework.ExpectNoError(err, "failed to scale up deploy")
		})
	})

	Context("runonce to scale up", func() {
		var initReplicas int32 = 2
		var targetSize int32 = 3

		BeforeEach(func() {
			err := PrepareScaleDeploy(cli, name, namespace, image, initReplicas)
			framework.ExpectNoError(err, "failed to patch deploy")
			_, err = client.CreateCronhpaRunonce(cronhpaCli, name, namespace, targetSize)
			framework.ExpectNoError(err, "faild to create cronhpa")
		})
		AfterEach(func() {
			err = CleanScaleDeploy(cli, name, namespace)
			framework.ExpectNoError(err, "failed to clean deploy")
			err := client.DeleteCronhpa(cronhpaCli, name, namespace)
			framework.ExpectNoError(err, "failed to delete cronhpa")
		})
		It("check replicas ok", func() {
			// check deploy scale
			err := wait.Poll(2*time.Second, 1*time.Minute, func() (done bool, err error) {
				deploy, err := client.GetDeployment(cli, name, namespace)
				if err != nil {
					return false, err
				}
				if deploy.Status.ReadyReplicas != targetSize {
					return false, nil
				}
				return true, nil
			})
			framework.ExpectNoError(err, "failed to scale deploy")
			// check runonce: 1. scale to 1
			var newReplicas int32 = 1
			_, err = client.PatchDeployReplicas(cli, name, namespace, newReplicas)
			framework.ExpectNoError(err, "failed to patch deploy")
			err = client.WaitDeployReplicas(cli, name, namespace, newReplicas)
			framework.ExpectNoError(err, "failed to wait deploy replicas")
			// check runonce: 2. deploy replicas unchanged
			err = wait.Poll(30*time.Second, 40*time.Second, func() (done bool, err error) {
				deploy, err := client.GetDeployment(cli, name, namespace)
				if err != nil {
					return false, err
				}
				if deploy.Status.ReadyReplicas != newReplicas {
					return false, nil
				}
				return true, nil
			})
			framework.ExpectNoError(err, "failed to check deploy unchanged")
		})
	})
})

func PrepareScaleDeploy(cli kubernetes.Interface, name, namespace, image string, replicas int32) error {
	_, err := client.CreateDeployment(cli, name, namespace, image, replicas)
	if err != nil {
		return err
	}
	err = client.WaitDeployReplicas(cli, name, namespace, replicas)
	return err
}

func CleanScaleDeploy(cli kubernetes.Interface, name, namespace string) error {
	err := client.DeleteDeployment(cli, name, namespace)
	return err
}
