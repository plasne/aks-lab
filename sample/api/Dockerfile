# syntax=docker/dockerfile:1
FROM golang:1.16-alpine

WORKDIR /app

COPY ./api/go.mod ./
COPY ./api/go.sum ./
RUN go mod download

COPY ./api/*.go ./

RUN go build -o /api

EXPOSE 80

CMD [ "/api" ]