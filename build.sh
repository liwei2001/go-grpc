protoc -I ./pb --go_out=plugins=grpc:./organization ./organization/*.proto

docker build -t local/server -f Dockerfile_server .
docker run -it -p 3000:3000 local/server

