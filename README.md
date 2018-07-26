# A Simple Go Blog

基于gin构建的简单的博客，编写这个博客主要是为了学习Go语言

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