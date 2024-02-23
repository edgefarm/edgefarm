package clusters

import (
	"fmt"
	"time"

	tmplutil "github.com/openyurtio/openyurt/pkg/util/templates"

	"github.com/edgefarm/edgefarm/pkg/k8s"
	"github.com/edgefarm/edgefarm/pkg/shared"
	"github.com/s0rg/retry"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

func PrepareKubeClient(kubeconfigPath string) (*rest.Config, error) {
	p, err := shared.Expand(kubeconfigPath)
	if err != nil {
		return nil, err
	}
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", p)
	if err != nil {
		return nil, err
	}
	return kubeconfig, nil
}

func WaitForKubeconfig(kubeconfig *rest.Config, cluster string, timeout time.Duration) (bool, error) {

	try := retry.New(
		retry.Count(int(timeout.Seconds())),
		retry.Sleep(time.Second),
		retry.Verbose(true),
	)
	if err := try.Single(fmt.Sprintf("Waiting for Kubeconfig for cluster %s", cluster),
		func() error {
			s, err := k8s.GetSecret(kubeconfig, cluster, "default")
			if err != nil {
				return err
			}
			if s != nil {
				_, err := k8s.SecretValue(s, "value")
				if err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
		return false, err
	}
	return true, nil
}

func RenderAndApply(name string, template string, context map[string]string, kubeconfig *rest.Config) error {
	klog.Infoln(fmt.Sprintf("Creating %s", name))
	manifest, err := prepareTemplate(template, context)
	if err != nil {
		return err
	}
	if err := k8s.Apply(kubeconfig, manifest); err != nil {
		return err
	}
	return nil
}

func prepareTemplate(tempalte string, context map[string]string) (string, error) {
	result, err := tmplutil.SubsituteTemplate(tempalte, context)
	if err != nil {
		return "", err
	}
	return result, nil
}

func WaitForAPIServerReachable(kubeconfig *rest.Config, timeout time.Duration) (bool, error) {
	try := retry.New(
		retry.Count(int(timeout.Seconds())),
		retry.Sleep(time.Second),
		retry.Verbose(true),
	)
	if err := try.Single("Waiting for API Server to be reachable",
		func() error {
			_, err := k8s.GetClientset(kubeconfig)
			if err != nil {
				return err
			}
			control, err := k8s.GetNodes(kubeconfig, &v1.LabelSelector{
				MatchLabels: map[string]string{
					"node-role.kubernetes.io/control-plane": "",
				},
			})
			if len(control) >= 1 {
				return nil
			}
			if err != nil {
				return err
			}

			return nil
		}); err != nil {
		return false, err
	}
	return true, nil
}

func WaitForAllNodes(kubeconfig *rest.Config, timeout time.Duration, desiredControlPlaneNodes, desiredWorkerNodes int) (bool, error) {
	try := retry.New(
		retry.Count(int(timeout.Seconds())),
		retry.Sleep(time.Second),
		retry.Verbose(true),
	)
	if err := try.Single("Waiting for all nodes to be ready",
		func() error {
			control, err := k8s.GetNodes(kubeconfig, &v1.LabelSelector{
				MatchLabels: map[string]string{
					"node-role.kubernetes.io/control-plane": "",
				},
			})
			if err != nil {
				return err
			}

			allNodes, err := k8s.GetNodes(kubeconfig, &v1.LabelSelector{
				MatchLabels: map[string]string{
					"kubernetes.io/os": "linux",
				},
			})
			if err != nil {
				return err
			}

			currentControlPlaneNodes := len(control)
			currentWorkerNodes := len(allNodes) - currentControlPlaneNodes
			if currentWorkerNodes == desiredWorkerNodes && currentControlPlaneNodes == desiredControlPlaneNodes {
				return nil
			}
			return func() error {
				errorStr := "waiting for "
				if currentControlPlaneNodes < desiredControlPlaneNodes {
					errorStr += fmt.Sprintf("%d additional control plane nodes ", desiredControlPlaneNodes-currentControlPlaneNodes)
				}
				if currentWorkerNodes < desiredWorkerNodes {
					errorStr += fmt.Sprintf("%d additional worker nodes", desiredWorkerNodes-currentWorkerNodes)
				}
				return fmt.Errorf(errorStr)
			}()
		}); err != nil {
		return false, err
	}
	return true, nil
}
