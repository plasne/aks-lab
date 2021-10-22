# Deploy Multiple Versions

In this section we will see how to deploy multiple versions of the songs service.

Imagine you have highly decoupled services as you often do with a microservices architecture pattern. One advantage to this strategy is a developer can write new functionality for a single service, deploying a new version of it easily because the service is small and easily understood. Updating this one building block is easy, but you may find there are dozens of other services that take a dependency on this service to operate. It is generally impractical to upgrade all those other services at the same time. Therefore, it becomes common for several versions of a service to be running in production at the same time. Dependency services can call the version they need until such time as all services have been upgraded to the point that the old service is no longer needed.

As discussed in the prior section, this strategy generally applies to major versions as minor versions and bug fixes should generally keep the existing data contracts and functionality intact. Minor versions can extend functionality provided that doesn't break existing workflows.

Feature flags or API version designations are often used in conjunction with this model to control when a consumer of the services uses the default workflow or an alternate workflow.

## Implementation

For this exercise, you will...

1. Write a new manifest for the new deployment strategy. This should probably include 2 versions of the songs service (v1 and v2), a Service, a VirtualService, and a DestinationRule.

    There are a number of ways to route traffic to a specific version, you could use different host names, paths, querystring parameters, or headers.
    
    For the sample ([sample/deploy/songs-versions-manifest.yaml](sample/deploy/songs-versions-manifest.yaml)), I chose to use a header called "x-api-version" which should either be "v1" or "v2". If a value outside of that range is specified, an HTTP 400 Bad Request is returned. If the header is not specified it goes to v1 - this maintains compatibility with the prior implementation. Note that the sample manifest is meant to replace the songs-manifest.yaml, so you will probably need to revoke that previous manifest if deployed.

2. Modify your API to route to the appropriate songs service based on the client specifying the intended version.

3. Test that you can hit v1 and that it returns the same response as in your prior implementation. Test that you can hit v2 and that it returns as expected. Test that you can hit v3 and that it returns a HTTP 400 Bad Request.

## Discussion

### Is there a better way to handle the fault in the sample?

In the sample, if an x-api-version header is supplied and is not in the supported range, then we route it to "service-that-can-return-proper-error-message", except that this service doesn't exist and there is a fault applied 100% of the time to return HTTP 400 Bad Request.

```yaml
  - match:
    - headers:
        x-api-version:
          regex: ".*"
    fault:
      abort:
        percentage:
          value: 100
        httpStatus: 400
    route:
      - destination:
          host: service-that-can-return-proper-error-message
```

This does provide us with the correct status code, but it doesn't supply us with a helpful error message (body). In a real implementation, I would probably deploy an NGINX container that hosts a bunch of error responses and simply `redirect` to the appropriate one (without the fault handling). An alternative is to use Envoy Filters, but NGINX tends to be more flexible in this regard and is probably easier to deploy and maintain.