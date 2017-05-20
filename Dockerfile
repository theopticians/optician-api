#build stage
FROM golang:latest AS build-env

RUN mkdir -p /go/src/github.com/theopticians/optician-api 
ADD . /go/src/github.com/theopticians/optician-api/
WORKDIR /go/src/github.com/theopticians/optician-api 
RUN go get ./...
RUN go install
CMD ["/go/bin/optician-api"]


# final stage
FROM alpine

WORKDIR /app
COPY --from=build-env /go/bin/optician-api /app/
ENTRYPOINT ./optician-api
