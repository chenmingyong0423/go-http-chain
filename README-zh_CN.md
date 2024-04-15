<h1 align="center">
  go-http-chain
</h1>

[![GitHub Repo stars](https://img.shields.io/github/stars/chenmingyong0423/go-http-chain)](https://github.com/chenmingyong0423/go-http-chain/stargazers)
[![GitHub issues](https://img.shields.io/github/issues/chenmingyong0423/go-http-chain)](https://github.com/chenmingyong0423/go-http-chain/issues)
[![GitHub License](https://img.shields.io/github/license/chenmingyong0423/go-http-chain)](https://github.com/chenmingyong0423/go-http-chain/blob/main/LICENSE)
[![GitHub release (with filter)](https://img.shields.io/github/v/release/chenmingyong0423/go-http-chain)](https://github.com/chenmingyong0423/go-http-chain)
[![Go Report Card](https://goreportcard.com/badge/github.com/chenmingyong0423/go-http-chain)](https://goreportcard.com/report/github.com/chenmingyong0423/go-http-chain)
[![All Contributors](https://img.shields.io/badge/all_contributors-1-orange.svg?style=flat-square)](#contributors-)

一个可链式调用的 Go HTTP 库，用于简化请求和响应管理。

---

[English](./README.md) | 中文简体

# Install
```bash
go get github.com/chenmingyong0423/go-http-chain
```
# Usage
```go
// 创建一个新默认的 client 并为 client 设置全局的 Header 参数
client := httpchain.NewDefault().SetHeader("X-Global-Param", "go-http-chain")
// 创建一个 Request 指定 GET 请求和设置 Header 参数，client 的 Header、Query 参数也会传入到该 Request
// 获取 *http.response
resp, err := client.Get("localhost:8080/test").SetHeader("name", "陈明勇").Call(context.Background()).Result()

// 直接解析响应结果到指定的结构体
var result map[string]any
err = client.Get("localhost:8080/test").Call(context.Background()).DecodeRespBody(&result)
```