# FIFSKY Blog for Golang

<a href="https://travis-ci.org/fifsky/blog"><img src="https://travis-ci.org/fifsky/blog.svg" alt="Build Status"></a>
<a href="https://opensource.org/licenses/mit-license.php" rel="nofollow"><img src="https://badges.frapsoft.com/os/mit/mit.svg?v=103"></a>


这是基于gin构建的简单的博客，这个博客从08年的ASP版本开始，一路经历了ASP、PHP、Golang版本，这是我个人的学习方法，每学习一门语言，都会使用新的语言重写这个博客，从而锻炼对一门语言细节的掌握

https://fifsky.com/

## 初始化
1、创建数据库blog

2、命令行执行
```
make install
```
3、运行
```
make build
make start
make stop
```

4、访问 127.0.0.1:8080

## 登录
默认登录用户名密码
test  123456

## 推荐
如果你要在同一个服务器上运行多个Go服务，或者想使用其他的静态文件服务，推荐使用 [Caddy](https://caddyserver.com/)

详见根目录`Caddyfile`

## 配置文件
```
config/local 开发环境
config/prod 线上环境
```

## 安装依赖包
由于一些依赖包需要代理才能安装，而最新的go module模式可以使用 https://goproxy.io/ 提供的代理