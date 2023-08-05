---
title: OIDC based kubernetes access
authors:
  - "@batthebee"
reviewers:
  - "@?"
creation-date: 2023-08-03
last-updated: 2023-08-03
status:  experimental
see-also:
  - "/docs/proposals/20230803-user-management-and-access.md"
replaces:
superseded-by:
---

# Technical Proposal: Kubernetes Cluster Access via HashiCorp Boundary with Keycloak OIDC Provider

## 1. Summary

This proposal outlines the integration of HashiCorp Boundary with the edgefarm.core Kubernetes cluster, enhancing security and access control through the use of Keycloak as an OIDC (OpenID Connect) provider. By implementing this solution, we aim to improve authentication and authorization mechanisms while ensuring a streamlined and secure access process to the Kubernetes cluster for authorized edgefarm users.

## 2. Motivation

Currently, when creating users, a kubeconfig must be created manually and made available to the user. This will be automated in the future. The user should be able to request a passed kubeconfig with appropriate rights via his centrally managed credentials.

Traditional methods of accessing Kubernetes clusters often involve complex network configurations and security concerns. HashiCorp Boundary provides a secure and dynamic way to access infrastructure, while Keycloak serves as a reliable OIDC identity provider, offering robust authentication and single sign-on capabilities.

### 2.1 Goals

- **Enhanced Security**: Implement a more secure access control mechanism to the Kubernetes cluster by utilizing HashiCorp Boundary's zero-trust architecture.
- **Simplified Access**: Provide a user-friendly and centralized point of access to the Kubernetes cluster for authorized personnel.
- **Centralized Identity Management**: Utilize Keycloak as an OIDC provider to enable secure, federated, and standards-based authentication.
- **Scalability**: Ensure the solution can handle a growing number of users and applications without compromising performance.

### 2.2 Non-Goals/Future Work

- This proposal does not cover the migration of existing authentication mechanisms to OIDC.
- Extending this setup to include multi-factor authentication (MFA) is considered as future work and is not included in this proposal.
- Creating and managing new roles is not part of this proposal
- Creating a customer friendly cli's will be covered in a separate proposal.
- Connecting to additional systems such as ssh access to the edge nodes is covered in other proposals. This focuses exclusively on authendification/authorization against Kubernetes.

## 3. Proposal

We propose the following steps to implement the integration between HashiCorp Boundary and the Kubernetes cluster using Keycloak as the OIDC provider:

### 3.1 HashiCorp Boundary Setup

Boundary serves as the secure access gateway, governing access to critical resources 
within EdgeFarm. By integrating Boundary into the infrastructure orchestration process, 
we ensure that the provisioning of access points aligns seamlessly with user management 
and access control.

1. **Create Terrajet Crossplane Boundary provider:** A new TerraJet-driven Crossplane 
  provider will facilitate the deployment of Boundary resources within the Kubernetes
  environment.

2. **Deploy HashiCorp Boundary Server:** Create a Helm chart to set up a Boundary server
   to act as an intermediary between users and the Kubernetes cluster.

3. **Define Boundary Host Catalog:** Configure a host catalog within Boundary, using 
the TerraJet provider, that defines the target Kubernetes cluster instances.

4. **Create Boundary Policies:** Extend the Helm chart by defining policies that control user access to the Kubernetes cluster based on identity, roles, and other attributes using the TerraJet provider.

### 3.2 HashiCorp Vault Setup

Vault creates and stores the required kubeconfigs and role bindings and makes them available to the Boundary user. With the integration of Vault, we have a secure and centralized repository for all accesses and can also create accesses with fast expiration dates.

1. **Create Terrajet Crossplane Vault provider:** Modify the existing TerraJet-driven Crossplane provider so that it facilitates the deployment of Vault Kubernetes resources (kubeconfigs and role bindings) within the Kubernetes environment.

2. **Configure Vault Kubernetes Authentication Backend:** Set up the Vault Kubernetes authentication backend to establish a trust relationship between Vault and the Kubernetes cluster. This involves configuring Kubernetes API server details, such as the API server URL and CA certificate, within Vault.

3. **Vault Policies and Kubernetes Roles:** Define Vault policies that correspond to the desired Kubernetes access levels. These policies align with the Boundary-defined roles and access scopes. Map Vault policies to Kubernetes roles, specifying the permissions users or service accounts should have within the cluster.

4. **Dynamic Secrets for Kubeconfigs:** Utilize Vault's dynamic secrets engine to generate and manage Kubernetes kubeconfigs dynamically. When a user or service account is authenticated by Boundary and granted access, Vault generates a short-lived kubeconfig with the necessary credentials. This ensures that access credentials are rotated regularly, enhancing security.

### 3.3 Kubernetes Integration

1. **Kubernetes RBAC Configuration:** Extend the Helm chart to configure Kubernetes Role-Based Access Control (RBAC) rules and policies that align with the roles defined within Boundary and Vault. Define role bindings, cluster roles, and cluster role bindings to enforce fine-grained access controls based on user attributes and policies.

2. **Boundary Credential Management:** Enhance the Helm chart by implementing a robust mechanism for Boundary to securely retrieve Kubernetes credentials and tokens for authenticated users. This mechanism leverages the Vault-generated kubeconfigs and tokens, allowing Boundary to seamlessly deliver the necessary access credentials to users in a secure and controlled manner.

With these steps completed, the HashiCorp Boundary and Vault integration ensures a comprehensive and secure approach to managing Kubernetes cluster access. Users are authenticated through Boundary, and their access is tightly controlled and audited using Vault-generated credentials and policies. This setup not only enhances security but also simplifies the access management process by providing a centralized and automated solution.

### 3.4 Keycloak Integration

1. Create Terrajet Crossplane Keycloak provider: A new TerraJet-driven crossplane provider 
   will facilitate the deployment of Keycloak resources within the Kubernetes environment. 

2. Deploy Keycloak Instance: Create Helm Chart to set up a Keycloak instance that will serve 
   as the OIDC identity provider.

3. Configure Keycloak Realm: Define a realm within Keycloak, specifying the required client 
   and identity provider settings.
   
4. Integrate Boundary with Keycloak: Configure Boundary to use Keycloak as the OIDC provider,
   including client registration and endpoint configurations.

## Conclusion

Integrating HashiCorp Boundary with the Kubernetes cluster using Keycloak as the OIDC provider offers a robust, scalable, and secure solution for managing access. This proposal outlines the necessary steps to achieve this integration, enhancing both the security and accessibility of the Kubernetes environment.

## Implementation History

#- [ ] MM/DD/YYYY: Proposed idea in an issue or [community meeting]
