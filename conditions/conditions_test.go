package conditions

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	. "github.com/testernetes/gkube"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Conditions", Ordered, func() {

	var deployment *appsv1.Deployment
	var k8s KubernetesHelper

	BeforeAll(func() {
		k8s = NewKubernetesHelper()
		appLabel := map[string]string{"app": "my-deployment"}
		deployment = &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-deployment",
				Namespace: "default",
			},
			Spec: appsv1.DeploymentSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: appLabel,
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: appLabel,
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Image:   "busybox:latest",
								Name:    "test",
								Command: []string{"/bin/sh"},
								Args:    []string{"-c", "sleep 60"},
							},
						},
					},
				},
			},
		}
	})

	It("Should create the deployment using the KubernetesHelper", func(ctx SpecContext) {
		Eventually(k8s.Create).WithContext(ctx).WithArguments(deployment).Should(Succeed())
	}, SpecTimeout(time.Minute))

	It("Should become healthy", func(ctx SpecContext) {
		Eventually(k8s.Object).WithContext(ctx).WithArguments(deployment).Should(
			HaveJSONPath(`{.status.conditions[?(@.type=="Available")]}`,
				MatchFields(IgnoreExtras, Fields{
					"Status": BeEquivalentTo(corev1.ConditionTrue),
					"Reason": Equal("MinimumReplicasAvailable"),
				}),
			))
	}, SpecTimeout(time.Minute))

	AfterAll(func(ctx SpecContext) {
		Eventually(k8s.Delete).WithContext(ctx).WithArguments(deployment).Should(Succeed())
	}, NodeTimeout(time.Minute))
})

func TestKubernetesHelper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Conditions")
}
