FROM golang:latest 

RUN mkdir -p /go/src/github.com/theopticians/optician-api 
ADD . /go/src/github.com/theopticians/optician-api/
WORKDIR /go/src/github.com/theopticians/optician-api 
RUN go get ./...
RUN go install
CMD ["/go/bin/optician-api"]
