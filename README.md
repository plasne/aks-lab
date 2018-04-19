# AKS Lab
A short lab for getting an application up and running using Azure Container Service (AKS).

While it is not required, it is probably easier to debug your code locally before pushing it out in a container to Kubernetes. You can install Node.js here: https://nodejs.org/en/. 

## Provision a Kubernetes Cluster

Provision an "Azure Container Service (AKS)" cluster in Azure using the default parameters:
* Node count: 3
* Node virtual machine size: 3x Standard D2 v2

While the resource provisions, you may move on to provisioning the Azure Container Registry.

Ref:
* https://docs.microsoft.com/en-us/azure/aks/kubernetes-walkthrough-portal
* https://docs.microsoft.com/en-us/azure/aks/kubernetes-walkthrough
* https://docs.microsoft.com/en-us/azure/virtual-machines/linux/mac-create-ssh-keys
* https://docs.microsoft.com/en-us/azure/virtual-machines/linux/ssh-from-windows
* https://docs.microsoft.com/en-us/azure/container-service/kubernetes/container-service-kubernetes-service-principal

## Provision an Azure Container Registry

Provision an "Azure Container Registry" in Azure using the default parameters:
* Admin user: Enable
* SKU: Standard

Ref:
* https://docs.microsoft.com/en-us/azure/container-registry/container-registry-get-started-portal
* https://docs.microsoft.com/en-us/azure/container-registry/container-registry-get-started-azure-cli
* https://docs.microsoft.com/en-us/azure/container-registry/container-registry-get-started-powershell

## Deploy a web front end

Perform the following steps to deploy a simple "Hello World" application:
1. Write a hello world Node.js application that will expose an HTTP endpoint on port 80.
2. 

Ref:
* https://expressjs.com/en/starter/hello-world.html

>! bash
>! const express = require("express");
>! ```

## Expose the container

## Add a Postgres database

## Add a content creation service

## Deploy a new version of the web front end

## Scale the content creation service
