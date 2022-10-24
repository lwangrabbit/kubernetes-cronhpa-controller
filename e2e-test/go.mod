module github.com/AliyunContainerService/kubernetes-cronhpa-controller/e2e-test

go 1.16

require (
	github.com/AliyunContainerService/kubernetes-cronhpa-controller v0.0.0-00010101000000-000000000000
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.21.1
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.44.1
	github.com/prometheus-operator/prometheus-operator/pkg/client v0.47.0
	k8s.io/api v0.20.12
	k8s.io/apimachinery v0.20.12
	k8s.io/cli-runtime v0.20.12
	k8s.io/client-go v0.20.12
	k8s.io/component-base v0.20.12
	k8s.io/klog/v2 v2.4.0
	k8s.io/kubernetes v1.20.12
	k8s.io/utils v0.0.0-20201110183641-67b214c5f920
)

replace (
	github.com/AliyunContainerService/kubernetes-cronhpa-controller => ../
	k8s.io/api => k8s.io/api v0.20.12
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.20.12
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.12
	k8s.io/apiserver => k8s.io/apiserver v0.20.12
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.20.12
	k8s.io/client-go => k8s.io/client-go v0.20.12
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.20.12
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.20.12
	k8s.io/code-generator => k8s.io/code-generator v0.20.12
	k8s.io/component-base => k8s.io/component-base v0.20.12
	k8s.io/component-helpers => k8s.io/component-helpers v0.20.12
	k8s.io/controller-manager => k8s.io/controller-manager v0.20.12
	k8s.io/cri-api => k8s.io/cri-api v0.25.0-alpha.0
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.20.12
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.20.12
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.20.12
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.20.12
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.20.12
	k8s.io/kubectl => k8s.io/kubectl v0.20.12
	k8s.io/kubelet => k8s.io/kubelet v0.20.12
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.20.12
	k8s.io/metrics => k8s.io/metrics v0.20.12
	k8s.io/mount-utils => k8s.io/mount-utils v0.20.12
	k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.20.12
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.20.12
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.20.12
	k8s.io/sample-controller => k8s.io/sample-controller v0.20.12
)
