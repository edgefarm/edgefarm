package clusters

import (
	"fmt"
	"time"

	tmplutil "github.com/openyurtio/openyurt/pkg/util/templates"

	"github.com/edgefarm/edgefarm/pkg/k8s"
	"github.com/edgefarm/edgefarm/pkg/shared"
	"github.com/s0rg/retry"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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

func RenderAndApply(template string, context map[string]string, kubeconfig *rest.Config) error {
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
