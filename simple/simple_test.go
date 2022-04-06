package simple

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/testernetes/gkube"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Simple use of the KubernetesHelper", Ordered, func() {

	var k8s KubernetesHelper
	var namespace *corev1.Namespace

	BeforeAll(func() {
		k8s = NewKubernetesHelper()
		namespace = &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "simple",
			},
		}
	})

	It("should create a namespace", func() {
		Eventually(k8s.Create(namespace)).Should(Succeed())
		Eventually(k8s.Object(namespace)).Should(
			WithJSONPath("{.status.phase}", Equal(corev1.NamespacePhase("Active"))),
		)
	})

	It("should filter a list of namespaces using a JSONPath", func() {
		Eventually(k8s.Objects(&corev1.NamespaceList{})).Should(
			WithJSONPath("{.items[*].metadata.name}", ContainElement("simple")),
		)
	})

	AfterAll(func() {
		Eventually(k8s.Delete(namespace, GracePeriodSeconds(30))).Should(Succeed())
	})
})

func TestKubernetesHelper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Simple")
}
