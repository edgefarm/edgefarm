![](../assets/architecture-edgefarm.core.png)

## Initial state

### 0.

cloudcore is set up, valid certificates are stored, and is waiting for edge nodes.

## The user wants to register a new device

### 1.

The user issue a new node token from vault. This token is only valid for a specific device.

### 2.

The user transfers the token to the NodeRegistration service of the device, whereupon the service has a certificate issued by vault and renews it again and again. The certificates have a very short validity.

### 3.

The certificate is transferred to the egdecore. The edgecore connects to the certificate at the CloudCore and synchronizes from there on via mtls.

## The user wants to deploy workload on edge

### 4.

After going through the previous steps, it is now possible to deploy workloads to edge devices using standard Kubernetes tools.
