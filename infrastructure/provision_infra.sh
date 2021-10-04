#!/bin/bash

# Set up the following variables (configure as needed)
REGION_NAME=eastus
RESOURCE_GROUP=akslabhv-rg
ACR_NAME=akslabhv
CLUSTER_NAME=akslabhv
ISTIO_VERSION=1.11.3

# Create resource group
az group create --name $RESOURCE_GROUP --location $REGION_NAME 

# Create Azure container registry
az acr create --resource-group $RESOURCE_GROUP --name $ACR_NAME --sku Standard

# Create cluster
az aks create --resource-group $RESOURCE_GROUP --name $CLUSTER_NAME --node-count 3 \
    --generate-ssh-keys  --attach-acr $ACR_NAME

# Get aks credentials to use kubectl
az aks get-credentials --resource-group $RESOURCE_GROUP --name $CLUSTER_NAME

# Deploy istio to your cluster
istioctl install --set profile=minimal -y

# Add a namespace label to instruct Istio to automatically inject Envoy sidecar proxies when you deploy your application later
kubectl label namespace default istio-injection=enabled