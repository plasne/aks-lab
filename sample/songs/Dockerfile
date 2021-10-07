FROM golang:1.17-alpine as build
WORKDIR /build
COPY ./go.mod .
COPY ./go.sum .
COPY ./*.go .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o songs .

FROM scratch as run
WORKDIR /app
COPY --from=build /build/songs .
EXPOSE 80
CMD [ "./songs" ]