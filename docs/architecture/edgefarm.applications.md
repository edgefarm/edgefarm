![](../assets/architecture-edgefarm.applications.png)

## Initial state

### 0.

KubeVela is configured and up and running and waits on application manifests.

## User wants to deploy an Application

### 1.

User apply his application, defined by a manifest against the kubernetes api.

KubeVela receives the manifest and converts the definition into argo rollouts ressources.
This ressources are pushed again against the kubernetes api and kubernetes creates the workload on
the specific nodes.

## User wants to deploy the Application in GitOps style

### 2.

The user needs to tell Argo CD how to access his git repository.

From now on, Argo CD watches the repository for changes. If there any chages, argocd will pickup
the ressouces and push them against the kubernetes api.

### 3.

The user commits his application mainfest and Argo CD pushes the ressources.

The remaining processing is equal to 1.
