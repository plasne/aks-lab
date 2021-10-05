# AKS Lab
A short lab for getting an application up and running using Azure Container Service (AKS).

The specific code and commands are collapsed. The intent of this lab is for you to use the "Ref" section to figure out how to do the steps. You can use the collapsed segments AFTER you review the documentation if you get stuck or to verify you did it correctly.

If you do not want to install these tools locally, you could also install them on a Linux or Windows VM in Azure.

You will need to install Docker Community Edition:
* https://docs.docker.com/install

It is not required to have any of these tools installed, but you could install them if you want to run some of this locally:
* Azure CLI 2.0: https://docs.microsoft.com/en-us/cli/azure/install-azure-cli?view=azure-cli-latest
* Node.js: https://nodejs.org/en/
* kubectl: https://kubernetes.io/docs/tasks/tools/install-kubectl/
* istioctl: https://istio.io/latest/docs/setup/getting-started/#download - Istio configuration command line utility

Azure CLI 2.0 and kubectl are already installed in Cloud Shell (https://azure.microsoft.com/en-us/features/cloud-shell/) so you could use that.

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

Provision an "Azure Container Service (AKS)" cluster in Azure using the following configuration:

* Node count: 3
* Node virtual machine size: 3x Standard D2 v2
* Integrated with the ACR you created above (so that an AcrPull role is set for the cluster's managed identity)

Install istio on your cluster.

While the resource provisions, you may move on to provisioning the Azure Container Registry.

Ref:

* https://docs.microsoft.com/en-us/azure/aks/kubernetes-walkthrough-portal
* https://docs.microsoft.com/en-us/azure/aks/kubernetes-walkthrough
* https://docs.microsoft.com/en-us/azure/virtual-machines/linux/mac-create-ssh-keys
* https://docs.microsoft.com/en-us/azure/virtual-machines/linux/ssh-from-windows
* https://docs.microsoft.com/en-us/azure/aks/use-managed-identity
* https://docs.microsoft.com/en-us/azure/aks/cluster-container-registry-integration?tabs=azure-cli
* https://docs.microsoft.com/en-us/azure/container-registry/container-registry-roles?tabs=azure-cli
* https://istio.io/latest/docs/

<details>
  <summary>Provision AKS</summary>

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

# Install istio to your cluster
istioctl install --set profile=minimal -y

# Add a namespace label to instruct Istio to automatically inject Envoy sidecar proxies when you deploy your application later
kubectl label namespace default istio-injection=enabled
```

</details>

&nbsp;

## Deploy a web front end application

Perform the following steps to deploy a simple "Hello World" application:

1. Write a Hello World Node.js application that will expose an HTTP endpoint on port 80.

Ref:
* https://expressjs.com/en/starter/hello-world.html

<details>
  <summary>Node.js Hello World Code Sample</summary>

If you have Node installed, you can provision a new app and install Express:

```bash
npm init
npm install express --save
```

Alternatively, you can simply create the package.json file manually:

```bash
{
  "name": "hello",
  "version": "1.0.0",
  "description": "",
  "main": "server.js",
  "dependencies": {
    "express": "^4.16.3"
  },
  "devDependencies": {},
  "scripts": {
    "test": "echo \"Error: no test specified\" && exit 1",
    "start": "node server.js"
  },
  "author": "",
  "license": "ISC"
}
```

The server.js file could look something like this:

```javascript
const express = require("express");
const app = express();

app.get("/", (req, res) => {
  res.send("Hello World!\n");
});

const port = process.env.PORT || 8800;
app.listen(port, () => {
  console.log(`listening on port ${port}...`);
});
```

You could test locally by:

```bash
node server.js
curl http://localhost:8800
```

</details>

&nbsp;

2. Write a Dockerfile.

Ref:
* https://docs.docker.com/engine/reference/builder
* https://hub.docker.com/_/node

<details>
  <summary>Dockerfile Sample</summary>
    
```Dockerfile
FROM node:latest
COPY server.js server.js
COPY package.json package.json
RUN npm install
ENV PORT 80
EXPOSE 80
CMD node server.js
```
    
</details>

&nbsp;

3. Build the container image.

Ref:
* https://docs.docker.com/engine/reference/commandline/build
* https://docs.docker.com/engine/reference/run

<details>
  <summary>Build Commands</summary>

You can build and view the built images by:

```bash
docker build -t hello:latest -t hello:1.0.0 .
docker images
```

You can test locally by:

```bash
docker run -d --name hello --publish 8800:80 hello:latest
curl http://localhost:8800
```

</details>

&nbsp;

4. Push the container to ACR.

Ref:
* https://docs.microsoft.com/en-us/azure/aks/tutorial-kubernetes-prepare-acr
* https://docs.microsoft.com/en-us/azure/aks/tutorial-kubernetes-deploy-cluster#configure-acr-authentication
* https://docs.docker.com/engine/reference/commandline/tag
* https://docs.docker.com/engine/reference/commandline/push

<details>
  <summary>Push Commands</summary>

```bash
az login
az acr login --name whatever
docker tag hello:1.0.0 whatever.azurecr.io/hello:1.0.0
docker tag hello:1.0.0 whatever.azurecr.io/hello:latest
docker images
docker push whatever.azurecr.io/hello:latest
docker push whatever.azurecr.io/hello:1.0.0
az acr repository list --name whatever --output table
az acr repository show-tags --name whatever --repository hello --output table
```

Alternatively, you can login to your ACR like this:

```bash
docker login whatever.azurecr.io -u whatever -p password
```

</details>

&nbsp;

5. Build a YAML file for the deployment.

Ref:
* https://v1-8.docs.kubernetes.io/docs/concepts/workloads/controllers/deployment/#creating-a-deployment

<details>
  <summary>Deployment YAML File Sample</summary>

The following is an example deployment hello.yaml file:

```yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: hello
  labels:
    app: hello
spec:
  replicas: 1
  selector:
    matchLabels:
      app: hello
  template:
    metadata:
      labels:
        app: hello
    spec:
      containers:
      - name: hello
        image: pelasneacr.azurecr.io/hello:1.0.0
        ports:
        - containerPort: 80
```

</details>

&nbsp;

6. Deploy as a single container (1 replica) to Kubernetes.

Ref:
* https://docs.microsoft.com/en-us/azure/aks/tutorial-kubernetes-deploy-cluster
* https://v1-8.docs.kubernetes.io/docs/concepts/workloads/controllers/deployment/#creating-a-deployment

<details>
  <summary>Deployment Commands</summary>

```bash
# login to Kubernetes
az aks get-credentials --resource-group whatever-rg --name whatever
kubectl get nodes

# grant the Kubernetes service principal access to ACR
CLIENT_ID=$(az aks show --resource-group pelasne-aks --name pelasne-aks --query "servicePrincipalProfile.clientId" --output tsv)
ACR_ID=$(az acr show --resource-group pelasne-acr --name pelasneacr --query "id" --output tsv)
az role assignment create --assignee $CLIENT_ID --role Reader --scope $ACR_ID

# create the deployment
kubectl create -f hello.yaml --record --save-config
kubectl get deployments
kubectl rollout status deployment hello
kubectl get rs
kubectl get pods --show-labels
```
  
</details>

&nbsp;

## Expose the container to the internet

Ref:
* https://kubernetes.io/docs/concepts/services-networking/service
* https://kubernetes.io/docs/concepts/services-networking/service/#type-loadbalancer
* https://github.com/Azure-Samples/azure-voting-app-redis/blob/master/azure-vote-all-in-one-redis.yaml

<details>
  <summary>Deploy Commands</summary>
  
The following is an example hello-expose.yaml file:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: hello
spec:
  type: LoadBalancer
  ports:
  - port: 80
  selector:
    app: hello
```

Then you can run the following commands:

```bash
kubectl create -f hello-expose.yaml --record --save-config
kubectl get service hello --watch
```

Once you have an external IP it is done. You can then curl or open a browser to that IP and see your response.

</details>

&nbsp;

## Create a new version of the service

Create a new container that responds with "Hello Lab!" instead of "Hello World!".

Ref:
* https://docs.docker.com/engine/reference/commandline/tag
* https://docs.docker.com/engine/reference/commandline/push

<details>
  <summary>Create New Version</summary>

After changing the server.js source code, you can:

```bash
docker build -t hello:latest -t hello:2.0.0 -t whatever.azurecr.io/hello:latest -t whatever.azurecr.io/hello:2.0.0 .
docker images
docker push whatever.azurecr.io/hello:latest
docker push whatever.azurecr.io/hello:2.0.0
```

</details>

&nbsp;

## Deploy the new version in Kubernetes

Ref:
* https://v1-8.docs.kubernetes.io/docs/concepts/cluster-administration/manage-deployment/#in-place-updates-of-resources
* https://tachingchen.com/blog/kubernetes-rolling-update-with-deployment

<details>
  <summary>Deploy New Version</summary>
  
One way to do this is to modify the hello.yaml file to change the container image version to 2.0.0 and then:

```bash
kubectl apply -f hello.yaml --record
```

Another way would be to modify the existing deployment by:

```bash
kubectl edit deployment hello
```
  
</details>

&nbsp;

## Some other things to try

* Scale the number of replicas of your container
* Scale the number of Kubernetes nodes
* Deploy a database and setup connectivity between containers
