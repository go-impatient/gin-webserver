// 交叉编译。
//
// 1. Mac 下编译 Linux 和 Windows 64位可执行程序
// CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o webserver
// CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -v -o webserver
//
// 2. Linux 下编译 Mac 和 Windows 64位可执行程序
// CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -v -o webserver
// CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -v -o webserver
//
// 3. Windows 下编译 Mac 和 Linux 64位可执行程序
// SET CGO_ENABLED=0
// SET GOOS=darwin
// SET GOARCH=amd64
// go build -v -o webserver
//
// SET CGO_ENABLED=0
// SET GOOS=linux
// SET GOARCH=amd64
// go build -v -o webserver

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

var (
	goos, goarch, goarm string
	race bool
)

func init() {
	flag.StringVar(&goos, "goos", "", "GOOS for which to build")
	flag.StringVar(&goarch, "goarch", "", "GOARCH for which to build")
	flag.StringVar(&goarm, "goarm", "", "GOARM for which to build")
	flag.BoolVar(&race, "race", false, "Enable race detector")
}

func main() {
	flag.Parse()

	gopath := os.Getenv("GOPATH")
	args := []string{
		"build",
		"-asmflags", fmt.Sprintf("-trimpath=%s", gopath),
		"-gcflags", fmt.Sprintf("-trimpath=%s", gopath),
	}
	if race {
		args = append(args, "-race")
	}

	env := os.Environ()
	env = append(env, "GOOS=" + goos, "GOARCH=" + goarch, "GOARM=" + goarm)
	if !race {
		env = append(env , "CGO_ENABLED=0")
	}

	cmd := exec.Command("go", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Env = env
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}