FROM golang:1.19-alpine3.16 as build-env
ENV appdir /app
ARG port=80


RUN mkdir -p /etc/ssl/certs/ca-certificates/ && update-ca-certificates --fresh && rm -rf /var/cache/apk/*
RUN update-ca-certificates


RUN mkdir -p $appdir
WORKDIR $appdir

ADD Makefile .
ADD ./scripts ./scripts

RUN echo http://dl-cdn.alpinelinux.org/alpine/edge/main >> /etc/apk/repositories
RUN echo http://dl-cdn.alpinelinux.org/alpine/edge/testing >> /etc/apk/repositories
RUN apk update
RUN apk add megatools

RUN apk --no-cache add gcc g++ make git bash
LABEL maintainer="Sample Org"

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .
# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/beAPIDev ./cmd/api/main.go
##RUN chmod +x script.sh

# Run the bin
FROM scratch

COPY --from=build-env /etc/ssl/certs/ca-certificates /etc/ssl/certs/
COPY --from=build-env /etc/ssl/certs/ /etc/ssl/certs/
COPY --from=build-env /usr/local/share/ca-certificates /usr/local/share/ca-certificates

WORKDIR /app
COPY --from=build-env /go/bin/beAPIDev /go/bin/beAPIDev
ENTRYPOINT ["/go/bin/beAPIDev"]

EXPOSE $port
