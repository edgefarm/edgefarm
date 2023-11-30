package k8s

import (
	"context"
	"errors"
	"os"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclientset "k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

var (
	// PropagationPolicy defines the propagation policy used when deleting a resource
	PropagationPolicy = metav1.DeletePropagationBackground
	// CheckServantJobPeriod defines the time interval between two successive ServantJob statu's inspection
	CheckServantJobPeriod = time.Second * 10
)

// RunJobAndCleanup runs the job, wait for it to be complete, and delete it
func RunJobAndCleanup(cliSet kubeclientset.Interface, job *batchv1.Job, timeout, period time.Duration, waitForTimeout bool) error {
	job, err := cliSet.BatchV1().Jobs(job.GetNamespace()).Create(context.Background(), job, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	waitJobTimeout := time.After(timeout)
	defer func() {
		labelSelector, err := metav1.LabelSelectorAsSelector(job.Spec.Selector)
		if err != nil {
			return
		}
		podList, err := cliSet.CoreV1().Pods(job.GetNamespace()).List(context.TODO(), metav1.ListOptions{
			LabelSelector: labelSelector.String(),
		})
		if err != nil {
			return
		}

		if len(podList.Items) == 0 {
			return
		}
		if err := PrintPodLog(cliSet, &podList.Items[0], os.Stderr); err != nil {
			klog.Errorf("failed to print job pod logs, %v", err)
		}
	}()

	for {
		select {
		case <-waitJobTimeout:
			return errors.New("wait for job to be complete timeout")
		case <-time.After(period):
			newJob, err := cliSet.BatchV1().Jobs(job.GetNamespace()).
				Get(context.Background(), job.GetName(), metav1.GetOptions{})
			if err != nil {
				if apierrors.IsNotFound(err) {
					return err
				}

				if waitForTimeout {
					klog.Infof("continue to wait for job(%s) to complete until timeout, even if failed to get job, %v", job.GetName(), err)
					continue
				}
				return err
			}

			if newJob.Status.Succeeded == *newJob.Spec.Completions {
				if err := cliSet.BatchV1().Jobs(job.GetNamespace()).
					Delete(context.Background(), job.GetName(), metav1.DeleteOptions{
						PropagationPolicy: &PropagationPolicy,
					}); err != nil {
					klog.Errorf("fail to delete succeeded servant job(%s): %s", job.GetName(), err)
					return err
				}
				return nil
			}
		}
	}
}
