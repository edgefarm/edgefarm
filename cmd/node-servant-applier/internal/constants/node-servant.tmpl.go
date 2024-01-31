package constants

const (
	ConvertServantJobTemplate = `apiVersion: batch/v1
kind: Job
metadata:
  name: {{.jobName}}
  namespace: kube-system
spec:
  ttlSecondsAfterFinished: 10
  template:
    spec:
      hostPID: true
      hostNetwork: true
      restartPolicy: OnFailure
      nodeName: {{.nodeName}}
      tolerations:
      - key: edgefarm.io
        effect: NoSchedule
      volumes:
      - name: host-root
        hostPath:
          path: /
          type: Directory
      - name: configmap
        configMap:
          defaultMode: 420
          name: {{.configmap_name}}
      serviceAccount: node-servant-convert
      containers:
      - name: node-servant
        image: {{.node_servant_image}}
        imagePullPolicy: IfNotPresent
        command:
        - /bin/sh
        - -c
        args:
        - 'TOKEN=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token) && apk add curl && /usr/local/bin/entry.sh convert --working-mode={{.working_mode}} --yurthub-image={{.yurthub_image}} {{if .yurthub_healthcheck_timeout}}--yurthub-healthcheck-timeout={{.yurthub_healthcheck_timeout}} {{end}}--join-token={{.joinToken}} {{if .enable_dummy_if}}--enable-dummy-if={{.enable_dummy_if}}{{end}} {{if .enable_node_pool}}--enable-node-pool={{.enable_node_pool}}{{end}} && curl -k -X PATCH https://$KUBERNETES_SERVICE_HOST:$KUBERNETES_SERVICE_PORT_HTTPS/api/v1/nodes/$NODE_NAME -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/merge-patch+json" --data "{\"metadata\":{\"labels\":{\"node.edgefarm.io/converted\":\"true\",\"node.edgefarm.io/to-be-converted\":\"false\"}}}"'
        securityContext:
          privileged: true
        volumeMounts:
        - mountPath: /openyurt
          name: host-root
        - mountPath: /openyurt/data
          name: configmap
        env:
        - name: NODE_NAME
          valueFrom:
            fieldRef:
              fieldPath: spec.nodeName
          {{if  .kubeadm_conf_path }}
        - name: KUBELET_SVC
          value: {{.kubeadm_conf_path}}
          {{end}}
`
)
