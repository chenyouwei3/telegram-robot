# 基于 Golang 的 Binlog 实时消息推送

## 简介

该项目基于 Golang 和 [bot 包](https://github.com/go-telegram/bot) 开发，使用 [go-mysql-org 包](https://github.com/go-mysql-org/go-mysql) 实现了 MySQL Binlog 的监听与解析，能够将数据库的变更实时推送至指定消息通道，适合需要监听数据库变化的实时处理场景。

## 功能

- 监听 MySQL Binlog 记录，捕捉数据增删改操作
- 实时推送数据库变更消息至指定消息队列或处理模块
- 具备高效和低延迟的消息推送能力

## 安装

注意将bot包当中的源码启动函数更改,将请求转发到运行程序的sock5端口

> proxyURL, _ := url.Parse("http://your-http-proxy.com:8080")
> transport := &http.Transport{ Proxy: http.ProxyURL(proxyURL),}


- 需要更改推送的用户,在notify.go当中的切片填入需要的uuid和chatid
- 在main.go当中修改自己推送机器人的token
- config.yaml编写数据库配置

