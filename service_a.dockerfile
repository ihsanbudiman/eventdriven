# service_a build stage
FROM golang:1.19.3-alpine3.16 AS service_a_builder

WORKDIR /go/bin/service_a

# Copy the go.mod and go.sum files
COPY go.mod .
COPY go.sum .

# get the dependencies
RUN go get -d -v ./...

COPY . .

# tidy
RUN go mod tidy

RUN go build -o service_a /go/bin/service_a/service_a

# service_a runner stage
FROM alpine:3.13
WORKDIR /go/bin/service_a

EXPOSE 8080
COPY --from=service_a_builder /go/bin/service_a/service_a .

# chmod
RUN chmod +x /go/bin/service_a/service_a

# chmod 755 to /go/bin/service_a/service_a
RUN chmod 755 /go/bin/service_a/service_a

ENTRYPOINT ["/go/bin/service_a/service_a"]
# CMD ["/go/bin/service_a/service_a"]
