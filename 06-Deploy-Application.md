# Deploy Application

Perform the following steps to deploy your application:

## Push the containers to ACR.

Ref:

* https://docs.microsoft.com/en-us/azure/aks/tutorial-kubernetes-prepare-acr?tabs=azure-cli#create-an-azure-container-registry
* https://docs.docker.com/engine/reference/commandline/tag
* https://docs.docker.com/engine/reference/commandline/push

<details>
  <summary>Push Commands</summary>

```bash
# Set variables (adjust as needed)
ACR_NAME=akslab

az login
az acr login --name $

# For each of the images you build in the `Deploy to Docker` section
# Tag it so that you can push it to your ACR 
docker tag songs:1.0.0 $ACR_NAME.azurecr.io/songs:1.0.0

# Push it to your ACR
docker push $ACR_NAME.azurecr.io/songs:1.0.0

# Check ACR repositories
az acr repository list --name $ACR_NAME --output table

# Check the tags in your repository
az acr repository show-tags --name $ACR_NAME --repository songs --output table
```

</details>

&nbsp;

## Create the YAML files (manifests)

Create a deploy directory for all manifest files. A Kubernetes manifest file allows you to describe your workloads in the YAML format declaratively and simplify Kubernetes object management.

Build a YAML file for the deployment of each service: songs, contracts, api.

Ref:
* https://kubernetes.io/docs/concepts/workloads/controllers/deployment/

<details>
  <summary>Deployment YAML File Sample</summary>

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: songs-app
  labels:
    app: songs
spec:
  replicas: 1
  selector:
    matchLabels:
      app: songs
  template:
    metadata:
      labels:
        app: songs
    spec:
      containers:
      - name: songs
        image: pelasneakslabacr.azurecr.io/songs:1.0.0
        ports:
        - containerPort: 80
```

</details>

&nbsp;

A Kubernetes service is an abstract way to expose an application running on a set of Pods as a network service.
A Kubernetes service acts as a load balancer and redirects traffic to the specific ports of specified ports by using port-forwarding rules.
Build a YAML file for each service (songs, contracts, api). The api service is the only service that needs to be publicly exposed, so not all services will be of the same type.

Ref:

* https://kubernetes.io/docs/concepts/services-networking/service
* https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types
* https://docs.microsoft.com/en-us/azure/aks/tutorial-kubernetes-deploy-cluster

<details>
  <summary>Service YAML File Sample</summary>

```yaml
apiVersion: v1
kind: Service
metadata:
  name: songs
spec:
  type: ClusterIP
  ports:
  - port: 80
  selector:
    app: songs
```

</details>

&nbsp;

## Deploy to your cluster

Run the following commands for each of your YAML files:

```bash
kubectl apply -f myfile.yaml
```

Check your deployments and services:

```bash
kubectl get pods 
kubectl get services
```

Once all pods are in running state and you have an external IP for your api service, you can then curl or open a browser to that IP and see your response.

```bash
curl http://MYEXTERNAL-IP/song?id=6
```

</details>

&nbsp;

## Some other things to try

* Scale the number of replicas of your container
* Scale the number of Kubernetes nodes
