package podexec

import (
	"context"
	"exec"
	"os/exec"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	. "github.com/testernetes/gkube"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Conditions", Ordered, func() {

	var pod *corev1.Pod
	var k8s KubernetesHelper

	BeforeAll(func() {
		k8s = NewKubernetesHelper()
		pod = &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "busy",
				Namespace: "default",
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
		}
		Eventually(k8s.Create(pod)).Should(Succeed())
		Eventually(k8s.Object(pod)).WithTimeout(time.Minute).Should(WithJSONPath(
			`{.status.phase}`, BeEquivalentTo(corev1.PodRunning)))
	})

	It("be able to resolve google", func() {
		cmd := exec.CommandContext(context.TODO(), "nslookup", "google.com")
		k8s.Exec(pod, "test", cmd, GinkgoWriter, GinkgoWriter)
	})

	AfterAll(func() {
		Eventually(k8s.Delete(pod)).Should(Succeed())
	})
})

func TestKubernetesHelper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Conditions")
}
