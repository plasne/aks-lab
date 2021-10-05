# Develop Process Services

In section 01, we developed some entity services. Those services allow us to do data operations on the entities in our solution, but they should not contain business logic for processing data in our solution. We will instead create some process services to act on the entities and outside data.

For this lab, we will not have any process services, but I have included it here for completeness. Some examples of the process services we might build include:

- A recommendation service that provides users with artists they may like based on following artists.

- A scheduled job to import songs as artists record new material.

- A scheduled job to ensure all data in the system is valid (no duplicates, all records have required information, etc.).