FROM golang:alpine AS build

WORKDIR /go/src/github.com/mntor/littlebill

COPY . .

RUN go install /go/src/github.com/mntor/littlebill

FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=build /go/bin /usr/bin

EXPOSE 8080

CMD ["/usr/bin/littlebill"]