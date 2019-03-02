FROM golang
ADD . /go/src/scytale/organization/server
WORKDIR /go/src/scytale/organization/server
RUN go install scytale/organization/server
ENTRYPOINT ["/go/bin/server"]
EXPOSE 5300
