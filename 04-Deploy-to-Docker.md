# Deploy to Docker

In this section, you will test your complete solution locally in Docker. Doing it locally before you deploy to Kubernetes is often easier for troubleshooting.

Ref:

* https://docs.docker.com/engine/reference/commandline/docker/
* https://docs.docker.com/language/golang/build-images/
* https://docs.docker.com/compose/networking/
* https://docs.docker.com/engine/reference/builder/#environment-replacement

## Create Dockerfiles

You should create a Dockerfile for each of your services (api, songs, and contracts). In your Dockerfile separate build from deployment.
Your Golang code will compile down to a single executable. Because the executable does not need a runtime, you can often get away without an operating system, and base off scratch.

<details>
  <summary>Dockerfile Sample</summary>

```Dockerfile
FROM golang:1.17-alpine as build
WORKDIR /build
COPY ./go.mod .
COPY ./go.sum .
COPY ./*.go .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o songs .

FROM scratch as run
WORKDIR /app
COPY --from=build /build/songs .
EXPOSE 80
CMD [ "./songs" ]
```

</details>

&nbsp;

## Build the container images

<details>
  <summary>Build Commands</summary>

You can build and view the built images by:

```bash
cd songs
docker build -t songs:1.0.0 .
docker images
```

You can test locally by:

```bash
docker run -d --name songs -p 9100:80 songs:1.0.0
curl http://localhost:9100/?id=7
```

</details>

&nbsp;

## Compose and Test

Compose is a tool for defining and running multi-container Docker applications. With Compose, you use a YAML file to configure your application’s services. Then, with a single command, you create and start all the services from your configuration.

Because the api will be using the songs and contracts entity services, a network overlay is required.
Compose also sets up a single network for your app by default. Create a `docker-compose.yaml` file at the root of your project and ensure the containers for the entity services are named appropriately so that they can be accessed by the api, and that the container and host ports are set correctly.

<details>
  <summary>Compose command</summary>

```bash
docker-compose up
```

</details>

&nbsp;

>Note: Environment variables could be used in the Dockerfiles but this is not required for this exercise.