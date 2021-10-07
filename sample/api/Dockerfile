FROM golang:1.17-alpine as build
WORKDIR /build
COPY ./go.mod .
COPY ./go.sum .
COPY ./*.go .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o api .

FROM scratch as run
WORKDIR /app
COPY --from=build /build/api .
EXPOSE 80
CMD [ "./api" ]