FROM golang
ADD . /go/src/github.com/liwei2001/go-grpc/server
WORKDIR /go/src/github.com/liwei2001/go-grpc/server
RUN go install github.com/liwei2001/go-grpc/server
ENTRYPOINT ["/go/bin/server"]
EXPOSE 5300
