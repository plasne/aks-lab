# Deploy Application

Perform the following steps to deploy your application:

## Push the images to ACR

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

## Debugging

You can get a shell to a running container on a cluster's pod by using the `kubectl exec` command. In our case we could use that to curl the internal services (songs, contract) by deploying a client and using exec to get into it.

Ref:

https://kubernetes.io/docs/tasks/debug-application-cluster/get-shell-running-container/

<details>
  <summary>Client YAML File Sample</summary>

```yml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: client
  name: client
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: client
  template:
    metadata:
      labels:
        app: client
        version: v1
      name: client
    spec:
      containers:
      - name: client
        image: ubuntu
        command: ["/bin/bash", "-ec", "while :; do echo '.'; sleep 5 ; done; apt-get update && apt-get install -y curl"]
```

</details>

&nbsp;

```bash
# Deploy your client (adjust filename as needed)
kubectl apply -f client-manifest.yaml

# View all the pods that are running in your cluster
kubectl get pods

# Check the containers that are running on your client pod (notice that an istio sidecar container is running on every pod -more on that later :)
kubectl get pods [POD_NAME_HERE] -o jsonpath='{.spec.containers[*].name}'

# Using the Pod name use kubectl exec 
kubectl exec -it [POD_NAME_HERE] -c [CONTAINER_NAME_HERE] -- /bin/bash
```

## Some other things to try

* Scale the number of replicas of your container
* Scale the number of Kubernetes nodes
