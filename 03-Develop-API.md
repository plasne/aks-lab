# Develop an API

In section 01, we developed a couple of Entity Services that allow us to do data operations on "songs" and "contracts". However, often we will not expose entity services directly, but rather put them behind an API. The API will be exposed outside of our AKS environment, but the individual entity services will not be.

There are a number of reasons why you might use an API service instead of exposing every service in your environment, including...

- This gives you a single place to do common tasks like authentication, authorization, rate limiting, etc. It is generally also possible to use an API management gateway to do these tasks.

- An API can be a view of your application taylored to a specific consumer or role, for example, you might have an API for your mobile and web applications, a different API for bulk operations, and a different API for administrators. You could change operations, protocols, models, etc. as needed to best optimize for the intended audience.

- Individual entity and process

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