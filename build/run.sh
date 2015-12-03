#!/bin/bash

#### Config Vars ####
# update these to reflect your service
ServiceName="base"
BasePath="/home/core/dev/"
port="8080"


# determine service path for volume mounting
CurrentDir=`pwd`
ServicePath="${CurrentDir/$BasePath/}"

# run the build
docker run -it -v `pwd`:"/go/src/$ServicePath" divideandconquer/godep:1.5.1 /bin/bash -c "cd /go/src/$ServicePath; ./build/build.sh"

# build the docker container with the new binary
docker build -t $ServiceName .

# run the container
docker run -it -p $port:8080 $ServiceName
