#!/bin/bash

web_repo="github.com/moocss/go-webserver"
build_repo="github.com/moocss/go-webserver/cmd"

install_path=$( cd `dirname $0`; pwd )/web-deploy

if [ -d ${install_path} ];then
    install_path=${install_path}-$( date +%Y%m%d%H%M%S )
fi

if [ -z ${GOPATH} ];then
    GOPATH=`go env GOPATH`
fi

build_web() {
    go get ${build_repo}
    cd $GOPATH/src/${build_repo}/web
    go run build.go
}

install_web() {
    mkdir ${install_path}
    cd ${install_path}
    mkdir bin log etc
    cp $GOPATH/src/${build_repo}/webserver ./bin/
    cp $GOPATH/src/${web_repo}/conf.yaml ./etc/conf.yaml
    cp -r $GOPATH/src/${web_repo}/client ./client
}


build_web

install_web

echo "Installing web binary: ${install_path}/bin"
echo "Installing web public: ${install_path}/client"
echo "Installing conf.yaml: ${install_path}/etc"
echo "Install complete"
