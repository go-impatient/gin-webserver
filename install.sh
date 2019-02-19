#!/bin/bash
build_numbar="0.1.0"
server_repo="github.com/moocss/go-webserver"
build_repo="github.com/moocss/go-webserver/cmd/webserver"
version_repo="github.com/moocss/go-webserver/src/pkg/version"

install_path=$( cd `dirname $0`; pwd )/server-deploy

if [ -d ${install_path} ];then
    install_path=${install_path}-$( date +%Y%m%d%H%M%S )
fi

if [ -z ${GOPATH} ];then
    GOPATH=`go env GOPATH`
fi

gotool() {
    # 格式化代码
	gofmt -w .
	# 代码检查并跳过vendor
	go tool vet . | grep -v vendor;true
}

build_all() {
    gotool
     # linux
    GOOS=linux  GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-w -X ${version_repo}.VersionDev=build.${build_numbar} -X ${versionDir}.VersionDate=$( date +%Y%m%d%H%M%S )" -v -a -installsuffix cgo -o release/webserver_linux
    # darwin:
    GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-w -X ${version_repo}.VersionDev=build.${build_numbar} -X ${versionDir}.VersionDate=$( date +%Y%m%d%H%M%S )" -v -a -installsuffix cgo -o release/webserver
    # windows
    GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-w -X ${version_repo}.VersionDev=build.${build_numbar} -X ${versionDir}.VersionDate=$( date +%Y%m%d%H%M%S )" -v -a -installsuffix cgo -o release/webserver.exe
}

build() {
    gotool
    go build -v -ldflags "-w -X ${version_repo}.VersionDev=build.${build_numbar} -X ${versionDir}.VersionDate=$( date +%Y%m%d%H%M%S )" -o release/webserver
}

dev_server() {
    cd $GOPATH/src/${build_repo}
    go run -dev
}

build_server() {
    # go get ${build_repo}
    cd $GOPATH/src/${build_repo}
    build_all
}   

install_server() {
    mkdir ${install_path}
    cd ${install_path}
    mkdir bin log etc
    cp $GOPATH/src/${build_repo}/release/webserver ./bin/
    cp $GOPATH/src/${build_repo}/release/webserver_linux ./bin/
    cp $GOPATH/src/${build_repo}/release/webserver.exe ./bin/
    cp $GOPATH/src/${server_repo}/src/config.yaml ./etc/config.yaml
    cp -r $GOPATH/src/${server_repo}/client ./client
}

build_web() {
    cd client && yarn install && yarn build
}

build_server

install_server

# build_web

echo "Installing server: ${install_path}/bin"
echo "Installing client: ${install_path}/client"
echo "Installing config: ${install_path}/etc"
echo "Install done..."
