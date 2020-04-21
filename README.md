## Kingfisher Istio

Kubernetes中的微服务治理模块

## 依赖

- Golang: `Go >= 1.13`
- Istio: `isito >= 1.3.5`

## 说明

- istio client-go 更新 
    - 下载[https://github.com/istio/client-go/tree/master/pkg/apis](https://github.com/istio/client-go/tree/master/pkg/apis)到pkg目录下面，然后执行hack/code_generate.sh

## Makefile的使用

- 根据需求修改对应的REGISTRY变量，即可修改推送的仓库地址
- 编译成二进制文件： make build
- 生成镜像推送到镜像仓库： make push

