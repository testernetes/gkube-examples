package simple

import (
	"testing"
	"time"

	. "github.com/matt-simons/gkube"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Simple use of the KubernetesHelper", func() {

	It("Should create and delete a namespace", func() {

		By("Creating a new default KubernetesHelper")
		k8s := NewKubernetesHelper()

		namespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "simple",
			},
		}

		By("Using the KubernetesHelper to create the namespace")
		Eventually(k8s.Create(namespace)).Should(Succeed())

		By("Filtering a list of namespaces using a JSON path")
		Eventually(k8s.Objects(&corev1.NamespaceList{})).Should(WithJSONPath("{.items[*].metadata.name}", ContainElement("simple")))

		By("Waiting additional time for any namespaced objects to delete")
		Eventually(k8s.Delete(namespace)).WithTimeout(time.Minute).Should(Succeed())
	})
})

func TestKubernetesHelper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Simple")
}
