
## 代码结构

```bash
.
├── bin                         # 存放编译好的 Go 二进制文件
├── src                         # 源码
│   ├── api                     # 后端 Go API 服务源码
│   │   ├── bootstrap           # 启动 Go API 服务相关参数包
│   │   ├── executor            # task 队列执行器
│   │   ├── handler             # HTTP handlers
│   │   ├── log                 # 基于 zap 封装的 log 包
│   │   ├── middleware          # HTTP 中间件
│   │   ├── model               # 全局通用 model
│   │   ├── router              # HTTP 路由
│   │   ├── service             # 封装好的服务
│   │   ├── db                  # 数据库
│   │   ├── util                # 工具类包
│   │   ├── version             # 提供运行时的版本信息等显示的支持
│   │   ├── worker              # task worker
│   │   └── main.go             # Go API 入口
│   └── client                  # 前端源码
│       ├── build               # Webpack 打包脚本
│       ├── src                 # 前端 js 源码
│       ├── package.json
│       ├── package-lock.json
│       └── README.md
├── tool                        # Makefile 可能会用到的一些编译脚本
│   ├── apppkg.sh
│   ├── build.sh
│   └── version.sh
├── vendor                      # Go 依赖
├── Gopkg.lock                  # dep 版本锁定文件，由 dep 生成
├── Gopkg.toml                  # dep 版本约束文件，用户可编辑
├── LICENSE
├── Makefile                    # Makefile文件
├── README_CN.md
├── README.md
└── run                         # 本地快速启动脚本
```
