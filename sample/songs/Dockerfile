# syntax=docker/dockerfile:1
FROM golang:1.16-alpine

WORKDIR /app

COPY ./songs/go.mod ./
COPY ./songs/go.sum ./
RUN go mod download

COPY ./songs/*.go ./

RUN go build -o /songs

EXPOSE 80

CMD [ "/songs" ]