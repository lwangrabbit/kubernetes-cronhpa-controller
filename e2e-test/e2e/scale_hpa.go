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

type ScaleHpaCase struct {
	InputHpaMinReplicas int32
	InputHpaMaxReplicas int32

	TargetCronhpaReplicas int32
	CurrentDeployReplicas int32

	OutputHpaMinReplicas int32
	OutputHpaMaxReplicas int32
	OutputDeployReplicas int32
}

// https://help.aliyun.com/document_detail/151557.html

var _ = Describe("Scale Hpa", func() {
	name := "test-hpa"
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

	metricName := "metric_hpa"

	Context("case1", func() {
		tcase := ScaleHpaCase{
			InputHpaMinReplicas:   1,
			InputHpaMaxReplicas:   10,
			TargetCronhpaReplicas: 5,
			CurrentDeployReplicas: 5,
			OutputHpaMinReplicas:  1,
			OutputHpaMaxReplicas:  10,
			OutputDeployReplicas:  5,
		}
		BeforeEach(func() {
			// create deploy
			_, err := client.CreateDeployment(cli, name, namespace, image, tcase.CurrentDeployReplicas)
			framework.ExpectNoError(err, "failed to create deploy")
			// wait deploy replicas ready
			err = client.WaitDeployReplicas(cli, name, namespace, tcase.CurrentDeployReplicas)
			framework.ExpectNoError(err, "failed to wait deploy replicas")
			// create hpa
			_, err = client.CreateHpa(cli, name, namespace, metricName, tcase.InputHpaMinReplicas, tcase.InputHpaMaxReplicas)
			framework.ExpectNoError(err, "failed to create hpa")
			// patch hpa.status.currentReplicas
			_, err = client.UpdateHpaStatus(cli, name, namespace, tcase.CurrentDeployReplicas)
			framework.ExpectNoError(err, "failed to patch hpa status")
			// create cronhpa
			_, err = client.CreateCronhpaHpa(cronhpaCli, name, namespace, tcase.TargetCronhpaReplicas)
			framework.ExpectNoError(err, "failed to create cronhpa")
		})
		AfterEach(func() {
			// delete cronhpa
			err := client.DeleteCronhpa(cronhpaCli, name, namespace)
			framework.ExpectNoError(err, "failed to delete cronhpa")
			// delete hpa
			err = client.DeleteHpa(cli, name, namespace)
			framework.ExpectNoError(err, "failed to delete hpa")
			// delete deploy
			err = client.DeleteDeployment(cli, name, namespace)
			framework.ExpectNoError(err, "failed to delete deploy")
		})
		It("check cronhpa behavior", func() {
			// Wait hpa replicas
			err := wait.Poll(2*time.Second, 1*time.Minute, func() (done bool, err error) {
				hpa, err := client.GetHpa(cli, name, namespace)
				if err != nil {
					return false, err
				}
				if *hpa.Spec.MinReplicas != tcase.OutputHpaMinReplicas {
					return false, nil
				}
				if hpa.Spec.MaxReplicas != tcase.OutputHpaMaxReplicas {
					return false, nil
				}
				return true, nil
			})
			framework.ExpectNoError(err, "failed to get expect hpa replicas")
			// wait deploy replicas
			err = wait.Poll(2*time.Second, 1*time.Minute, func() (done bool, err error) {
				deploy, err := client.GetDeployment(cli, name, namespace)
				if err != nil {
					return false, err
				}
				if *deploy.Spec.Replicas != tcase.OutputDeployReplicas {
					return false, nil
				}
				return true, nil
			})
			framework.ExpectNoError(err, "failed to get expect deploy replicas")
		})
	})

	Context("case2", func() {
		tcase := ScaleHpaCase{
			InputHpaMinReplicas:   1,
			InputHpaMaxReplicas:   10,
			TargetCronhpaReplicas: 4,
			CurrentDeployReplicas: 5,
			OutputHpaMinReplicas:  1,
			OutputHpaMaxReplicas:  10,
			OutputDeployReplicas:  5,
		}
		BeforeEach(func() {
			// create deploy
			_, err := client.CreateDeployment(cli, name, namespace, image, tcase.CurrentDeployReplicas)
			framework.ExpectNoError(err, "failed to create deploy")
			// wait deploy replicas ready
			err = client.WaitDeployReplicas(cli, name, namespace, tcase.CurrentDeployReplicas)
			framework.ExpectNoError(err, "failed to wait deploy replicas")
			// create hpa
			_, err = client.CreateHpa(cli, name, namespace, metricName, tcase.InputHpaMinReplicas, tcase.InputHpaMaxReplicas)
			framework.ExpectNoError(err, "failed to create hpa")
			// patch hpa.status.currentReplicas
			_, err = client.UpdateHpaStatus(cli, name, namespace, tcase.CurrentDeployReplicas)
			framework.ExpectNoError(err, "failed to patch hpa status")
			// create cronhpa
			_, err = client.CreateCronhpaHpa(cronhpaCli, name, namespace, tcase.TargetCronhpaReplicas)
			framework.ExpectNoError(err, "failed to create cronhpa")
		})
		AfterEach(func() {
			// delete cronhpa
			err := client.DeleteCronhpa(cronhpaCli, name, namespace)
			framework.ExpectNoError(err, "failed to delete cronhpa")
			// delete hpa
			err = client.DeleteHpa(cli, name, namespace)
			framework.ExpectNoError(err, "failed to delete hpa")
			// delete deploy
			err = client.DeleteDeployment(cli, name, namespace)
			framework.ExpectNoError(err, "failed to delete deploy")
		})
		It("check cronhpa behavior", func() {
			// Wait hpa replicas
			err := wait.Poll(2*time.Second, 1*time.Minute, func() (done bool, err error) {
				hpa, err := client.GetHpa(cli, name, namespace)
				if err != nil {
					return false, err
				}
				if *hpa.Spec.MinReplicas != tcase.OutputHpaMinReplicas {
					return false, nil
				}
				if hpa.Spec.MaxReplicas != tcase.OutputHpaMaxReplicas {
					return false, nil
				}
				return true, nil
			})
			framework.ExpectNoError(err, "failed to get expect hpa replicas")
			// wait deploy replicas
			err = wait.Poll(2*time.Second, 1*time.Minute, func() (done bool, err error) {
				deploy, err := client.GetDeployment(cli, name, namespace)
				if err != nil {
					return false, err
				}
				if *deploy.Spec.Replicas != tcase.OutputDeployReplicas {
					return false, nil
				}
				return true, nil
			})
			framework.ExpectNoError(err, "failed to get expect deploy replicas")
		})
	})

	Context("case3", func() {
		tcase := ScaleHpaCase{
			InputHpaMinReplicas:   1,
			InputHpaMaxReplicas:   10,
			TargetCronhpaReplicas: 6,
			CurrentDeployReplicas: 5,
			OutputHpaMinReplicas:  6,
			OutputHpaMaxReplicas:  10,
			OutputDeployReplicas:  6,
		}
		BeforeEach(func() {
			// create deploy
			_, err := client.CreateDeployment(cli, name, namespace, image, tcase.CurrentDeployReplicas)
			framework.ExpectNoError(err, "failed to create deploy")
			// wait deploy replicas ready
			err = client.WaitDeployReplicas(cli, name, namespace, tcase.CurrentDeployReplicas)
			framework.ExpectNoError(err, "failed to wait deploy replicas")
			// create hpa
			_, err = client.CreateHpa(cli, name, namespace, metricName, tcase.InputHpaMinReplicas, tcase.InputHpaMaxReplicas)
			framework.ExpectNoError(err, "failed to create hpa")
			// patch hpa.status.currentReplicas
			_, err = client.UpdateHpaStatus(cli, name, namespace, tcase.CurrentDeployReplicas)
			framework.ExpectNoError(err, "failed to patch hpa status")
			// create cronhpa
			_, err = client.CreateCronhpaHpa(cronhpaCli, name, namespace, tcase.TargetCronhpaReplicas)
			framework.ExpectNoError(err, "failed to create cronhpa")
		})
		AfterEach(func() {
			// delete cronhpa
			err := client.DeleteCronhpa(cronhpaCli, name, namespace)
			framework.ExpectNoError(err, "failed to delete cronhpa")
			// delete hpa
			err = client.DeleteHpa(cli, name, namespace)
			framework.ExpectNoError(err, "failed to delete hpa")
			// delete deploy
			err = client.DeleteDeployment(cli, name, namespace)
			framework.ExpectNoError(err, "failed to delete deploy")
		})
		It("check cronhpa behavior", func() {
			// Wait hpa replicas
			err := wait.Poll(2*time.Second, 1*time.Minute, func() (done bool, err error) {
				hpa, err := client.GetHpa(cli, name, namespace)
				if err != nil {
					return false, err
				}
				if *hpa.Spec.MinReplicas != tcase.OutputHpaMinReplicas {
					return false, nil
				}
				if hpa.Spec.MaxReplicas != tcase.OutputHpaMaxReplicas {
					return false, nil
				}
				return true, nil
			})
			framework.ExpectNoError(err, "failed to get expect hpa replicas")
			// wait deploy replicas
			err = wait.Poll(2*time.Second, 1*time.Minute, func() (done bool, err error) {
				deploy, err := client.GetDeployment(cli, name, namespace)
				if err != nil {
					return false, err
				}
				if *deploy.Spec.Replicas != tcase.OutputDeployReplicas {
					return false, nil
				}
				return true, nil
			})
			framework.ExpectNoError(err, "failed to get expect deploy replicas")
		})
	})

	Context("case4", func() {
		tcase := ScaleHpaCase{
			InputHpaMinReplicas:   5,
			InputHpaMaxReplicas:   10,
			TargetCronhpaReplicas: 4,
			CurrentDeployReplicas: 5,
			OutputHpaMinReplicas:  4,
			OutputHpaMaxReplicas:  10,
			OutputDeployReplicas:  5,
		}
		BeforeEach(func() {
			// create deploy
			_, err := client.CreateDeployment(cli, name, namespace, image, tcase.CurrentDeployReplicas)
			framework.ExpectNoError(err, "failed to create deploy")
			// wait deploy replicas ready
			err = client.WaitDeployReplicas(cli, name, namespace, tcase.CurrentDeployReplicas)
			framework.ExpectNoError(err, "failed to wait deploy replicas")
			// create hpa
			hpa, err := client.CreateHpa(cli, name, namespace, metricName, tcase.InputHpaMinReplicas, tcase.InputHpaMaxReplicas)
			framework.ExpectNoError(err, "failed to create hpa")
			// patch hpa.status.currentReplicas
			hpa.Status.CurrentReplicas = tcase.CurrentDeployReplicas
			_, err = client.UpdateHpaStatus(cli, name, namespace, tcase.CurrentDeployReplicas)
			framework.ExpectNoError(err, "failed to patch hpa status")
			// create cronhpa
			_, err = client.CreateCronhpaHpa(cronhpaCli, name, namespace, tcase.TargetCronhpaReplicas)
			framework.ExpectNoError(err, "failed to create cronhpa")
		})
		AfterEach(func() {
			// delete cronhpa
			err := client.DeleteCronhpa(cronhpaCli, name, namespace)
			framework.ExpectNoError(err, "failed to delete cronhpa")
			// delete hpa
			err = client.DeleteHpa(cli, name, namespace)
			framework.ExpectNoError(err, "failed to delete hpa")
			// delete deploy
			err = client.DeleteDeployment(cli, name, namespace)
			framework.ExpectNoError(err, "failed to delete deploy")
		})
		It("check cronhpa behavior", func() {
			// Wait hpa replicas
			err := wait.Poll(2*time.Second, 1*time.Minute, func() (done bool, err error) {
				hpa, err := client.GetHpa(cli, name, namespace)
				if err != nil {
					return false, err
				}
				if *hpa.Spec.MinReplicas != tcase.OutputHpaMinReplicas {
					return false, nil
				}
				if hpa.Spec.MaxReplicas != tcase.OutputHpaMaxReplicas {
					return false, nil
				}
				return true, nil
			})
			framework.ExpectNoError(err, "failed to get expect hpa replicas")
			// wait deploy replicas
			err = wait.Poll(2*time.Second, 1*time.Minute, func() (done bool, err error) {
				deploy, err := client.GetDeployment(cli, name, namespace)
				if err != nil {
					return false, err
				}
				if *deploy.Spec.Replicas != tcase.OutputDeployReplicas {
					return false, nil
				}
				return true, nil
			})
			framework.ExpectNoError(err, "failed to get expect deploy replicas")
		})
	})

	Context("case5", func() {
		tcase := ScaleHpaCase{
			InputHpaMinReplicas:   5,
			InputHpaMaxReplicas:   10,
			TargetCronhpaReplicas: 11,
			CurrentDeployReplicas: 5,
			OutputHpaMinReplicas:  11,
			OutputHpaMaxReplicas:  11,
			OutputDeployReplicas:  11,
		}
		BeforeEach(func() {
			// create deploy
			_, err := client.CreateDeployment(cli, name, namespace, image, tcase.CurrentDeployReplicas)
			framework.ExpectNoError(err, "failed to create deploy")
			// wait deploy replicas ready
			err = client.WaitDeployReplicas(cli, name, namespace, tcase.CurrentDeployReplicas)
			framework.ExpectNoError(err, "failed to wait deploy replicas")
			// create hpa
			_, err = client.CreateHpa(cli, name, namespace, metricName, tcase.InputHpaMinReplicas, tcase.InputHpaMaxReplicas)
			framework.ExpectNoError(err, "failed to create hpa")
			// patch hpa.status.currentReplicas
			_, err = client.UpdateHpaStatus(cli, name, namespace, tcase.CurrentDeployReplicas)
			framework.ExpectNoError(err, "failed to patch hpa status")
			// create cronhpa
			_, err = client.CreateCronhpaHpa(cronhpaCli, name, namespace, tcase.TargetCronhpaReplicas)
			framework.ExpectNoError(err, "failed to create cronhpa")
		})
		AfterEach(func() {
			// delete cronhpa
			err := client.DeleteCronhpa(cronhpaCli, name, namespace)
			framework.ExpectNoError(err, "failed to delete cronhpa")
			// delete hpa
			err = client.DeleteHpa(cli, name, namespace)
			framework.ExpectNoError(err, "failed to delete hpa")
			// delete deploy
			err = client.DeleteDeployment(cli, name, namespace)
			framework.ExpectNoError(err, "failed to delete deploy")
		})
		It("check cronhpa behavior", func() {
			// Wait hpa replicas
			err := wait.Poll(2*time.Second, 1*time.Minute, func() (done bool, err error) {
				hpa, err := client.GetHpa(cli, name, namespace)
				if err != nil {
					return false, err
				}
				if *hpa.Spec.MinReplicas != tcase.OutputHpaMinReplicas {
					return false, nil
				}
				if hpa.Spec.MaxReplicas != tcase.OutputHpaMaxReplicas {
					return false, nil
				}
				return true, nil
			})
			framework.ExpectNoError(err, "failed to get expect hpa replicas")
			// wait deploy replicas
			err = wait.Poll(2*time.Second, 1*time.Minute, func() (done bool, err error) {
				deploy, err := client.GetDeployment(cli, name, namespace)
				if err != nil {
					return false, err
				}
				if *deploy.Spec.Replicas != tcase.OutputDeployReplicas {
					return false, nil
				}
				return true, nil
			})
			framework.ExpectNoError(err, "failed to get expect deploy replicas")
		})
	})

})
