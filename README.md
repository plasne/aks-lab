# AKS Lab
A short lab for getting an application up and running using Azure Container Service (AKS).

It is not required to have any of these tools installed, but you could install them if you want to run some of this locally:
* Azure CLI 2.0: https://docs.microsoft.com/en-us/cli/azure/install-azure-cli?view=azure-cli-latest
* Node.js: https://nodejs.org/en/
* 

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

1. Write a Hello World Node.js application that will expose an HTTP endpoint on port 80.

<details>
  <summary>Node.js Hello World Code Sample</summary>

You can provision a new app and install Express by:

```bash
npm init
npm install express --save
```

The server.js file could look something like this:

```javascript
const express = require("express");
const app = express();

app.get("/", (req, res) => {
  res.send.("Hello World!\n");
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

<details>
  <p>
    
```Dockerfile
FROM node:latest
COPY server.js server.js
COPY package.json package.json
RUN npm install
ENV PORT 80
EXPOSE 80
CMD node server.js
```
    
  </p>
</details>

&nbsp;
3. Build the container image.

<details>
  <p>

You can build and view the built images by:

```bash
docker build -t hello:latest -t hello:1.0.0 .
docker images
```

You can test locally by:

```bash
docker run --name hello --publish 8800:80 hello:lastest
curl http://localhost:8800
```

  </p>
</details>

&nbsp;
4. Push the container to ACR.

<details>
  <p>

```bash
az acr login --name whatever.azurecr.io
docker tag hello whatever.azure.io/hello:latest whatever.azure.io/hello:1.0.0
docker push whatever.azure.io/hello:latest whatever.azure.io/hello:1.0.0
az acr repository list --name whatever --output table
```

  </p>
</details>

&nbsp;
Ref:
* https://expressjs.com/en/starter/hello-world.html
* https://docs.docker.com/engine/reference/builder
* https://hub.docker.com/_/node
* https://docs.docker.com/engine/reference/commandline/build
* https://docs.docker.com/engine/reference/run
* https://docs.microsoft.com/en-us/azure/aks/tutorial-kubernetes-prepare-acr

## Expose the container

## Add a Postgres database

## Add a content creation service

## Deploy a new version of the web front end

## Scale the content creation service
