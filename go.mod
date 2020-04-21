module kingfisher/king-istio

go 1.12

require (
	github.com/gin-gonic/gin v1.4.0
	istio.io/api v0.0.0-20191210204543-73d66eb8c8d4
	k8s.io/api v0.0.0-20191206001707-7edad22604e1
	k8s.io/apimachinery v0.0.0-20191203211716-adc6f4cd9e7d
	k8s.io/client-go v11.0.0+incompatible
	kingfisher/kf v0.0.0-00010101000000-000000000000
	kingfisher/king-k8s v0.0.0-00010101000000-000000000000
)

replace (
	github.com/docker/docker => github.com/docker/docker v0.7.3-0.20190924004649-91870ed38213
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190620085101-78d2af792bab
	kingfisher/kf => ../kf
	kingfisher/king-k8s => ../king-k8s
)
