package k8s

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PatchCoreDNSDeployment patches the CoreDNS deployment to contain the edgefarm.io/NoSchedule toleration
func PatchCoreDNSDeployment() error {
	clientset, err := GetClientset(nil)
	if err != nil {
		return err
	}

	deployment, err := clientset.AppsV1().Deployments("kube-system").Get(context.Background(), "coredns", metav1.GetOptions{})
	if err != nil {
		return err
	}

	tolerations := deployment.Spec.Template.Spec.Tolerations
	if tolerations == nil {
		tolerations = []v1.Toleration{}
	}

	tolerations = append(tolerations, v1.Toleration{
		Key:    "edgefarm.io",
		Effect: v1.TaintEffectNoSchedule,
	})

	deployment.Spec.Template.Spec.Tolerations = tolerations
	_, err = clientset.AppsV1().Deployments("kube-system").Update(context.Background(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}
