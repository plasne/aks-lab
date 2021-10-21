FROM golang:1.17-alpine as build
WORKDIR /build
COPY ./go.mod .
COPY ./go.sum .
COPY ./*.go .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o songs .
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

FROM scratch as run
WORKDIR /app
COPY --from=build /build/songs .
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 80
CMD [ "./songs" ]