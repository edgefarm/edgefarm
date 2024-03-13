---
title: User Management and Access
authors:
  - "@batthebee"
reviewers:
  - "@?"
creation-date: 2023-08-03
last-updated: 2023-08-03
status: experimental
see-also:
  - "/docs/proposals/20230803-oidc-kubernetes-access.md"
replaces:
superseded-by:
---

# Technical Proposal: Comprehensive User Management and Access Solution for EdgeFarm

## Introduction

This proposal outlines a robust and comprehensive user management and access solution for EdgeFarm, aiming to streamline user authentication, centralize access control, and enhance security across various components within the system. The solution entails the implementation of a centralized user management system, the integration of key system components, and the establishment of granular role-based access control.

## Objectives

The primary objectives of this proposal are as follows:

1. **Single Sign-On (SSO)**: Implement a Single Sign-On mechanism to facilitate seamless and 
   secure access for users across EdgeFarm components.
2. **Centralized User Management**: Establish a central user management system to 
   efficiently manage user credentials, permissions, and roles.
3. **System Integration**: Integrate key components within the EdgeFarm ecosystem, 
   including Grafana, Backstage, and Kubernetes, with the central user management system.
4. **Default Role Mapping**: Design and implement a default role structure that aligns 
   with the specific functionalities and responsibilities of users within EdgeFarm.
5. **Enhanced Security**: Elevate security levels by enforcing least privilege principles 
   and minimizing the reliance on root access.
6. **Scalability**: Set the foundation for future scalability, paving the way for a 
   multi-user system and potential integration with external identity providers.

## Motivation

The current state of EdgeFarm lacks a unified user administration system, resulting in fragmented access management practices. User-specific credentials are being generated separately for each system, leading to inefficiencies and security vulnerabilities. This necessitates the implementation of a holistic solution.

### Non-Goals/Future Work

#### Future Prospects 

While the proposed solution focuses on immediate priorities, the groundwork will be laid for future enhancements, including:

1. **Multi-User System**: Transitioning towards a multi-user system to accommodate broader user engagement.
2. **Node Access**: Integrating node access via SSH, expanding the scope of the central user management system.
3. **Invitation-Based User Addition**: Implementing an invitation-based system for adding users to EdgeFarm.
4. **External Identity Providers**: Exploring integration with external identity providers like GitHub for seamless user authentication.

#### Non-Goals

1. **Advanced Role Management**: Developing advanced role management capabilities to fine-tune permissions.
2. **Import/Export Interface**: Integration of accessing import/export interfaces of edgefarm.network

## Proposal

### Central User Management System

Keycloak, a battle-tested open-source user management platform, has been selected as the central pillar of this solution. Keycloak's robust feature set, proven reliability, and active community make it the optimal choice for EdgeFarm's requirements.

### System Integration

The following key systems will be seamlessly integrated with the central user management system:

- **Grafana**: Enable users to access Grafana dashboards with their unified EdgeFarm credentials.
- **Backstage**: Provide users with streamlined entry to Backstage's services via centralized authentication.
- **Kubernetes**: Implement single sign-on for Kubernetes, simplifying user access to cluster resources.

### Role-Based Access Control

Roles will be established as distinct groups within Keycloak, aligning with specific user functions. This ensures that each user is granted appropriate permissions based on their responsibilities. Additionally, default role mappings will be defined to swiftly assign necessary rights upon user onboarding.

### Implementation Approach

The integration of EdgeFarm components with the central user management system will be facilitated through the OpenID Connect (OIDC) protocol, guaranteeing secure and standardized communication.

## Implementation History

 - [ ] MM/DD/YYYY: Proposed idea in an issue or [community meeting]

