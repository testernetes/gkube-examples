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
	var cm *corev1.ConfigMap

	BeforeAll(func() {
		k8s = NewKubernetesHelper()
		cm = &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "simple-example",
				Namespace: "default",
			},
		}
	})

	It("should create a configmap", func() {
		Eventually(k8s.Create(cm)).Should(Succeed())
	})

	It("should update the configmap", func() {
		Eventually(k8s.Update(cm, func() error {
			cm.Data = map[string]string{
				"something": "simple",
			}
			return nil
		})).Should(Succeed())
	})

	It("should contain something simple ", func() {
		Eventually(k8s.Object(cm)).Should(
			HaveJSONPath("{.data.something}", Equal("simple")),
		)
	})

	AfterAll(func() {
		Eventually(k8s.Delete(cm)).Should(Succeed())
	})
})

func TestKubernetesHelper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Simple")
}
