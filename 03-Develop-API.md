# Develop an API

In the other sections, we developed entity services that allow us to do data operations on "songs" and "contracts". We may have also developed process services that allow us to add business value on top of our data. However, often we will not expose these services directly, but rather put them behind an API. The API will be exposed outside of our AKS environment, but the individual entity and process services will not be.

There are a number of reasons why you might use an API service instead of exposing every service in your environment, including...

- This gives you a single place to do common tasks like authentication, authorization, licensing, rate limiting, etc. It is generally also possible to use an API management gateway to do these tasks.

- An API can be a view of your application taylored to a specific consumer or role, for example, you might have an API for your mobile and web applications, a different API for bulk operations, and a different API for administrators. You could change operations, protocols, models, etc. as needed to best optimize for the intended audience.

- Often the useful activity that your API provides will involve contacting multiple services and mashing-up those results.

- Using a small number of API services in front of your entity and process services can improve security in the following ways...
  - You are presenting a smaller attack surface.
  - You are reducing unintended change - most of your code updates will happen in entity and process services.
  - The developers working on entity and process services don't need as much security insight (authentication, licensing, etc.).

- As mentioned previously, this design also allows for an easy separation of concerns. An API developer doesn't have to know anything about how or where data is stored, or anything about how computations are done, he/she simply needs to understand the services available to them and their contracts.

## API

You will develop an API that will allow you to store (in-memory) and retrieve songs and payment information by ID. It should adhere to the following contract...

```yml
swagger: '2.0'
info:
  description: API
  version: 1.0.0
  title: API
paths:
  /song:
    get:
      produces:
        - application/json
      parameters:
        - name: id
          in: query
          required: false
          type: string
          x-example: '7'
      responses:
        '200':
          description: Definition generated from Swagger Inspector
          schema:
            $ref: '#/definitions/songWithPayment'
          x-examples:
            application/json: |-
              {
                  "id": 7,
                  "artist": "Eminem",
                  "title": "The Ringer",
                  "payment": 0.05,
                  "genre": "Rap"
              }
    post:
      consumes:
        - application/json
      produces:
        - application/json
      parameters:
        - in: body
          name: body
          required: false
          schema:
            $ref: '#/definitions/song'
          x-examples:
            application/json: |-
              {
                  "title": "Boreas",
                  "artist": "Oh Hellos",
                  "genre": "Folk"
              }
      responses:
        '200':
          description: Definition generated from Swagger Inspector
          schema:
            $ref: '#/definitions/song'
          x-examples:
            application/json: |-
              {
                  "id": 25,
                  "title": "Boreas",
                  "artist": "Oh Hellos",
                  "genre": "Folk"
              }
definitions:
  song:
    required:
      - title
      - artist
      - genre
    properties:
      id:
        type: integer
      title:
        type: string
      artist:
        type: string
      genre:
        type: string
  songWithPayment:
    properties:
      id:
        type: integer
        format: int32
      title:
        type: string
      artist:
        type: string
      payment:
        type: number
        format: double
      genre:
        type: string
```

The GET method will return information about the song as well as how the artist is paid for streaming that song.

## Tips

- When running locally, you will probably want the services running on different ports. When running in production, you will probably want the services all running on port 80. Environment Variables are a good way to do this.

- When the API calls the entity services, those services could return certain HTTP error messages (for example, what if a song ID is not provided). Make sure your API returns those HTTP codes and messages when appropriate. The sample shows how to do this.