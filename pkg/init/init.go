/*
Copyright 2022 The OpenYurt Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package init

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	kubeclientset "k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"

	mycontext "github.com/edgefarm/edgefarm/pkg/context"
	"github.com/edgefarm/edgefarm/pkg/k8s"
	"github.com/edgefarm/edgefarm/pkg/k8s/tokens"
	"github.com/edgefarm/edgefarm/pkg/packages"

	"github.com/openyurtio/openyurt/pkg/projectinfo"
	strutil "github.com/openyurtio/openyurt/pkg/util/strings"
	"github.com/openyurtio/openyurt/pkg/yurthub/util"
	"github.com/openyurtio/openyurt/test/e2e/cmd/init/lock"
	kubeutil "github.com/openyurtio/openyurt/test/e2e/cmd/init/util/kubernetes"
)

const (
	// defaultYurthubHealthCheckTimeout defines the default timeout for yurthub health check phase
	defaultYurthubHealthCheckTimeout = 2 * time.Minute

	yssYurtHubCloudName = "yurt-static-set-yurt-hub-cloud"
	yssYurtHubName      = "yurt-static-set-yurt-hub"
)

type ClusterConverter struct {
	RootDir                   string
	ComponentsBuilder         *kubeutil.Builder
	ClientSet                 kubeclientset.Interface
	CloudNodes                []string
	EdgeNodes                 []string
	WaitServantJobTimeout     time.Duration
	YurthubHealthCheckTimeout time.Duration
	KubeConfigPath            string
	YurtManagerImage          string
	NodeServantImage          string
	YurthubImage              string
	EnableDummyIf             bool
}

func (c *ClusterConverter) Run() error {
	if err := lock.AcquireLock(c.ClientSet); err != nil {
		return err
	}
	defer func() {
		if releaseLockErr := lock.ReleaseLock(c.ClientSet); releaseLockErr != nil {
			klog.Error(releaseLockErr)
		}
	}()

	klog.Info("Deploying yurt-manager")
	if err := c.deployYurtManager(); err != nil {
		klog.Errorf("failed to deploy yurt-manager with image %s, %s", c.YurtManagerImage, err)
		return err
	}

	if err := packages.Install(packages.ClusterBootstrapYurtHub); err != nil {
		klog.Errorf("error occurs when deploying Yurthub, %v", err)
		return err
	}

	if err := c.prepareyNodeServantApplier(); err != nil {
		klog.Errorf("error occurs when preparing node servant applier, %v", err)
		return err
	}
	return nil
}
func AddWorkerLabelAndAutonomyAnnotation(cliSet kubeclientset.Interface, node *corev1.Node, lVal, aVal string) (*corev1.Node, error) {
	node.Labels["node.edgefarm.io/type"] = "virtual"
	node.Labels[projectinfo.GetEdgeWorkerLabelKey()] = lVal
	node.Annotations[projectinfo.GetAutonomyAnnotation()] = aVal
	newNode, err := cliSet.CoreV1().Nodes().Update(context.Background(), node, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return newNode, nil
}

func LabelEdgeNodes(edgeNodes []string) error {
	clientset, err := k8s.GetClientset(nil)
	if err != nil {
		return err
	}
	nodeLst, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list nodes, %w", err)
	}
	for _, node := range nodeLst.Items {
		isEdge := strutil.IsInStringLst(edgeNodes, node.Name)
		if !isEdge {
			continue
		}
		_, err := AddWorkerLabelAndAutonomyAnnotation(clientset, &node, "true", "true")
		if err != nil {
			return fmt.Errorf("failed to add label to edge node %s, %w", node.Name, err)
		}
	}
	return nil
}

func LabelCloudNodes(cloudNodes []string) error {
	reducedCloudNodes := []string{}
	// remove *-control-plane nodes from cloud nodes
	for _, name := range cloudNodes {
		if strings.Contains(name, "control-plane") {
			continue
		}
		reducedCloudNodes = append(reducedCloudNodes, name)
	}

	clientset, err := k8s.GetClientset(nil)
	if err != nil {
		return err
	}
	nodeLst, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list nodes, %w", err)
	}
	for _, node := range nodeLst.Items {
		isCloud := strutil.IsInStringLst(reducedCloudNodes, node.Name)
		if !isCloud {
			continue
		}
		_, err := AddWorkerLabelAndAutonomyAnnotation(clientset, &node, "false", "false")
		if err != nil {
			return fmt.Errorf("failed to add label to edge node %s, %w", node.Name, err)
		}
	}
	return nil
}

// func (c *ClusterConverter) labelEdgeNodes() error {
// 	nodeLst, err := c.ClientSet.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
// 	if err != nil {
// 		return fmt.Errorf("failed to list nodes, %w", err)
// 	}
// 	for _, node := range nodeLst.Items {
// 		isEdge := strutil.IsInStringLst(c.EdgeNodes, node.Name)
// 		if _, err = kubeutil.AddEdgeWorkerLabelAndAutonomyAnnotation(
// 			c.ClientSet, &node, strconv.FormatBool(isEdge), "false"); err != nil {
// 			return fmt.Errorf("failed to add label to edge node %s, %w", node.Name, err)
// 		}
// 	}
// 	return nil
// }

// func (c *ClusterConverter) deployYurtHub() error {
// 	return packages.Install(packages.ClusterBootstrapYurtHub)
// }

// // waiting yurt-manager pod ready
// return wait.PollImmediate(10*time.Second, 5*time.Minute, func() (bool, error) {
// 	podList, err := c.ClientSet.CoreV1().Pods("kube-system").List(context.TODO(), metav1.ListOptions{
// 		LabelSelector: labels.SelectorFromSet(map[string]string{"app.kubernetes.io/name": "yurt-hub"}).String(),
// 	})
// 	if err != nil {
// 		klog.Errorf("failed to list yurt-hub pod, %v", err)
// 		return false, nil
// 	} else if len(podList.Items) == 0 {
// 		klog.Infof("no yurt-manager pod: %#v", podList)
// 		return false, nil
// 	}

// 	if podList.Items[0].Status.Phase == corev1.PodRunning {
// 		for i := range podList.Items[0].Status.Conditions {
// 			if podList.Items[0].Status.Conditions[i].Type == corev1.PodReady &&
// 				podList.Items[0].Status.Conditions[i].Status == corev1.ConditionTrue {
// 				return true, nil
// 			}
// 		}
// 	}
// 	klog.Infof("pod(%s/%s): %#v", podList.Items[0].Namespace, podList.Items[0].Name, podList.Items[0])
// 	return false, nil
// })
// }

func (c *ClusterConverter) prepareyNodeServantApplier() error {
	// if err := prepareClusterInfoConfigMap(c.ClientSet, c.KubeConfigPath); err != nil {
	// 	return err
	// }

	joinToken, err := tokens.GetOrCreateJoinTokenString(c.ClientSet)
	if err != nil || joinToken == "" {
		return fmt.Errorf("fail to get join token: %w", err)
	}

	convertCtx := map[string]interface{}{
		"node_servant_image": c.NodeServantImage,
		"yurthub_image":      c.YurthubImage,
		"joinToken":          joinToken,
		// The node-servant will detect the kubeadm_conf_path automatically
		// It will be either "/usr/lib/systemd/system/kubelet.service.d/10-kubeadm.conf"
		// or "/etc/systemd/system/kubelet.service.d/10-kubeadm.conf".
		"kubeadm_conf_path": "",
		"working_mode":      string(util.WorkingModeEdge),
		"enable_dummy_if":   strconv.FormatBool(c.EnableDummyIf),
	}
	if c.YurthubHealthCheckTimeout != defaultYurthubHealthCheckTimeout {
		convertCtx["yurthub_healthcheck_timeout"] = c.YurthubHealthCheckTimeout.String()
	}

	npExist, err := nodePoolResourceExists(c.ClientSet)
	if err != nil {
		return err
	}
	convertCtx["enable_node_pool"] = strconv.FormatBool(npExist)
	convertCtx["configmap_name"] = yssYurtHubName
	klog.Infof("convert context for edge nodes(%q): %#+v", c.EdgeNodes, convertCtx)

	// Create context for node servant applier helm chart installed in clusterCreate
	mycontext.Context("node-servant-applier", mycontext.WithData(convertCtx))

	// if len(c.EdgeNodes) != 0 {
	// 	convertCtx["working_mode"] = string(util.WorkingModeEdge)
	// 	job, err := RenderNodeServantJob(convertCtx)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if err = kubeutil.RunServantJobs(c.ClientSet, c.WaitServantJobTimeout, func(nodeName string) (*batchv1.Job, error) {
	// 		return nodeservant.RenderNodeServantJob("convert", convertCtx, nodeName)
	// 	}, c.EdgeNodes, os.Stderr, false); err != nil {
	// 		// print logs of yurthub
	// 		for i := range c.EdgeNodes {
	// 			hubPodName := fmt.Sprintf("yurt-hub-%s", c.EdgeNodes[i])
	// 			pod, logErr := c.ClientSet.CoreV1().Pods("kube-system").Get(context.TODO(), hubPodName, metav1.GetOptions{})
	// 			if logErr == nil {
	// 				kubeutil.PrintPodLog(c.ClientSet, pod, os.Stderr)
	// 			}
	// 		}

	// 		// print logs of yurt-manager
	// 		podList, logErr := c.ClientSet.CoreV1().Pods("kube-system").List(context.TODO(), metav1.ListOptions{
	// 			LabelSelector: labels.SelectorFromSet(map[string]string{"app.kubernetes.io/name": "yurt-manager"}).String(),
	// 		})
	// 		if logErr != nil {
	// 			klog.Errorf("failed to get yurt-manager pod, %v", logErr)
	// 			return err
	// 		}

	// 		if len(podList.Items) == 0 {
	// 			klog.Errorf("yurt-manager pod doesn't exist")
	// 			return err
	// 		}
	// 		if logErr = kubeutil.PrintPodLog(c.ClientSet, &podList.Items[0], os.Stderr); logErr != nil {
	// 			return err
	// 		}
	// 		return err
	// 	}
	// }
	return nil
}

// // RenderNodeServantJob return k8s job
// // to start k8s job to run convert/revert on specific node
// func RenderNodeServantJob(renderCtx map[string]string, nodeName string) (*batchv1.Job, error) {
// 	tmplCtx := make(map[string]string)
// 	for k, v := range renderCtx {
// 		tmplCtx[k] = v
// 	}

// 	servantJobTemplate := constants.ConvertServantJobTemplate
// 	jobBaseName := "node-servant-convert"

// 	tmplCtx["jobName"] = jobBaseName + "-" + nodeName
// 	tmplCtx["nodeName"] = nodeName
// 	jobYaml, err := tmplutil.SubsituteTemplate(servantJobTemplate, tmplCtx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	srvJobObj, err := k8s.YamlToObject([]byte(jobYaml))
// 	if err != nil {
// 		return nil, err
// 	}
// 	srvJob, ok := srvJobObj.(*batchv1.Job)
// 	if !ok {
// 		return nil, fmt.Errorf("fail to assert node-servant job")
// 	}

// 	return srvJob, nil
// }

// func validate(action string, tmplCtx map[string]string, nodeName string) error {
// 	if nodeName == "" {
// 		return fmt.Errorf("nodeName empty")
// 	}

// 	switch action {
// 	case "convert":
// 		keysMustHave := []string{"node_servant_image", "yurthub_image", "joinToken"}
// 		return checkKeys(keysMustHave, tmplCtx)
// 	case "revert":
// 		keysMustHave := []string{"node_servant_image"}
// 		return checkKeys(keysMustHave, tmplCtx)
// 	default:
// 		return fmt.Errorf("action invalied: %s ", action)
// 	}
// }

// func checkKeys(arr []string, tmplCtx map[string]string) error {
// 	for _, k := range arr {
// 		if _, ok := tmplCtx[k]; !ok {
// 			return fmt.Errorf("key %s not found", k)
// 		}
// 	}
// 	return nil
// }

// 	// deploy yurt-hub and reset the kubelet service on cloud nodes
// 	convertCtx["working_mode"] = string(util.WorkingModeCloud)
// 	klog.Infof("convert context for cloud nodes(%q): %#+v", c.CloudNodes, convertCtx)
// 	if err = kubeutil.RunServantJobs(c.ClientSet, c.WaitServantJobTimeout, func(nodeName string) (*batchv1.Job, error) {
// 		return nodeservant.RenderNodeServantJob("convert", convertCtx, nodeName)
// 	}, c.CloudNodes, os.Stderr, false); err != nil {
// 		return err
// 	}

// 	klog.Info("If any job fails, you can get job information through 'kubectl get jobs -n kube-system' to debug.\n" +
// 		"\tNote that before the next conversion, please delete all related jobs so as not to affect the conversion.")

// 	return nil
// }

// func prepareYurthubStart(cliSet kubeclientset.Interface, kcfg string) (string, error) {
// 	// prepare kube-public/cluster-info configmap before convert

// 	// // prepare global settings(like RBAC, configmap) for yurthub
// 	// if err := kubeutil.DeployYurthubSetting(cliSet); err != nil {
// 	// 	return "", err
// 	// }

// 	// prepare join-token for yurthub

// 	return joinToken, nil
// }

// // prepareClusterInfoConfigMap will create cluster-info configmap in kube-public namespace if it does not exist
// func prepareClusterInfoConfigMap(client kubeclientset.Interface, file string) error {
// 	info, err := client.CoreV1().ConfigMaps(metav1.NamespacePublic).Get(context.Background(), bootstrapapi.ConfigMapClusterInfo, metav1.GetOptions{})
// 	if err != nil && apierrors.IsNotFound(err) {
// 		// Create the cluster-info ConfigMap with the associated RBAC rules
// 		if err := kubeadmapi.CreateBootstrapConfigMapIfNotExists(client, file); err != nil {
// 			return fmt.Errorf("error creating bootstrap ConfigMap, %w", err)
// 		}
// 		if err := kubeadmapi.CreateClusterInfoRBACRules(client); err != nil {
// 			return fmt.Errorf("error creating clusterinfo RBAC rules, %w", err)
// 		}
// 	} else if err != nil || info == nil {
// 		return fmt.Errorf("fail to get configmap, %w", err)
// 	} else {
// 		klog.V(4).Infof("%s/%s configmap already exists, skip to prepare it", info.Namespace, info.Name)
// 	}
// 	return nil
// }

func nodePoolResourceExists(client kubeclientset.Interface) (bool, error) {
	groupVersion := schema.GroupVersion{
		Group:   "apps.openyurt.io",
		Version: "v1alpha1",
	}
	apiResourceList, err := client.Discovery().ServerResourcesForGroupVersion(groupVersion.String())
	if err != nil && !apierrors.IsNotFound(err) {
		klog.Errorf("failed to discover nodepool resource, %v", err)
		return false, err
	} else if apiResourceList == nil {
		return false, nil
	}

	for i := range apiResourceList.APIResources {
		if apiResourceList.APIResources[i].Name == "nodepools" && apiResourceList.APIResources[i].Kind == "NodePool" {
			return true, nil
		}
	}
	return false, nil
}

func (c *ClusterConverter) deployYurtManager() error {
	err := packages.Install(packages.ClusterBootstrapYurtManager)
	if err != nil {
		return err
	}

	// waiting yurt-manager pod ready
	return wait.PollImmediate(10*time.Second, 5*time.Minute, func() (bool, error) {
		podList, err := c.ClientSet.CoreV1().Pods("kube-system").List(context.TODO(), metav1.ListOptions{
			LabelSelector: labels.SelectorFromSet(map[string]string{"app.kubernetes.io/name": "yurt-manager"}).String(),
		})
		if err != nil {
			klog.Errorf("failed to list yurt-manager pod, %v", err)
			return false, nil
		} else if len(podList.Items) == 0 {
			klog.Infof("no yurt-manager pod: %#v", podList)
			return false, nil
		}

		if podList.Items[0].Status.Phase == corev1.PodRunning {
			for i := range podList.Items[0].Status.Conditions {
				if podList.Items[0].Status.Conditions[i].Type == corev1.PodReady &&
					podList.Items[0].Status.Conditions[i].Status == corev1.ConditionTrue {
					return true, nil
				}
			}
		}
		klog.Infof("pod(%s/%s): %#v", podList.Items[0].Namespace, podList.Items[0].Name, podList.Items[0])
		return false, nil
	})
}
