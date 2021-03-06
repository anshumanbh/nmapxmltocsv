# Build Container
FROM golang:1.9.4-alpine3.7 AS build-env
RUN apk add --no-cache --upgrade git openssh-client ca-certificates
RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/anshumanbh/nmapxmltocsv

# Cache the dependencies early
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only -v

COPY main.go ./

RUN go install

# Final Container
FROM alpine:3.7
LABEL maintainer="Anshuman Bhartiya"
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-env /go/bin/nmapxmltocsv /usr/bin/nmapxmltocsv

COPY in.xml /

ENTRYPOINT ["/usr/bin/nmapxmltocsv"]
