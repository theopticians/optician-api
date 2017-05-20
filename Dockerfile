#build stage
FROM golang:latest AS build

WORKDIR /go/src/github.com/theopticians/optician-api 
ADD . .
RUN go get ./...
RUN go install


# final stage
FROM debian:jessie-slim 
WORKDIR /app/
COPY --from=build /go/bin/optician-api .
CMD ["./optician-api"]  
