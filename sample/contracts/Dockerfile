# syntax=docker/dockerfile:1
FROM golang:1.16-alpine

WORKDIR /app

COPY ./contracts/go.mod ./
COPY ./contracts/go.sum ./
RUN go mod download

COPY ./contracts/*.go ./

RUN go build -o /contracts

EXPOSE 80

CMD [ "/contracts" ]