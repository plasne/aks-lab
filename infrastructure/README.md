# Infrastructure Overview - WIP & merge to main doc

The Infrastructure that needs to be provisioned for this lab includes:

- [Azure Container Registry](https://docs.microsoft.com/en-us/azure/container-registry/) - registry for your container images and helm charts
- [Azure Kubernetes Cluster](https://docs.microsoft.com/en-us/azure/aks/)
- [Istio](https://istio.io/latest/) - service mesh

## Log into Azure

Using the Azure CLI ensure you are logged into the right subscription

```bash
# Set a variable for your subscription
SUBSCRIPTION=<your subscription Id or name>

# Login to Azure
az login

# Set your default subscription
az account set -s $SUBSCRIPTION

# Confirm it is set correctly
az account show
```

## Download istioctl

On MacOS or Linux

```bash
# This will download version 1.11.3
curl -L https://istio.io/downloadIstio | ISTIO_VERSION=1.11.3 sh -

# Navigate to the istio package directory
cd istio-1.11.3

# Add the istioctl client to your path
export PATH=$PWD/bin:$PATH
```

On Windows:

- Download [version 1.11.3](https://github.com/istio/istio/releases/tag/1.11.3) and add  _`<your_path_to_istio_directory>/istio-1.11.3/bin`_ to your Path environment variable.

## Provision you resources

```bash
# Navigate to the infrastructure directory in this repo
cd infrastructure

# Run the following script
. provision_infra.sh
```
