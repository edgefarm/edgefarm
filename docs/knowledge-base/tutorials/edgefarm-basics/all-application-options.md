# Deploying an application

## Objectives

* Learn abiout the EdgeFarm application model
* Deploy your first application with kubectl

## EdgeFarm application model

Once you have a [running EdgeFarm cluster](../../cluster/create-local-cluster/), you can deploy applications to it. To do so, you create a edgefarm.application resource. This resource is a custom resource definition (CRD) that is specific to EdgeFarm. It defines a set of Kubernetes resources that are deployed together as a single unit. The edgefarm.application resource is the primary resource in the EdgeFarm application model, and it represents a single instance of your application. Once you've created an edgefarm.application resource, the Kubernetes control plane schedules the application's Pods to run on your selected nodes.

## Writing the manifest

The edgefarm.application resource is a Kubernetes custom resource. You can create it by writing a manifest file that describes the resource. The manifest file is a YAML file that contains the edgefarm.application resource definition. The following example shows a manifest file that creates an edgefarm.application resource. 
See the [edgefarm.application reference](../../../reference/reference/api/applications/overview) for a complete description of the edgefarm.application resource.

Let's create a file called `basic.yaml` and add the following content:
```yaml
apiVersion: core.oam.dev/v1beta1
kind: Application
metadata:
  name: kubernetes-bootcamp
  namespace: default
spec:
  components:
    - name: kubernetes-bootcamp
      type: edgefarm-applications 
      properties:
        image: gcr.io/google-samples/kubernetes-bootcamp:v1 #(1)!
        nodepoolSelector:
          matchLabels:
            bootcamp: "true" #(2)!
```

1. We'll deploy the kubernetes-bootstrap OCI image. This is a simple webserver that will be deployed to the cluster.
2. We'll deploy to nodepools that have the label `bootcamp=true`. This label will be added in the next few steps.

## Deploying the manifest

To view the Edge Nodes in the cluster, run the `kubectl get nodes -l openyurt.io/is-edge-worker=true` command.

You see the available edge nodes. Later, we will choose where to deploy our application based on Node available resources.

```console
$ kubectl get nodes -l openyurt.io/is-edge-worker=true
NAME               STATUS   ROLES    AGE    VERSION
edgefarm-worker2   Ready    <none>   9m2s   v1.22.7
edgefarm-worker3   Ready    <none>   9m3s   v1.22.7

$ kubectl get nodepools                 
NAME               TYPE   READYNODES   NOTREADYNODES   AGE
edgefarm-worker2   Edge   1            0               9m5s
edgefarm-worker3   Edge   1            0               9m5s
```

!!! note
    Every Edge Node is mapped to a corresponding nodepool. This can be seen as a 1:1 relationship. Via labels on the nodepool, we can control which applications are deployed to which node.

Now let's deploy the application and label a nodepool as `bootcamp=true`.

```console
$ kubectl apply -f basic.yaml #(1)!
application.core.oam.dev/kubernetes-bootcamp created

$ kubectl label nodepools.apps.openyurt.io edgefarm-worker3 bootcamp=true #(2)!
nodepool.apps.openyurt.io/edgefarm-worker3 labeled

$ kubectl get deployments.apps #(3)!
NAME                                         READY   UP-TO-DATE   AVAILABLE   AGE
kubernetes-bootcamp-edgefarm-worker3-8krt7   1/1     1            1           21s

$ kubectl get pods -o wide #(4)!                   
NAME                                                          READY   STATUS    RESTARTS   AGE   IP           NODE              NOMINATED NODE   READINESS GATES
kubernetes-bootcamp-edgefarm-worker3-8krt7-6b4fc49596-56f2h   1/1     Running   0          37s   10.244.2.5   edgefarm-worker3   <none>           <none>
```

1. We deploy the application
2. We label the nodepool `edgefarm-worker3` as `bootcamp=true`
3. We see that the deployment on node `edgefarm-worker3` was successful
4. We see that the pod is running on node `edgefarm-worker3`

## Testing the application

Now, let's test the application. We'll use the `kubectl exec` command to run a command in the pod and print its output.

```console
$ kubectl exec -it kubernetes-bootcamp-edgefarm-worker3-8krt7-6b4fc49596-56f2h -- curl http://localhost:8080/version
Hello Kubernetes bootcamp! | Running on: kubernetes-bootcamp-edgefarm-worker3-8krt7-6b4fc49596-56f2h | v=1
```

Great! The application is running. You've learned how to deploy an application and to control where it will be deployed.