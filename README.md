#  Golang实现高性能微服务网关

##  前言：

> 这个项目官方推荐学习时间为30个小时，但是我应该上了50个小时左右
>
> 因为是第一个练手，到项目那一块几乎所有代码我都敲了一遍。（推荐）

慕课链接：https://coding.imooc.com/class/chapter/436.html

导师项目地址：https://github.com/e421083458

本仓库将在原有的基础上添加更多的注释。这是我学完Go的**第一个实战项目**（PHP转Go）

(**明白原理后可以自己添加后端功能，前端设计可以自己修改下，开发毕设是没问题的**)

虽然`Nginx`是不二之选，但是用Go学习实现不仅复习**计网**的部分知识，还学会了大量**技术**。

这个网关是完全能**日常使用**的，并且在匹配服务的是否**加上自己的逻辑**就能实现各种效果。

##  技术选择

|     描述     |      技术栈       |                    参考地址                     |
| :----------: | :---------------: | :---------------------------------------------: |
|     前端     | Vue-Element-Admin | https://github.com/PanJiaChen/vue-element-admin |
|     图表     |      Echarts      |    https://echarts.apache.org/zh/index.html     |
|     后端     |        Gin        |        https://github.com/gin-gonic/gin         |
|    数据库    |   Mysql + Reids   |           https://www.google.com.hk/            |
|   开发文档   |      Swagger      |     https://github.com/swaggo/swag/releases     |
|   公共类库   |   Golang_Common   |   https://github.com/e421083458/golang_common   |
|   用户管理   |     OAuth2.0      |              https://oauth.net/2/               |
|   服务发现   |     Zookeeper     |      https://alsritter.icu/posts/20bb8062/      |
| 熔断降级限流 |      Hystrix      |       https://github.com/Netflix/Hystrix        |
|   压力测试   |   Apache Bench    |            https://httpd.apache.org/            |

##  文件结构

> 主要说一下后端结构，详细的看课与导师的仓库
>
> 后端是把`http`与`tcp`与`grpc`三项服务与`管理后端`写在一起的

```
├── certFile # SSL配置文件
├── conf # 全局配置文件
├── controller # 控制器
├── dao # 模型
├── docs # 文档
├── dto # 输入输出
├── grpcProxyMiddleware # grpc中间件
├── grpcProxyRouter # grpc路由
├── httpProxyMiddleware # http中间件
├── httpProxyRouter # http路由
├── logs # 日志
├── middleware # 基础中间件
├── public # 公共类库
├── reverseProxy # 代理插件
├── router # 主路由
├── services # 服务
├── tcpProxyMiddleware # tcp中间件
├── tcpProxyRouter # tcp路由
├── tcpserver # tcp服务
├── main.go 
├── go.mod
└── go.sum
```

##  快速开始

- 克隆

```bash
$ git clone https://github.com/HengY1Sky/Golang-Gate-Way.git
$ cd Golang-Gate-Way
```

- 后端

```bash
$ export GO111MODULE=on && export GOPROXY=https://goproxy.cn
$ cd backEnd
$ go mod tidy
$ go run main.go -endpoint dashboard # 开启面板管理
$ go run main.go -endpoint server # 开启服务
```

- 前端

```bash
$ npm-cli.js install --scripts-prepend-node-path=auto
$ npm run dev
```

- 数据库

使用任何的数据库管理工具，例如`Navicate`的`创建查询`在对应的数据库下创建数据表格

之后在`/backEnd/conf/dev/mysql_map`下进行编辑能让后端顺利启动就行

- Nginx

在面板的前台是要使用Nginx的Web服务进行转发到对应的接口中

首先要留意`/frontEnd/.env.production`下的生产模式下的转发路径

```nginx
  server {
        listen 8000;

        root /var/www/admin/dist; # 自己选择放哪里
        index index.html index.htm index.nginx-debian.html index.php;

        location / {
                try_files $uri $uri/ /index.html?$args;
        }

        location /prod-api/ {
            proxy_pass http://127.0.0.1:8002/;
        }
  }
```

> 其他的换成对应的生产环境，打包部署就好了，这个完全是能用的
>
> 但是缺点就是不能支持热更新，所以可以自己把`LoadOnce`重构一下就好了

总的来说，从中间件到一个个实现服务，从分开的一个个部分到网关的实现，知识点挺多的。

最后，**共勉**吧！🚀
