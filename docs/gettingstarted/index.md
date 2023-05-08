# Getting Started

Welcome to the Getting Started page for EdgeFarm.

## Overview

Here you'll learn about the terminology and how to manage EdgeFarm components like Edge Nodes, Applications, Networks using the EdgeFarm Portal.

As a side note, EdgeFarm is GitOps ready. Commit your changes to your Git repository and let the CD system do the rest. You find more information about GitOps [here](https://www.cncf.io/blog/2021/09/28/gitops-101-whats-it-all-about/). 

## Basic Terminology

Let's talk about basic terminology. 

### Edge Node

A node represents a Kubernetes node. Thus, an Edge Node is a remote device e.g. Raspbrry Pi connected to the Kubernetes cluster running your workload. Edge Nodes are managed by `EdgeFarm.core`. 

### Application

Writing workload definitions can be hard. To make this more convinient there is `EdgeFarm.applications`. It allows you to define your workload in a very minimalist format that can be extended to your needs. Using `EdgeFarm.applications` you can roll out your custom OCI images to your Edge Nodes.

### Network

How to transfer application data? For this purpose `EdgeFarm.network` was developed. Create a Network and let your applications communicate no matter if running in Cloud, Edge or even exported to a third party system.

## Portal

You can find the EdgeFarm Portal here: https://go.edgefarm.io

### The portal simply explained

The Portal is designed to make it as easy as possible to create new EdgeFarm components and to give an overview of the cluster. In the portal you will find workflows to create new components as well as information and metrics about existing components.
Creating components, no matter which component, results in a code change in your upstream Git Repository - remember GitOps, right?

### How to start?

Quick overview:

1. Create a system
2. Add edge node(s)
3. Create network(s)
4. Create application(s)

### 1. System

First, you start by creating a `system` in the portal. The `system` is a kind of project or workspace that is specifically linked to a Git repository where your manifest files will be stored.

To do this, navigate to https://go.edgefarm.io/create and select the workflow to create a new `system`.

A wizard will ask you for some information. Fill it out and click on 'Create'. After completing this wizard, a new Git repository is created, whose link you can find in the summary of the wizard. Also, an entry has been created in the portal to represent your `system`.

### 2. Edge Nodes

Next, you have to make the `edge nodes` known to the `system`. There is also a predefined workflow for this. Navigate to https://go.edgefarm.io/create and select the workflow which creates a new node.

As before follow the wizard to the end. This workflow has now opened a pull request in your Git repository, which was previously created byt the `system` component. Click on the link in the wizard summary or navigate to your pull request from the Git interface to merge it. After the merge of the pull request the node can be viewed in the portal.

Note: After the merge it may take a few minutes until the node is visible in the portal.

To display all `edge nodes` known to the portal, navigate to the menu item 'Nodes' in the page navigation.

### 3. Networks

As mentioned before, a `network` is used to let `applications` communicate. Define which `edge node` shall be part of the `network` and define the `streams` and `users` that belong to that network. `users` have specific user defined rights to interact with `streams`. `streams` are used to either buffer your data e.g. during poor network connections or aggregate other `streams` to collect the data of many `streams`.

Note: if your application doesn't need to communicate or send data to other applications, you might want to skip the creation of the `network`. 

To create a `network`, navigate to https://go.edgefarm.io/create and select the workflow that adds a `network`. Follow the steps of the wizards to the end. Merge the Pull-Request and see your `network` ocurring in the `Networks` section in the portal.

### 4. Applications

Define an `application` by providing a name and a OCI image. You can customize your application spec to your needs by e.g. adding envs, volumes, custom commands and args, ...

If your application shall be allowed to communicate with a `network` define the network specific parts of the application.

To create an `application`, navigate to https://go.edgefarm.io/create and select the workflow that adds a `application`. Follow the steps of the wizards to the end. Merge the Pull-Request and see your `application` ocurring in the `Applications` section in the portal.