package conditions

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	. "github.com/tkube/gkube"

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

	It("Should create the deployment using the KubernetesHelper", func() {
		Eventually(k8s.Create(deployment)).Should(Succeed())
	})

	It("Should become healthy", func() {
		Eventually(k8s.Object(deployment)).WithTimeout(time.Minute).Should(WithJSONPath(
			`{.status.conditions[?(@.type=="Available")]}`,
			MatchFields(IgnoreExtras, Fields{
				"Status": BeEquivalentTo(corev1.ConditionTrue),
				"Reason": Equal("MinimumReplicasAvailable"),
			}),
		))
	})

	AfterAll(func() {
		Eventually(k8s.Delete(deployment)).Should(Succeed())
	})
})

func TestKubernetesHelper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Conditions")
}
