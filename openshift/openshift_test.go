package simple

import (
	"testing"

	. "github.com/matt-simons/gkube"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	openshiftapi "github.com/openshift/api"
	configv1 "github.com/openshift/api/config/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var _ = Describe("customized use of the KubernetesHelper", func() {

	It("Should have completed an Openshift Version upgrade", func() {

		By("Creating a new KubernetesHelper with custom schema")

		scheme := runtime.NewScheme()
		openshiftapi.Install(scheme)

		k8s := NewKubernetesHelper(WithScheme(scheme))

		clusterversion := &configv1.ClusterVersion{
			ObjectMeta: metav1.ObjectMeta{
				Name: "version",
			},
		}

		By("Filtering the avaiable condition using a JSON path")
		Eventually(k8s.Object(clusterversion)).Should(WithJSONPath(
			`{.status.conditions[?(@.type=="Available")].status}`,
			BeEquivalentTo(corev1.ConditionTrue)),
		)
	})
})

func TestKubernetesHelper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Simple")
}
