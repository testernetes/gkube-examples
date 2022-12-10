package trackedclient

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/testernetes/gkube"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	tc "github.com/testernetes/trackedclient"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

var _ = Describe("TrackedClient", Ordered, func() {

	var c tc.TrackedClient
	var k8s KubernetesHelper

	BeforeAll(func() {
		var err error
		c, err = tc.New(config.GetConfigOrDie(), client.Options{Scheme: scheme.Scheme})
		Expect(err).ShouldNot(HaveOccurred())

		k8s = NewKubernetesHelper(WithClient(c))
	})

	It("should create a namespace", func(ctx SpecContext) {
		namespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "my-namespace",
			},
		}
		Eventually(k8s.Create(ctx, namespace)).Should(Succeed())
		Eventually(k8s.Object(ctx, namespace)).Should(
			HaveJSONPath("{.status.phase}", Equal(corev1.NamespacePhase("Active"))),
		)
	}, SpecTimeout(time.Minute))

	AfterAll(func(ctx SpecContext) {
		Expect(c.DeleteAllTracked(ctx)).Should(Succeed())
	}, NodeTimeout(time.Minute))
})

func TestKubernetesHelper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TrackedClient")
}
