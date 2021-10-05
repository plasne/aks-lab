# Deploy Azure Resources
You will need the following tools. If you do not want to install these tools locally, you could also install them on a Linux or Windows VM in Azure.

You will need to install Docker Community Edition:
* https://docs.docker.com/install

It is not required to have any of these tools installed, but you could install them if you want to run some of this locally:
* Azure CLI 2.x: https://docs.microsoft.com/en-us/cli/azure/install-azure-cli?view=azure-cli-latest
* kubectl: `az aks install-cli` - Kubernetes command-line utility.
* istioctl: https://istio.io/latest/docs/setup/getting-started/#download - Istio configuration command-line utility

Azure CLI 2.x and kubectl are already installed in Cloud Shell (https://azure.microsoft.com/en-us/features/cloud-shell/) so you could use that.

In each section below, the az-cli commands to deploy the resource are hidden. You can open that section and run those commands, but you might try to do it on your own first.

&nbsp;

## Provision an Azure Container Registry

Provision an "Azure Container Registry" in Azure using the default parameters:
* SKU: Standard

Ref:
* https://docs.microsoft.com/en-us/azure/container-registry/container-registry-get-started-portal
* https://docs.microsoft.com/en-us/azure/container-registry/container-registry-get-started-azure-cli
* https://docs.microsoft.com/en-us/azure/container-registry/container-registry-get-started-powershell

<details>
  <summary>Provision ACR</summary>

```bash
# Set up the following variables (configure as needed)
SUBSCRIPTION=<your subscription Id or name>
REGION_NAME=eastus
RESOURCE_GROUP=akslabhv-rg
ACR_NAME=akslabhv
ACR_SKU=Standard

# Login to Azure
az login

# Set your default subscription
az account set -s $SUBSCRIPTION

# Confirm it is set correctly
az account show

# Create resource group
az group create --name $RESOURCE_GROUP --location $REGION_NAME 

# Create Azure container registry
az acr create --resource-group $RESOURCE_GROUP --name $ACR_NAME --sku $ACR_SKU
```

</details>

&nbsp;

## Provision a Kubernetes Cluster

Provision an "Azure Kubernetes Service (AKS)" cluster in Azure using the following configuration:

* Node count: 3
* Node virtual machine size: 3x Standard D2 v2
* Integrated with the ACR you created above (so that an AcrPull role is set for the cluster's managed identity)

Ref:

* https://docs.microsoft.com/en-us/azure/aks/kubernetes-walkthrough-portal
* https://docs.microsoft.com/en-us/azure/aks/kubernetes-walkthrough
* https://docs.microsoft.com/en-us/azure/virtual-machines/linux/mac-create-ssh-keys
* https://docs.microsoft.com/en-us/azure/virtual-machines/linux/ssh-from-windows
* https://docs.microsoft.com/en-us/azure/aks/use-managed-identity
* https://docs.microsoft.com/en-us/azure/aks/cluster-container-registry-integration?tabs=azure-cli
* https://docs.microsoft.com/en-us/azure/container-registry/container-registry-roles?tabs=azure-cli

<details>
  <summary>Provision AKS</summary>

```bash
# Set up the following variables (configure as needed)
RESOURCE_GROUP=akslabhv-rg # created above
ACR_NAME=akslabhv # created above
CLUSTER_NAME=akslabhv
ISTIO_VERSION=1.11.3
NODE_COUNT=3
NODE_VM_SIZE=Standard_DS2_v2

# Create cluster
az aks create --resource-group $RESOURCE_GROUP --name $CLUSTER_NAME --node-count $NODE_COUNT \
    --node-vm-size $NODE_VM_SIZE --generate-ssh-keys --enable-managed-identity

# Integrate ACR
az aks update --attach-acr $ACR_NAME

# Get aks credentials to use kubectl
az aks get-credentials --resource-group $RESOURCE_GROUP --name $CLUSTER_NAME
```

</details>

&nbsp;

## Install Istio

Install Istio on the AKS cluster.

Ref:

* https://istio.io/latest/docs/

<details>
  <summary>Install Istio</summary>

Download istioctl

* MacOS or Linux:

  ```bash
  # This will download version 1.11.3
  curl -L https://istio.io/downloadIstio | ISTIO_VERSION=1.11.3 sh -

  # Navigate to the istio package directory
  cd istio-1.11.3

  # Add the istioctl client to your path
  export PATH=$PWD/bin:$PATH
  ```

* Windows:

  Download [version 1.11.3](https://github.com/istio/istio/releases/tag/1.11.3) and add  _`<your_path_to_istio_directory>/istio-1.11.3/bin`_ to your Path

Provision resources

```bash
# Install istio to your cluster
istioctl install --set profile=minimal -y

# Add a namespace label to instruct Istio to automatically inject Envoy sidecar proxies when you deploy your application later
kubectl label namespace default istio-injection=enabled
```

</details>

&nbsp;

