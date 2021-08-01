#!/bin/bash
echo What is proto file name\?
read path
exec protoc $path/$path"_pb"/$path.proto --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative