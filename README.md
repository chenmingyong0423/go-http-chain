<h1 align="center">
  go-http-chain
</h1>

[![GitHub Repo stars](https://img.shields.io/github/stars/chenmingyong0423/go-http-chain)](https://github.com/chenmingyong0423/go-http-chain/stargazers)
[![GitHub issues](https://img.shields.io/github/issues/chenmingyong0423/go-http-chain)](https://github.com/chenmingyong0423/go-http-chain/issues)
[![GitHub License](https://img.shields.io/github/license/chenmingyong0423/go-http-chain)](https://github.com/chenmingyong0423/go-http-chain/blob/main/LICENSE)
[![GitHub release (with filter)](https://img.shields.io/github/v/release/chenmingyong0423/go-http-chain)](https://github.com/chenmingyong0423/go-http-chain)
[![Go Report Card](https://goreportcard.com/badge/github.com/chenmingyong0423/go-http-chain)](https://goreportcard.com/report/github.com/chenmingyong0423/go-http-chain)
[![All Contributors](https://img.shields.io/badge/all_contributors-1-orange.svg?style=flat-square)](#contributors-)

A chainable Go HTTP library for streamlined request and response management.

---

English | [中文简体](./README-zh_CN.md)

# Install
```bash
go get github.com/chenmingyong0423/go-http-chain
```
# Usage
```go
// Create a new default client and set global Header parameters for the client
client := httpchain.NewDefault().SetHeader("X-Global-Param", "go-http-chain")
// Create a Request specifying a GET request; the client's Headers and Query parameters will be passed to this Request
// Retrieve the *http.Response
resp, err := client.Get("localhost:8080/test").
		SetHeader("name", "Mingyong Chen").
		SetQuery("name", "Mingyong Chen").
		Do(context.Background())

// Directly parse the response result into a specified struct
var result map[string]any
err = client.Get("localhost:8080/test").
		SetHeader("name", "Mingyong Chen").
		SetQuery("name", "Mingyong Chen").
		DoAndParse(context.Background(), &result)
```
