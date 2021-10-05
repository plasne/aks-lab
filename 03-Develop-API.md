# Develop an API

Many applications will require a REST API. The purpose 


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