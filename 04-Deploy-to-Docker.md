# Deploy to Docker

In this section, you will test your complete solution locally in Docker. Doing it locally before you deploy to Kubernetes is often easier for troubleshooting.

Ref:

* https://docs.docker.com/engine/reference/commandline/docker/
* https://docs.docker.com/language/golang/build-images/
* https://docs.docker.com/compose/networking/
* https://docs.docker.com/engine/reference/builder/#environment-replacement

## Create Dockerfiles

You should create a Dockerfile for each of your services (api, songs, and contracts). Your Golang code will compile down to a single executable.

## Compose and Test

Because the api will be using the songs and contracts entity services, a network overlay is required. Docker-compose sets up a single network for your app by default. Create a `docker-compose.yaml` file at the root of your project and ensure the containers for the entity services are named appropriately so that they can be accessed by the api, and that the container and host ports are set correctly.

>Note: Environment variables could be used in the Dockerfiles but this is not required for this exercise.
