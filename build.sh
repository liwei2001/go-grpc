protoc -I ./pb --go_out=plugins=grpc:./organization ./organization/*.proto

docker build -t local/server -f Dockerfile_server .
