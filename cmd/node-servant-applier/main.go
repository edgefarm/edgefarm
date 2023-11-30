package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	tmplutil "github.com/openyurtio/openyurt/pkg/util/templates"

	"github.com/edgefarm/edgefarm/cmd/node-servant-applier/internal/constants"
	"github.com/edgefarm/edgefarm/pkg/k8s/tokens"
	"github.com/edgefarm/edgefarm/pkg/k8s/yaml"
	kubeclientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

const (
	yssYurtHubName = "yurt-static-set-yurt-hub"
)

func main() {
	var config string
	flag.StringVar(&config, "config", "", "use kubeconfig if defined. If undefined use Service account auth")
	flag.Parse()
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		log.Fatal(err)
	}

	client, err := kubeclientset.NewForConfig(kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	joinToken, err := tokens.GetOrCreateJoinTokenString(client)
	if err != nil || joinToken == "" {
		log.Fatal(err)
	}

	nodeName := os.Getenv("NODE_NAME")
	nodeServantImage := getEnv("NODE_SERVANT_IMAGE", "ghcr.io/openyurtio/openyurt/node-servant:v1.4.0")
	yurthubImage := getEnv("YURTHUB_IMAGE", "ghcr.io/openyurtio/openyurt/yurthub:v1.4.0")
	workingMode := getEnv("WORKING_MODE", "edge")
	enableDummyIf := getEnv("ENABLE_DUMMY_IF", "true")

	convertCtx := map[string]string{
		"node_servant_image": nodeServantImage,
		"yurthub_image":      yurthubImage,
		"joinToken":          joinToken,
		// The node-servant will detect the kubeadm_conf_path automatically
		// It will be either "/usr/lib/systemd/system/kubelet.service.d/10-kubeadm.conf"
		// or "/etc/systemd/system/kubelet.service.d/10-kubeadm.conf".
		"kubeadm_conf_path": "",
		"working_mode":      workingMode,
		"enable_dummy_if":   enableDummyIf,
		"configmap_name":    yssYurtHubName,
		"nodeName":          nodeName,
		"jobName":           "node-servant-convert-" + nodeName,
		"enable_node_pool":  "true",
	}
	klog.Infof("convert context for edge node %s: %#+v", nodeName, convertCtx)

	job, err := RenderNodeServantJob(convertCtx)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := client.BatchV1().Jobs(job.GetNamespace()).Get(context.Background(), job.GetName(), metav1.GetOptions{}); err == nil {
		log.Printf("job %s already exists", job.GetName())
		os.Exit(0)
	}
	_, err = client.BatchV1().Jobs(job.GetNamespace()).Create(context.Background(), job, metav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

// RenderNodeServantJob return k8s job
// to start k8s job to run convert/revert on specific node
func RenderNodeServantJob(renderCtx map[string]string) (*batchv1.Job, error) {
	servantJobTemplate := constants.ConvertServantJobTemplate
	jobYaml, err := tmplutil.SubsituteTemplate(servantJobTemplate, renderCtx)
	if err != nil {
		return nil, err
	}

	jobObj, err := yaml.YamlToObject([]byte(jobYaml))
	if err != nil {
		return nil, err
	}
	job, ok := jobObj.(*batchv1.Job)
	if !ok {
		return nil, fmt.Errorf("fail to assert node-servant job")
	}

	return job, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
