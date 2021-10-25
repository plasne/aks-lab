# Add a Database

In this section we will add an external database to the songs entity service to overcome a few limitations of an in-memory database.

## The Problem

In section 05-Deploy-Azure-Resources, if you followed the sample, you would have noticed that only 1 replica of each service was deployed. That doesn't provide much resiliency, so what we could scale those...

```bash
kubectl scale deployment api-app --replicas=3
kubectl scale deployment contracts-app --replicas=3
kubectl scale deployment songs-app --replicas=3
```

The API service is stateless so there is no problem with scaling it. While unrealistic, our contracts service also cannot be modified (the contracts are hard-coded), so it can also be scaled. However, if you scaled the songs service and posted a new song to it...

```bash
> curl -i -H "x-api-version: v1" -X POST -d '{ "artist": "The Oh Hellos", "title": "Bitter Water", "genre": "Folk" }' "http://52.150.27.198/song"
```

...and then tried to retrieve it, you would sometimes be successful and sometimes not...

```bash
> curl -i -H "x-api-version: v1" "http://52.150.27.198/song?id=25"
{"artist":"The Oh Hellos","genre":"Folk","id":25,"payment":0.05,"title":"Bitter Water"}

> curl -i -H "x-api-version: v1" "http://52.150.27.198/song?id=25"
HTTP/1.1 400 Bad Request
Date: Thu, 21 Oct 2021 18:55:47 GMT
Content-Type: text/plain; charset=utf-8
Content-Length: 26
Connection: keep-alive
x-content-type-options: nosniff
x-envoy-upstream-service-time: 2
server: istio-envoy

the ID was out-of-range.
```

The reason should hopefully be obvious - there is no external database, replication, or shared volume for the 3 replicas. A change to one replica does not change the state of another. This problem becomes even worse as you add more new songs because you will start duplicating IDs - in short order you will have replicas that are completely out-of-sync.

One other problem with this design is that if a replicas is destroyed for any reason, the changes are lost.

## Solutions

As mentioned, there are a number of solutions to this problem:

- External Database - in Azure, on-premises, in the AKS cluster, etc.
- Replication - you could make changes to a single primary and replicate those changes to the secondary (read-only) replicas.
- Shared Volume - you could share a volume between the replicas and have them all read/write to the same file(s).

The case you are most likely to encounter is using an external database which is what we will use in this chapter.

## Create an Azure Cosmos instance using the Mongo API

### Why use the Mongo API?

When this lab was written, there is no SDK for using the Cosmos Core API wire protocol in Golang.
We could use the REST API, Mongo API, or Cassandra API. I chose the Mongo API because it is the easiest to use.

Ref:

- https://docs.microsoft.com/en-us/azure/cosmos-db/mongodb/mongodb-introduction
- https://docs.microsoft.com/en-us/cli/azure/cosmosdb?view=azure-cli-latest#az_cosmosdb_create

<details>
  <summary>Provision an Azure Cosmos instance with the Mongo API</summary>

```bash
# set your variables
RESOURCE_GROUP=akslabhv-rg
DB_NAME=akslabhv-mongodb

# create the cosmos account
az cosmosdb create --name $DB_NAME --resource-group $RESOURCE_GROUP --kind MongoDB --server-version 4.0

# create the database
az cosmosdb mongodb database create --account-name $DB_NAME --name db --resource-group $RESOURCE_GROUP

# create the collection
az cosmosdb mongodb collection create --account-name $DB_NAME --database-name db --resource-group $RESOURCE_GROUP --name col

# get the connection string options
az cosmosdb keys list --type connection-strings --name $DB_NAME --resource-group $RESOURCE_GROUP
```

</details>

&nbsp;

## Create a v2 version of the songs service to read and write songs to Cosmos

Create a v2 version of the songs service to read and write songs to Cosmos.

Build a new Dockerfile and test it locally as part of your solution.

Do not deploy your service to AKS yet, we will do that in section 09.

Ref:

- https://golang.org/doc/modules/major-version
- https://pkg.go.dev/go.mongodb.org/mongo-driver#readme-usage
- https://pkg.go.dev/context
- https://docs.microsoft.com/en-us/azure/cosmos-db/mongodb/create-mongodb-go
- https://docs.mongodb.com/drivers/go/current/fundamentals/bson/

### Tips for updating the song entity service

<details>
  <summary>Song IDs</summary>

When our song database was in-memory on a single replica it was easy enough to use a mutex to ensure that we could increment our IDs. However, changing to a distributed system makes that more complex. There are multiple ways to address this:

- We could use a system that handles external mutexes for distributed systems (ex. Zookeeper).
- We could use a trigger in Mongo to build an incrementing ID.
- We could let Mongo generate a unique alphanumeric ID on write.

In the sample, I chose the last option as this is the easiest to implement and does not require a blocking operation which would reduce throughput (the other 2 options suffer from this issue).

</details>

&nbsp;

## Update your API so that both versions of the songs service are supported

### Tips for updating the API

<details>
  <summary>Versioning the API</summary>

You may notice from the sample, that I do NOT create a major version of the API, like I suggested for the song service. Why?

The song service needs to be a new major version because:

- The data contract is changing - the ID will be a string now instead of an int.
- The behaviors are changing - changes are persisted, multiple replicas are now supported.

However, the API service can continue to operate in the exact same way as it did before (provided we leave a v1 version of the song service up in our cluster - more on this in the next section). It continues to be stateless and continues to return the data from the songs and contracts services.

</details>

&nbsp;

<details>
  <summary>Varying Schemas</summary>

Once we have 2 versions of our songs service with 2 different schemas (the ID is an integer in v1 and a string in v2), we cannot deserialize to a single "song" struct. There are several ways we could solve this problem:

- We could have 2 separate song structs and determine the one to use based on the version of the service we called.
- We could simply deserialize to `map[string]interface{}`, in which case we don't care what the underlying service returns.

I prefer the 2nd approach because our API is not an authority on what should be in the song schema. If we wanted to create a v3 songs service that has 10 extra fields, so be it - the API could simply ignore all those fields and return them being none the wiser.

Having said that, there are legitimate cases to be made for both options.

</details>

&nbsp;

## Test your solution

<details>
  <summary>Routing in Docker</summary>

When we deploy to AKS and use Istio, we will have the capability of routing to specific service versions by using the host name, path, querystring, headers, etc. Since we do not have these advanced routing capabilities in native Docker, you can...

- Simply use docker-compose to test your solution with v1, then bring it up again and test with v2.
- Give different names to your service versions (ex. songs-v1 and songs-v2).

</details>
