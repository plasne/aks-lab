# AKS and Golang Lab

This lab will help you design a Golang application, deploy it in Azure Kubernetes Service, and will cover some traffic routing via Istio and Azure Application Gateway.

The application we will develop includes a few simple services to support a music streaming service. You will be able to get song metadata by ID from an API which will also include how much the artist is to be paid for the song.

Interwoven in this lab are some general design concepts for microservices applications that I have employed over the years. I tend to refer to service types as entity, process, and API. I found this article by someone else that describes an extremely similar design pattern with some different names: <https://www.linkedin.com/pulse/types-microservices-maxim-alexandrovich/>.

There is sample code in the sample folder, but please attempt to do these activies on your own before consulting them.

First, write your application and verify it works locally...

- [01: Develop Entity Services](./01-Develop-Entity-Services.md)
- [02: Develop Process Services](./02-Develop-Process-Services.md)
- [03: Develop API](./03-Develop-API.md)

Second, test the solution locally in Docker...

- [04: Deploy to Docker](./04-Deploy-to-Docker.md)

Third, deploy the Azure infrastructure to run your application...

- [05: Deploy Azure Resources](./05-Deploy-Azure-Resources.md)
- [06: Deploy Application](./06-Deploy-Application.md)
- [07: Configure Ingress](./07-Configure-Ingress.md)

Last, add a database for your service and deploy a second version...

- [08: Add a Database](./08-Add-a-Database.md)
- [09: Deploy Multiple Versions](./09-Deploy-Multiple-Versions.md)
