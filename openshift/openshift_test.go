package simple

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/testernetes/gkube"

	openshiftapi "github.com/openshift/api"
	configv1 "github.com/openshift/api/config/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var _ = Describe("Openshift", Ordered, func() {

	var k8s KubernetesHelper

	BeforeAll(func() {
		scheme := runtime.NewScheme()
		openshiftapi.Install(scheme)

		k8s = NewKubernetesHelper(WithScheme(scheme))
	})

	It("Should have completed an Openshift Version upgrade", func(ctx SpecContext) {
		clusterversion := &configv1.ClusterVersion{
			ObjectMeta: metav1.ObjectMeta{
				Name: "version",
			},
		}
		Eventually(k8s.Object).WithContext(ctx).WithArguments(clusterversion).Should(
			HaveJSONPath(`{.status.conditions[?(@.type=="Available")].status}`,
				BeEquivalentTo(corev1.ConditionTrue)),
		)
	}, SpecTimeout(5*time.Second))
})

func TestKubernetesHelper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Openshift")
}
