module github.com/open-kingfisher/king-istio

go 1.14

require (
	github.com/gin-gonic/gin v1.6.2
	github.com/open-kingfisher/king-k8s v0.0.0-20200428040900-99215549fa1a
	github.com/open-kingfisher/king-utils v0.0.0-20200422073733-6505a8c88560
	istio.io/client-go v0.0.0-20200427001039-e36a26d24fd2
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
)

replace (
	k8s.io/api => k8s.io/api v0.18.1
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.1
	k8s.io/client-go => k8s.io/client-go v0.18.1
)
