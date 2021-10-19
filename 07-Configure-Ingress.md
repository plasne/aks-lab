# Configure Ingress using Istio

In Kubernetes, an Ingress exposes routes for HTTP and HTTPS traffic from outside a cluster to services inside a cluster.
Ingress may provide load balancing, SSL termination and name-based virtual hosting.

Istio is a service mesh that allows you to transparently add capabilities like traffic management, observability, and security, to your cluster without adding them to your code.

Ref:

* https://kubernetes.io/docs/concepts/services-networking/ingress/
* https://istio.io/latest/about/service-mesh/

## Create an Istio Gateway and configure routes for traffic

Istio traffic management relies on the Envoy proxies that are deployed along with your services, and lets you easily control the flow of traffic and API calls between services.

Along with support for Kubernetes Ingress, Istio offers another configuration model, Istio Gateway. A Gateway provides more extensive customization and flexibility than Ingress, and allows Istio features such as route rules to be applied to traffic entering the cluster.

Ref:

* https://istio.io/latest/docs/concepts/traffic-management/
* https://istio.io/latest/docs/tasks/traffic-management/ingress/ingress-control/

<details>
  <summary>Create YAML for the Istio Gateway and VirtualService</summary>

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: song-gateway
spec:
  selector:
    istio: ingressgateway # use Istio default gateway implementation
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "*"
---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: song
spec:
  hosts:
  - "*"
  gateways:
  - default/song-gateway
  http:
  - match:
    - uri:
        exact: /song
    route:
    - destination:
        host: api
        port:
          number: 80
```

>NOTE: For the purpose of this lab, you can use a wildcard `*` value for the host in the Gateway and VirtualService configurations. In a real world scenario, you would use your host's domain name.
</details>
&nbsp;

Deploy your manifest YAML file.

```bash
kubect apply -f <manifest-file-name>.yaml
```

Now that Ingress is configured, update the type of your api service to ClusterIP (vs. LoadBalancer).

## Test your configuration

Get the internal IP for your istio-ingressgateway to test your API using curl or a browser.

```bash
# Get istio gateway internal IP
ISTIO_GATEWAY_INTERNALIP=$(kubectl get svc istio-ingressgateway -n istio-system -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
echo $ISTIO_GATEWAY_INTERNALIP

# exec into the client pod to test
kubectl exec -it <client-pod-name> -c client -- /bin/bash
curl -i http://<use your ISTIO_GATEWAY_INTERNALIP>/song?id=7
```

## Create an App Gateway instance

TODO:
- Add why AppGateway is needed (WAF, ssl termination, https) and how Public IP address is proxied to private IP (versus using an external load balancer)
- Add clarification on why health endpoint was needed - Retest and include commands to add health probe to app gateway

<details>
  <summary>Provision App Gateway</summary>

```bash
# Set variables
RESOURCE_GROUP=akslabhv-rg # created in previous chapter
VNET_NAME=akslabhv-vnet # created in previous chapter
APP_GATEWAY=akslabhv-gw
APP_GATEWAY_SKU=Standard_v2
PUBLIC_IP_ADDRESS=akslabhv-appgw-ip
GW_SUBNET_NAME=appgw-subnet

# Create subnet for the App Gateway
az network vnet subnet create -g $RESOURCE_GROUP --vnet-name $VNET_NAME -n $GW_SUBNET_NAME --address-prefixes 10.0.0.0/24

GW_SUBNET_ID=$(az network vnet subnet show --resource-group $RESOURCE_GROUP --vnet-name $VNET_NAME --name $GW_SUBNET_NAME --query id -o tsv)

# Create App Gateway
az network application-gateway create -g $RESOURCE_GROUP -n $APP_GATEWAY --sku $APP_GATEWAY_SKU --subnet $GW_SUBNET_ID --servers $ISTIO_GATEWAY_INTERNALIP --public-ip-address $PUBLIC_IP_ADDRESS
```

</details>
&nbsp;
