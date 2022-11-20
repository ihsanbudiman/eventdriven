# service_b build stage
FROM golang:1.19.3-alpine3.16 AS service_b_builder
WORKDIR /go/bin/service_b

# Copy the go.mod and go.sum files
COPY go.mod .
COPY go.sum .

# get the dependencies
RUN go get -d -v ./...

COPY . .

# tidy
RUN go mod tidy

RUN go build -o service_b /go/bin/service_b/service_b

# service_b runner stage
FROM alpine:3.13
WORKDIR /go/bin/service_b

COPY --from=service_b_builder /go/bin/service_b/service_b .

# chmod
RUN chmod +x /go/bin/service_b/service_b

# chmod 755 to /go/bin/service_b/service_b
RUN chmod 755 /go/bin/service_b/service_b

ENTRYPOINT ["/go/bin/service_b/service_b"]
# CMD ["/go/bin/service_b/service_b"]