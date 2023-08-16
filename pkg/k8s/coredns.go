/*
Copyright © 2023 EdgeFarm Authors

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

package k8s

import (
	"context"
	"log"

	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8syaml "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

const (
	coreDNSDaemonSet = `apiVersion: apps/v1
kind: DaemonSet
metadata:
  labels:
    k8s-app: kube-dns
  name: coredns
  namespace: kube-system
spec:
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      k8s-app: kube-dns
  template:
    metadata:
      labels:
        k8s-app: kube-dns
    spec:
      containers:
      - args:
        - -conf
        - /etc/coredns/Corefile
        image: k8s.gcr.io/coredns/coredns:v1.8.4
        imagePullPolicy: IfNotPresent
        livenessProbe:
          failureThreshold: 5
          httpGet:
            path: /health
            port: 8080
            scheme: HTTP
          initialDelaySeconds: 60
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 5
        name: coredns
        ports:
        - containerPort: 53
          hostPort: 53
          name: dns
          protocol: UDP
        - containerPort: 53
          hostPort: 53
          name: dns-tcp
          protocol: TCP
        - containerPort: 9153
          hostPort: 9153
          name: metrics
          protocol: TCP
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /ready
            port: 8181
            scheme: HTTP
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        resources:
          limits:
            memory: "170Mi"
          requests:
            cpu: "100m"
            memory: "70Mi"
        securityContext:
          allowPrivilegeEscalation: false
          capabilities:
            add:
            - NET_BIND_SERVICE
            drop:
            - all
          readOnlyRootFilesystem: true
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /etc/coredns
          name: config-volume
          readOnly: true
      dnsPolicy: "None"
      dnsConfig:
        nameservers:
          - 8.8.8.8
          - 8.8.4.4
      nodeSelector:
        kubernetes.io/os: linux
      priorityClassName: system-cluster-critical
      restartPolicy: Always
      schedulerName: default-scheduler
      serviceAccount: coredns
      serviceAccountName: coredns
      terminationGracePeriodSeconds: 30
      tolerations:
      - key: CriticalAddonsOnly
        operator: Exists
      - effect: NoSchedule
        key: node-role.kubernetes.io/master
      - effect: NoSchedule
        key: node-role.kubernetes.io/control-plane
      - effect: NoSchedule
        key: edgefarm.io
      volumes:
      - configMap:
          defaultMode: 420
          items:
          - key: Corefile
            path: Corefile
          name: coredns
        name: config-volume`
)

func convertYamlToDaemonSet(manifest string) (*unstructured.Unstructured, error) {
	ds := &appsv1.DaemonSet{}
	err := yaml.Unmarshal([]byte(manifest), ds)
	if err != nil {
		return nil, err
	}

	// Decode the DaemonSet manifest into an unstructured object
	obj := &unstructured.Unstructured{}
	dec := k8syaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	_, _, err = dec.Decode([]byte(manifest), nil, obj)
	if err != nil {
		log.Fatal(err)
	}

	return obj, nil
}

// ReplaceCoreDNS deletes the CoreDNS deployment and replace it with a DaemonSet
func ReplaceCoreDNS() error {
	clientset, err := GetClientset(nil)
	if err != nil {
		return err
	}

	err = clientset.AppsV1().Deployments("kube-system").Delete(context.Background(), "coredns", metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	err = Apply(coreDNSDaemonSet)
	if err != nil {
		return err
	}
	return nil
}