
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

### Go 依赖

我们目前使用 [dep](https://github.com/golang/dep) 管理依赖。

```bash
# 安装 dep
go get -u github.com/golang/dep
dep init
dep ensure -add github.com/foo/bar github.com/baz/quux
dep ensure -v # 安装 Go 依赖
```