package constants

const (

	// ConvertServantJobTemplate defines the node convert servant job in yaml format
	ConvertServantJobTemplate = `
apiVersion: batch/v1
kind: Job
metadata:
  name: {{.jobName}}
  namespace: kube-system
spec:
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
      containers:
      - name: node-servant-servant
        image: {{.node_servant_image}}
        imagePullPolicy: IfNotPresent
        command:
        - /bin/sh
        - -c
        args:
        - "/usr/local/bin/entry.sh convert --working-mode={{.working_mode}} --yurthub-image={{.yurthub_image}} {{if .yurthub_healthcheck_timeout}}--yurthub-healthcheck-timeout={{.yurthub_healthcheck_timeout}} {{end}}--join-token={{.joinToken}} {{if .enable_dummy_if}}--enable-dummy-if={{.enable_dummy_if}}{{end}} {{if .enable_node_pool}}--enable-node-pool={{.enable_node_pool}}{{end}}"
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
