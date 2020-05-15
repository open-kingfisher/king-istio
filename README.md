# Kingfisher Istio
[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)
[![Go Report Card](https://goreportcard.com/badge/github.com/open-kingfisher/king-istio)](https://goreportcard.com/report/github.com/open-kingfisher/king-istio)

Kubernetes中的微服务治理模块

## 依赖

- Golang: `Go >= 1.13`
- Istio: `isito >= 1.3.5`

## 说明

- istio client-go 使用
    - 下载[https://github.com/istio/client-go](https://github.com/istio/client-go)

## Makefile的使用

- 根据需求修改对应的REGISTRY变量，即可修改推送的仓库地址
- 编译成二进制文件： make build
- 生成镜像推送到镜像仓库： make push

## 联系我们
- [交流群](https://github.com/open-kingfisher/community/blob/master/contact_us/README.md)