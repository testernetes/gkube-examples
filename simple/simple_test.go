package simple

import (
	"testing"
	"time"

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

	It("should create a configmap", func(ctx SpecContext) {
		Eventually(k8s.Create(ctx, cm)).Should(Succeed())
	}, SpecTimeout(time.Minute))

	It("should update the configmap", func(ctx SpecContext) {
		Eventually(k8s.Update(ctx, cm, func() error {
			cm.Data = map[string]string{
				"something": "simple",
			}
			return nil
		})).Should(Succeed())
	}, SpecTimeout(time.Minute))

	It("should contain something simple ", func(ctx SpecContext) {
		Eventually(k8s.Object(ctx, cm)).Should(
			HaveJSONPath("{.data.something}", Equal("simple")),
		)
	}, SpecTimeout(time.Minute))

	AfterAll(func(ctx SpecContext) {
		Eventually(k8s.Delete(ctx, cm)).Should(Succeed())
	}, NodeTimeout(time.Minute))
})

func TestKubernetesHelper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Simple")
}
