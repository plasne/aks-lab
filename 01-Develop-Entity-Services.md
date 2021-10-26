# Develop Entity Services

In a microservices architecture it is desirable to have a service for each entity that can be responsible for the general data operations related to that service (storing, retrieving, searching, etc.). Services that depend on this entity should ask the entity service to handle all data operations. This provides a number of benefits, including:

- You can chose the best solution for storing this specific data. Some data is always stored and retrieved by ID, a simple blob solution can be very cost effective and performant. Some data requires complex searches to find the right records, a database might be more appropriate. Other characteristics, like the size of the data, how the information will be projected, etc. might drive you towards specific storage solutions. Each entity service might make a different decision based on the requirements.

- Schemas may change from version to version. Having an entity service that is tightly coupled with the data storage solution allows this one service to make a schema change without affecting other services.

- Compartmentalizing the implementation details of an entity's storage to this service allows for an easy separation of concern. For instance, developers that need to use the entity don't need to know anything about how the data is stored or retrieved; they can simply operate on the public contract of this service.

## Songs

You will develop a Songs Entity Service that will allow you to store (in-memory) and retrieve songs by ID. It should adhere to the following contract...

```yml
swagger: '2.0'
info:
  description: Song Entity Service
  version: 1.0.0
  title: Song Entity Service
paths:
  /:
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
            $ref: '#/definitions/song'
          x-examples:
            application/json: |-
              {
                  "id": 7,
                  "artist": "Eminem",
                  "title": "The Ringer",
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
```

You can use the following data in your code...

```go
type song struct {
	Id     int    `json:"id"`
	Artist string `json:"artist"`
	Title  string `json:"title"`
	Genre  string `json:"genre"`
}

var songs = []song{
	{0, "Drake", "In My Feelings", "HipHop"},
	{1, "Maroon 5", "Girls Like You", "Pop"},
	{2, "Cardi B", "I Like It", "HipHop"},
	{3, "6ix9ine", "FEFE", "Pop"},
	{4, "Post Malone", "Better Now", "Rap"},
	{5, "Eminem", "Lucky You", "Rap"},
	{6, "Juice WRLD", "Lucid Dreams", "Rap"},
	{7, "Eminem", "The Ringer", "Rap"},
	{8, "Travis Scott", "Sicko Mode", "HipHop"},
	{9, "Tyga", "Taste", "HipHop"},
	{10, "Khalid & Normani", "Love Lies", "HipHop"},
	{11, "5 Seconds Of Summer", "Youngblood", "Pop"},
	{12, "Ella Mai", "Boo'd Up", "HipHop"},
	{13, "Ariana Grande", "God Is A Woman", "Pop"},
	{14, "Imagine Dragons", "Natural", "Rock"},
	{15, "Ed Sheeran", "Perfect", "Pop"},
	{16, "Taylor Swift", "Delicate", "Pop"},
	{17, "Florida Georgia Line", "Simple", "Country"},
	{18, "Luke Bryan", "Sunrise, Sunburn, Sunset", "Country"},
	{19, "Jason Aldean", "Drowns The Whiskey", "Country"},
	{20, "Childish Gambino", "Feels Like Summer", "HipHop"},
	{21, "Weezer", "Africa", "Rock"},
	{22, "Panic! At The Disco", "High Hopes", "Rock"},
	{23, "Eric Church", "Desperate Man", "Country"},
	{24, "Nicki Minaj", "Barbie Dreams", "Rap"},
}
```

For example, you could run the following curl and get the following response...

```bash
> curl -i "http://localhost:9100/?id=5"
HTTP/1.1 200 OK
Content-Type: application/json
Date: Tue, 05 Oct 2021 14:06:50 GMT
Content-Length: 60

{"id":5,"artist":"Eminem","title":"Lucky You","genre":"Rap"}%
```

## Contracts

You will also develop a Contracts Entity Service that will allow you to retrieve payment terms by artist. It should adhere to the following contract...

```yml
swagger: '2.0'
info:
  description: Contracts Entity Service
  version: 1.0.0
  title: Contracts Entity Service
paths:
  /:
    get:
      produces:
        - application/json
      parameters:
        - name: artist
          in: query
          required: false
          type: string
          x-example: Taylor Swift
      responses:
        '200':
          description: Definition generated from Swagger Inspector
          schema:
            $ref: '#/definitions/contract'
definitions:
  contract:
    properties:
      artist:
        type: string
      payment:
        type: number
        format: double
```

You should default the payment to any artist that doesn't have a contract to 0.05. You can build some records that specify other contracts, such as...

```go
type contract struct {
	Artist  string  `json:"artist"`
	Payment float64 `json:"payment"`
}

var contracts = []contract{
	{"Drake", 0.2},
	{"Taylor Swift", 0.25},
	{"Khalid & Normani", 0.1},
}
```

For example, you could run the following curl and get the following response...

```bash
> curl -i "http://localhost:9200/?artist=Taylor%20Swift"
HTTP/1.1 200 OK
Content-Type: application/json
Date: Tue, 05 Oct 2021 14:14:47 GMT
Content-Length: 40

{"artist":"Taylor Swift","payment":0.25}
```

In reality, you would probably want to be able to store new contracts and adjust existing contracts with this entity service, but we are not going to do that as part of this lab.

## Debugging

To complete this activity, make sure your VSCODE development environment is configured to support debugging.

- https://code.visualstudio.com/docs/languages/go

## Tips

- When running locally, you will probably want the services running on different ports. When running in production, you will probably want the services all running on port 80. Environment Variables are a good way to do this.

- When making changes to a variable that could be accessed concurrently (for example, multiple incoming HTTP requests), you should use one of the [mutex](https://gobyexample.com/mutexes) options. There is an example of using a RWMutex in the sample.