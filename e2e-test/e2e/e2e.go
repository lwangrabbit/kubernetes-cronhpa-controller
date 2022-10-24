// Copyright 2022 Cronhpa Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package e2e

import (
	"context" // load pprof
	_ "net/http/pprof"

	"github.com/onsi/ginkgo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/test/e2e/framework"
	e2elog "k8s.io/kubernetes/test/e2e/framework/log"
	utilnet "k8s.io/utils/net"

	e2econfig "github.com/AliyunContainerService/kubernetes-cronhpa-controller/e2e-test/e2e/config" // ensure auth plugins are loaded
)

var _ = ginkgo.SynchronizedBeforeSuite(func() []byte {
	if e2econfig.TestConfig.CheckCronhpa {
		c, err := framework.LoadClientset()
		if err != nil {
			klog.Fatal("Error loading client: ", err)
		}
		_, err = c.AppsV1().Deployments(metav1.NamespaceSystem).Get(context.TODO(), "kubernetes-cronhpa-controller", metav1.GetOptions{})
		if err != nil {
			e2elog.Failf("Failed to get Cronhpa deployment: %v", err)
		}
	}
	return nil
}, func(data []byte) {
	// Run on all Ginkgo nodes
	setupSuitePerGinkgoNode()
})

func setupSuitePerGinkgoNode() {
	c, err := framework.LoadClientset()
	if err != nil {
		klog.Fatal("Error loading client: ", err)
	}
	framework.TestContext.IPFamily = getDefaultClusterIPFamily(c)
	e2elog.Logf("Cluster IP family: %s", framework.TestContext.IPFamily)
}

// getDefaultClusterIPFamily obtains the default IP family of the cluster
// using the Cluster IP address of the kubernetes service created in the default namespace
// This unequivocally identifies the default IP family because services are single family
// TODO: dual-stack may support multiple families per service
// but we can detect if a cluster is dual stack because pods have two addresses (one per family)
func getDefaultClusterIPFamily(c kubernetes.Interface) string {
	// Get the ClusterIP of the kubernetes service created in the default namespace
	svc, err := c.CoreV1().Services(metav1.NamespaceDefault).Get(context.TODO(), "kubernetes", metav1.GetOptions{})
	if err != nil {
		e2elog.Failf("Failed to get kubernetes service ClusterIP: %v", err)
	}

	if utilnet.IsIPv6String(svc.Spec.ClusterIP) {
		return "ipv6"
	}
	return "ipv4"
}
