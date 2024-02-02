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
      - name: node-convert-script
        configMap:
          defaultMode: 420
          name: node-convert-script
      serviceAccount: node-servant-convert
      containers:
      - name: node-servant
        image: {{.node_servant_image}}
        imagePullPolicy: IfNotPresent
        command:
        - sh 
        - -c
        args: 
        - cp /script/run.sh /run.sh && chmod +x /run.sh && /run.sh
        securityContext:
          privileged: true
        volumeMounts:
        - mountPath: /openyurt
          name: host-root
        - mountPath: /openyurt/data
          name: configmap
        - mountPath: /script
          name: node-convert-script
        env:
        - name: NODE_NAME
          valueFrom:
            fieldRef:
              fieldPath: spec.nodeName
          {{if  .kubeadm_conf_path }}
        - name: KUBELET_SVC
          value: {{.kubeadm_conf_path}}
          {{end}}
        - name: WORKING_MODE
          value: {{.working_mode}}
        - name: YURTHUB_IMAGE
          value: {{.yurthub_image}}
        - name: JOIN_TOKEN
          value: {{.joinToken}}`
)
