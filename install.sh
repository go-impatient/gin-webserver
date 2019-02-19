#!/bin/bash

server_repo="github.com/moocss/go-webserver/src"
build_repo="github.com/moocss/go-webserver/src/cmd"

install_path=$( cd `dirname $0`; pwd )/server-deploy

if [ -d ${install_path} ];then
    install_path=${install_path}-$( date +%Y%m%d%H%M%S )
fi

if [ -z ${GOPATH} ];then
    GOPATH=`go env GOPATH`
fi

build_server() {
    # go get ${build_repo}
    cd $GOPATH/src/${build_repo}
    go run build.go
}

install_server() {
    mkdir ${install_path}
    cd ${install_path}
    mkdir bin log etc
    cp $GOPATH/src/${build_repo}/webserver ./bin/
    cp $GOPATH/src/${server_repo}/config.yaml ./etc/config.yaml
    cp -r $GOPATH/src/${server_repo}/client ./client
}

build_web() {
    cd client && yarn install && yarn build
}

build_server

install_server

# build_web

echo "Installing server binary: ${install_path}/bin"
echo "Installing server public: ${install_path}/client"
echo "Installing config.yaml: ${install_path}/etc"
echo "Install complete"
