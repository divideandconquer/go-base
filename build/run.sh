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
docker run -it -v `pwd`:"/go/src/$ServicePath" divideandconquer/godep:1.5.1 /bin/bash -c "cd /go/src/$ServicePath; ./build/build.sh" || { echo 'build failed' ; exit 1; }

# update the config 
ConsulAddr=`ifconfig eth1 | grep "inet " | sed -e 's/^[[:space:]]*//' | cut -d" " -f2`
ConfigPath="`pwd`/config/dev.json"
docker run -v $ConfigPath:/config.json divideandconquer/go-consul-client  -file /config.json -namespace dev/$ServiceName -consul $ConsulAddr:8500



# build the docker container with the new binary
docker build -t $ServiceName .

# run the container
docker run -it -p $port:8080 $ServiceName
