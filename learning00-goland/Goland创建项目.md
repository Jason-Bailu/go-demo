[Goland创建项目](https://blog.csdn.net/2401_84911504/article/details/138903184?ops_request_misc=&request_id=&biz_id=102&utm_term=%E4%BD%BF%E7%94%A8goland%E5%88%9B%E5%BB%BA%E9%A1%B9%E7%9B%AE&utm_medium=distribute.pc_search_result.none-task-blog-2~all~sobaiduweb~default-0-138903184.142^v100^pc_search_result_base7&spm=1018.2226.3001.4187)

- **关于代理：** 留下以下几种资源地址,注意：代理修改完，重启 GoLand 生效

​	阿里云https://mirrors.aliyun.com/goproxy/
​	nexus 社区提供的https://gonexus.dev
​	goproxy.io 的https://goproxy.io/
​	官方提供的https://proxy.golang.org
​	七牛云赞助支持的https://goproxy.cn

- **是否使用 go mod 的区别：**

​	开启mod：go env -w GO111MODULE=on ，会将包下载到 gopath 下的 pkg 下的 mod 文件夹中
​	关闭mod：go env -w GO111MODULE=off ，会将包下载到 gopath 下的 src 下
​	go env GO111MODULE=auto 只有当前目录在 GOPATH/src 目录之外而且当前目录包含 go.mod 文件或者其子目录包含 go.mod文件才会启用

- **go mod 概念：**官方推荐 mod 是go内置的模块管理器，管理的依赖包的版本，能保证在不同地方构建，获得的依赖模块是一致的

  ```bash
  go mod init 生成go.mod文件
  go mod download 下载go.mod中指定的所有依赖
  go mod tidy 整理现在的依赖
  go mod graph 查看现有的依赖结构
  go mod edit 编辑go.mod文件
  go mod vendor 导出项目所有依赖到vendor目录
  go mod verify 校验一个模块是否被篡改过
  go mod why 查看为什么需要依赖某模块
  ```
  
- **go项目结构：**[go项目目录布局](https://github.com/golang-standards/project-layout/blob/master/README_zh.md) 注：不要在项目中创建src目录，让go项目像java一样

  - /cmd：本项目的主干，main函数的位置，每个应用程序的目录名应该与你想要的可执行文件的名称相匹配(例如，`/cmd/myapp`)
  - /internal：私有应用程序和库代码。这是你不希望其他人在其应用程序或库中导入代码。
  - /pkg：外部应用程序可以使用的库代码(例如 `/pkg/mypubliclib`)。其他项目会导入这些库，希望它们能正常工作。
  - /vendor：应用程序依赖项(手动管理或使用你喜欢的依赖项管理工具，如新的内置 [`Go Modules`](https://go.dev/wiki/Modules) 功能)。
  - /api：OpenAPI/Swagger 规范，JSON 模式文件，协议定义文件。服务器目录
  - /web：特定于 Web 应用程序的组件:静态 Web 资源、服务器端模板和 SPAs。Web应用目录
  - /configs：配置文件模板或默认配置。confd 和 consul-template
  - /init：System init（systemd，upstart，sysv）和 process manager/supervisor（runit，supervisor）配置。
  - /scripts：执行各种构建、安装、分析等操作的脚本。
  - /build：打包和持续集成。
  - /deployments：IaaS、PaaS、系统和容器编排部署配置和模板(docker-compose、kubernetes/helm、mesos、terraform、bosh)。
  - /test：额外的外部测试应用程序和测试数据。
  - /docs：设计和用户文档(除了 godoc 生成的文档之外)。
  - /tools：这个项目的支持工具。注意，这些工具可以从 `/pkg` 和 `/internal` 目录导入代码。
  - /examples：你的应用程序和/或公共库的示例。
  - /third_party：外部辅助工具，分叉代码和其他第三方工具(例如 Swagger UI)。
  - /githooks：Git hooks。
  - /assets：与存储库一起使用的其他资源(图像、徽标等)。
  - /website：如果你不使用 Github 页面，则在这里放置项目的网站数据。