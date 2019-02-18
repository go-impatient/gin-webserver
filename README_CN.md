
## 代码结构

```bash
.
├── bin                             # 存放编译好的 Go 二进制文件
├── src                             # 源码
│   ├── api                         # 后端 Go API 服务源码
│   │   ├── config                  # 配置和配置文件
│   │   ├── bootstrap               # 启动 Go API 服务相关参数包
│   │   ├── executor                # task 队列执行器
│   │   ├── handler                 # HTTP handlers
│   │   ├── log                     # 基于 zap 封装的 log 包
│   │   ├── model                   # 全局通用 model
│   │   ├── router                  # HTTP 路由相关
│   │   │   ├── middleware          # HTTP 中间件存放位置
│   │   │   └── router.go
│   │   ├── service                 # 封装好的服务
│   │   ├── db                      # 数据库
│   │   ├── util                    # 工具类包
│   │   ├── pkg                     # 引用的包
│   │   ├── worker                  # task worker
│   │   ├── swagger                 # swagger
│   │   └── main.go                 # Go API 入口
│   └── client                      # 前端源码
│       ├── build                   # Webpack 打包脚本
│       ├── src                     # 前端 js 源码
│       ├── package.json
│       ├── package-lock.json
│       └── README.md
├── tool                            # Makefile 可能会用到的一些编译脚本
│   ├── apppkg.sh
│   ├── build.sh
│   └── version.sh
├── vendor                          # Go 依赖
├── Gopkg.lock                      # dep 版本锁定文件，由 dep 生成
├── Gopkg.toml                      # dep 版本约束文件，用户可编辑
├── LICENSE
├── Makefile                        # Makefile文件
├── README_CN.md
├── README.md
└── run                             # 本地快速启动脚本
```
### 普通编译运行
```bash
$ cd $GOPATH/src/github.com/moocss/go-webserver
$ gofmt -w .
$ go tool vet .
$ go build -v .
$ ./apiserver
```
### 源码交叉编译并安装
```bash
curl $GOPATH/src/github.com/moocss/go-webserver/src/install.sh |bash
```

### Go 依赖
#### dep
```bash
# 安装 dep
go get -u github.com/golang/dep
dep init
dep ensure -add github.com/foo/bar github.com/baz/quux
dep ensure -v # 安装 Go 依赖
```

#### govendor
```bash
# Setup your project.
cd "my project in GOPATH"
govendor init

# Add existing GOPATH files to vendor.
govendor add +external

# View your work.
govendor list

# Look at what is using a package
govendor list -v fmt

# Specify a specific version or revision to fetch
govendor fetch golang.org/x/net/context@a4bbce9fcae005b22ae5443f6af064d80a6f5a55
govendor fetch golang.org/x/net/context@v1   # Get latest v1.*.* tag or branch.
govendor fetch golang.org/x/net/context@=v1  # Get the tag or branch named "v1".

# Update a package to latest, given any prior version constraint
govendor fetch golang.org/x/net/context

# Format your repository only
govendor fmt +local

# Build everything in your repository only
govendor install +local

# Test your repository only
govendor test +local
```