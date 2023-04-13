# Getting Started

## Overview

The Edgefarm is built up declaratively according to GitOps standard. Whenever you want to add or change something, a manifest must be created or edited in Yaml format. These manifests are usually stored in a Git repository.

## Terminology

Here we explain some of the terms that are important for Edgefarm

### Node

A node represents a Kubernetes node. Usually this is the edge instance running on one of your remote devices and is connected to the main node.

### Network

We have developed a sophisticated network layer that connects the main node and all edge nodes. Such a network contains one or more sub-networks, one or more NATS streams and one or more users.
All these components together form a network

### Application

This important part of the Edgefarm is your application, which you can deploy into the Edgefarm. How such an application looks like is up to you. Mostly it will be an application that reads data from the remote device and sends it over the network to the cloud.

## Portal

For the purpose of simplifying the configuration of your Edgefarm, we have developed a portal that will help you with the initial setup of your Edgefarm.

You can find this portal at https://portal.edgefarm.io

### The portal simply explained

Portal is designed to make it as easy as possible for you to create new EdgeFarm components and to give you an overview of your EdgeFarm. In the portal you will find workflows to create new components as well as data and metrics about your existing components.

### How to start?

Quick overview:

1. Create a system
2. Add node(s)
3. Create network(s)
4. Add application(s)

### 1. System

First, you start by creating a system in the portal. The system is a kind of project or workspace that is specifically linked to a Git repository where your manifest files will be stored.

To do this, navigate to https://portal.edgefarm.io/create and select the workflow to create a new system.

In the following a wizard will ask you for some information. Fill it out and click on 'Create' at the end of the wizard. After completing this wizard, a new Git repository is created, whose link you can find in the summary of the wizard. Also, an entry has now been created in the portal to represent your system.

### 2. (Edge-) Nodes

Next, you have to make the (edge-) nodes known to the system. There is also a predefined workflow for this. Navigate to https://portal.edgefarm.io/create and select the workflow which creates a new node.

As before, follow the wizard to the end. This workflow has now opened a pull request in your Git repository, which was created by building the system. Click on the link in the wizard summary or navigate to your pull request from the Git interface to merge it. Only after the merge of the pull request the node can be visible in the portal.

Note: After the merge it may take a few minutes until the node is visible in the portal.

To display all nodes known to the portal, navigate to the menu item 'Nodes' in the page navigation.

### 3. Networks

To ensure basic communication with the core node, you must now create at least one network. This process will drop a Kubernetes custom resource as a manifest file on the already mentioned Git repository and deploy it to the cluster.
To do this, navigate to https://portal.edgefarm.io/create and select the workflow that adds a network. Follow the steps of the wizards to the end. Similar to creating a node, at the end the created pull request must be remembered before the network can be displayed in Portal.

To display all your networks that are known to the portal, navigate to the menu item 'Networks' in the page navigation.

### 4. Applications

Comming soon ...
