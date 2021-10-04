# Develop an Entity Service

In a microservices architecture it is desirable to have a service for each entity that can be responsible for the general data operations related to that service (storing, retrieving, searching, etc.). Services that depend on this entity should ask the entity service to handle all data operations. This provides a number of benefits, including:

- You can chose the best solution for storing this specific data. Some data is always stored and retrieved by ID, a simple blob solution can be very cost effective and performant. Some data requires complex searches to find the right records, a database might be more appropriate. Other characteristics, like the size of the data, how the information will be projected, etc. might drive you towards specific storage solutions. Each entity service might make a different decision based on the requirements.

- Schemas may change from version to version. Having an entity service that is tightly coupled with the data storage solution allows this one service to make a schema change without affecting other services.

- Compartmentalizing the implementation details of an entity's storage to this service allows for an easy separation of concern. For instance, developers that need to use the entity don't need to know anything about how the data is stored or retrieved; they can simply operate on the public contract of this service.

## Data Contract

