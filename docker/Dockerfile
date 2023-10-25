FROM golang:1.21-alpine3.17 AS build

RUN apk update

RUN apk --no-cache add gcc g++ make ca-certificates
RUN apk update -qq && apk add git

RUN mkdir -p /go/src/github.com/ramin/waypoint

ADD . /go/src/github.com/ramin/waypoint
WORKDIR /go/src/github.com/ramin/waypoint

# build the various artifacts
RUN go build -o ./bin/waypoint ./cmd/waypoint

# Build production image
FROM alpine:3.17

RUN apk update
RUN apk add ca-certificates
RUN update-ca-certificates

RUN mkdir -p /var/app

COPY --from=build /go/src/github.com/ramin/waypoint/bin/waypoint /var/bin/waypoint

CMD ["/var/bin/waypoint"]
